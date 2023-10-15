package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/operations"
	"github.com/sentinel-official/explorer/types"
	sessiontypes "github.com/sentinel-official/explorer/types/session"
	subscriptiontypes "github.com/sentinel-official/explorer/types/subscription"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "06_node-statistics"
)

var (
	fromHeight int64
	toHeight   int64
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.Int64Var(&fromHeight, "from-height", 901_801, "")
	flag.Int64Var(&toHeight, "to-height", 5_125_000, "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "type", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
				bson.E{Key: "session_id", Value: 1},
			},
		},
	}

	_, err := database.EventIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "address", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.NodeStatisticIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	return nil
}

func run(db *mongo.Database, height int64) (ops []types.DatabaseOperation, err error) {
	filter := bson.M{
		"height": height,
	}
	projection := bson.M{
		"height":           1,
		"time":             1,
		"end_block_events": 1,
	}

	dBlock, err := database.BlockFindOne(context.TODO(), db, filter, options.FindOne().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	if dBlock == nil {
		return nil, fmt.Errorf("block %d does not exist", height)
	}

	filter = bson.M{
		"height":      height,
		"result.code": 0,
	}
	projection = bson.M{
		"hash":       1,
		"messages":   1,
		"result.log": 1,
	}

	dTxs, err := database.TxFind(context.TODO(), db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	log.Println("TxsLen", len(dTxs))
	for tIndex := 0; tIndex < len(dTxs); tIndex++ {
		log.Println("TxHash", dTxs[tIndex].Hash)
		log.Println("MessagesLen", tIndex, len(dTxs[tIndex].Messages))

		txResultLog := types.NewABCIMessageLogs(dTxs[tIndex].Result.Log)
		for mIndex := 0; mIndex < len(dTxs[tIndex].Messages); mIndex++ {
			log.Println("Type", dTxs[tIndex].Messages[mIndex].Type)
			switch dTxs[tIndex].Messages[mIndex].Type {
			case "/sentinel.session.v1.MsgStartRequest", "/sentinel.session.v1.MsgService/MsgStart":
				msg, err := sessiontypes.NewMsgStartRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				ops = append(
					ops,
					operations.NewNodeStatisticUpdateSessionStartCount(
						db, msg.Node, utils.DayDate(dBlock.Time), 1,
					),
				)
			case "/sentinel.session.v1.MsgUpdateRequest", "/sentinel.session.v1.MsgService/MsgUpdate":
				msg, err := sessiontypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				ops = append(
					ops,
					operations.NewNodeStatisticUpdateSessionDetails(
						db, msg.From, utils.DayDate(dBlock.Time), msg.ID, msg.Bandwidth, msg.Duration,
					),
				)
			case "/sentinel.subscription.v1.MsgSubscribeToNodeRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToNode":
				msg, err := subscriptiontypes.NewMsgSubscribeToNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				event, err := subscriptiontypes.NewEventAddQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				ops = append(
					ops,
					operations.NewNodeStatisticUpdateSubscriptionStartCount(
						db, msg.Address, utils.DayDate(dBlock.Time), 1,
					),
					operations.NewNodeStatisticUpdateSubscriptionBytes(
						db, msg.Address, utils.DayDate(dBlock.Time), utils.MustIntFromString(event.Allocated),
					),
				)
			default:

			}
		}
	}

	log.Println("EndBlockEventsLen", dBlock.Height, len(dBlock.EndBlockEvents))
	for eIndex := 0; eIndex < len(dBlock.EndBlockEvents); eIndex++ {
		log.Println("Type", eIndex, dBlock.EndBlockEvents[eIndex].Type)
		switch dBlock.EndBlockEvents[eIndex].Type {
		case "sentinel.session.v1.EventEndSession":
			event, err := sessiontypes.NewEventEndSession(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewNodeStatisticUpdateSessionEndCount(
					db, utils.DayDate(dBlock.Time), event.ID, 1,
				),
			)
		case "sentinel.session.v1.EventPay":
			event, err := sessiontypes.NewEventPay(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewNodeStatisticUpdateEarningsForBytes(
					db, utils.DayDate(dBlock.Time), event.ID, event.Payment,
				),
			)
		case "sentinel.subscription.v1.EventCancelSubscription":
			event, err := subscriptiontypes.NewEventCancelSubscription(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewNodeStatisticUpdateSubscriptionEndCount(
					db, utils.DayDate(dBlock.Time), event.ID, 1,
				),
			)
		default:

		}
	}

	return ops, nil
}

func main() {
	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Panicln(err)
	}

	if err = db.Client().Ping(context.TODO(), nil); err != nil {
		log.Panicln(err)
	}

	if err := createIndexes(context.TODO(), db); err != nil {
		log.Panicln(err)
	}

	filter := bson.M{
		"app_name": appName,
	}

	dSyncStatus, err := database.SyncStatusFindOne(context.TODO(), db, filter)
	if err != nil {
		log.Panicln(err)
	}
	if dSyncStatus == nil {
		dSyncStatus = &models.SyncStatus{
			AppName:   appName,
			Height:    fromHeight - 1,
			Timestamp: time.Time{},
		}
	}

	height := dSyncStatus.Height + 1
	for height < toHeight {
		now := time.Now()
		log.Println("Height", height)

		ops, err := run(db, height)
		if err != nil {
			log.Panicln(err)
		}

		log.Println("OperationsLen", len(ops))
		if len(ops) == 0 {
			height++
			continue
		}

		err = db.Client().UseSession(
			context.TODO(),
			func(ctx mongo.SessionContext) error {
				err := ctx.StartTransaction(
					options.Transaction().
						SetReadConcern(readconcern.Snapshot()).
						SetWriteConcern(writeconcern.Majority()),
				)
				if err != nil {
					return err
				}

				abort := true
				defer func() {
					if abort {
						_ = ctx.AbortTransaction(ctx)
					}
				}()

				for i := 0; i < len(ops); i++ {
					if err := ops[i](ctx); err != nil {
						return err
					}
				}

				filter := bson.M{
					"app_name": appName,
				}
				update := bson.M{
					"$set": bson.M{
						"height": height,
					},
				}
				projection := bson.M{
					"_id": 1,
				}

				_, err = database.SyncStatusFindOneAndUpdate(ctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
				if err != nil {
					return err
				}

				height++

				abort = false
				return ctx.CommitTransaction(ctx)
			},
		)

		log.Println("Duration", time.Since(now))
		log.Println("")
		if err != nil {
			log.Panicln(err)
		}
	}
}
