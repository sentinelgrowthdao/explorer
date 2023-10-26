package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/operations"
	"github.com/sentinel-official/explorer/types"
	nodetypes "github.com/sentinel-official/explorer/types/node"
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
	flag.Int64Var(&fromHeight, "from-height", 12_310_005, "")
	flag.Int64Var(&toHeight, "to-height", math.MaxInt64, "")
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
				bson.E{Key: "session_id", Value: 1},
				bson.E{Key: "type", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
			},
			Options: options.Index().SetPartialFilterExpression(
				bson.M{
					"session_id": bson.M{
						"$exists": true,
					},
					"type": types.EventTypeSessionUpdateDetails,
				},
			),
		},
	}

	_, err := database.EventIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "addr", Value: 1},
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

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "id", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.SessionIndexesCreateMany(ctx, db, indexes)
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
		"begin_block_events": 1,
		"end_block_events":   1,
		"height":             1,
		"time":               1,
	}

	dBlock, err := database.BlockFindOne(context.TODO(), db, filter, options.FindOne().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	if dBlock == nil {
		return nil, fmt.Errorf("block %d does not exist", height)
	}

	log.Println("BeginBlockEventsLen", dBlock.Height, len(dBlock.BeginBlockEvents))
	for eIndex := 0; eIndex < len(dBlock.BeginBlockEvents); eIndex++ {
		log.Println("Type", eIndex, dBlock.BeginBlockEvents[eIndex].Type)
		switch dBlock.BeginBlockEvents[eIndex].Type {
		case "sentinel.subscription.v2.EventPayForPayout":
			event, err := subscriptiontypes.NewEventPayForPayout(dBlock.BeginBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewNodeStatisticUpdateEarningsForHours(
					db, utils.DayDate(dBlock.Time), event.ID, event.Payment,
				),
			)
		default:

		}
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
		dTxs[tIndex].Messages = dTxs[tIndex].Messages.WithAuthzMsgExecMessages()
		log.Println("TxHash", dTxs[tIndex].Hash)
		log.Println("MessagesLen", tIndex, len(dTxs[tIndex].Messages))

		for mIndex := 0; mIndex < len(dTxs[tIndex].Messages); mIndex++ {
			log.Println("Type", dTxs[tIndex].Messages[mIndex].Type)
			switch dTxs[tIndex].Messages[mIndex].Type {
			case "/sentinel.session.v2.MsgStartRequest", "/sentinel.session.v2.MsgService/MsgStart":
				msg, err := sessiontypes.NewMsgStartRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				ops = append(
					ops,
					operations.NewNodeStatisticIncreaseSessionStartCount(
						db, msg.NodeAddress, utils.DayDate(dBlock.Time), 1,
					),
				)
			case "/sentinel.session.v2.MsgUpdateDetailsRequest", "/sentinel.session.v2.MsgService/MsgUpdateDetails":
				msg, err := sessiontypes.NewMsgUpdateDetailsRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				ops = append(
					ops,
					operations.NewNodeStatisticUpdateSessionDetails(
						db, msg.From, utils.DayDate(dBlock.Time), msg.ID, msg.Bandwidth, msg.Duration,
					),
				)
			case "/sentinel.node.v2.MsgSubscribeRequest", "/sentinel.node.v2.MsgService/MsgSubscribe":
				msg, err := nodetypes.NewMsgSubscribeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				ops = append(
					ops,
					operations.NewNodeStatisticIncreaseSubscriptionStartCount(
						db, msg.NodeAddress, utils.DayDate(dBlock.Time), 1,
					),
				)
				if msg.Gigabytes != 0 {
					ops = append(
						ops,
						operations.NewNodeStatisticUpdateSubscriptionBytes(
							db, msg.NodeAddress, utils.DayDate(dBlock.Time), hubtypes.Gigabyte.MulRaw(msg.Gigabytes),
						),
					)
				}
				if msg.Hours != 0 {
					ops = append(
						ops,
						operations.NewNodeStatisticUpdateSubscriptionHours(
							db, msg.NodeAddress, utils.DayDate(dBlock.Time), msg.Hours,
						),
					)
				}
			default:

			}
		}
	}

	log.Println("EndBlockEventsLen", dBlock.Height, len(dBlock.EndBlockEvents))
	for eIndex := 0; eIndex < len(dBlock.EndBlockEvents); eIndex++ {
		log.Println("Type", eIndex, dBlock.EndBlockEvents[eIndex].Type)
		switch dBlock.EndBlockEvents[eIndex].Type {
		case "sentinel.session.v2.EventUpdateStatus":
			event, err := sessiontypes.NewEventUpdateStatus(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}
			if event.Status != hubtypes.StatusInactive.String() {
				continue
			}

			ops = append(
				ops,
				operations.NewNodeStatisticIncreaseSessionEndCount(
					db, utils.DayDate(dBlock.Time), event.ID, 1,
				),
			)
		case "sentinel.subscription.v2.EventPayForSession":
			event, err := subscriptiontypes.NewEventPayForSession(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewNodeStatisticUpdateEarningsForBytes(
					db, utils.DayDate(dBlock.Time), event.ID, event.Payment,
				),
			)
		case "sentinel.subscription.v2.EventUpdateStatus":
			event, err := subscriptiontypes.NewEventUpdateStatus(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}
			if event.Status != hubtypes.StatusInactive.String() {
				continue
			}

			ops = append(
				ops,
				operations.NewNodeStatisticIncreaseSubscriptionEndCount(
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
