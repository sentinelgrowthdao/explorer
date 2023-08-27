package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
	nodetypes "github.com/sentinel-official/explorer/types/node"
	providertypes "github.com/sentinel-official/explorer/types/provider"
	sessiontypes "github.com/sentinel-official/explorer/types/session"
	subscriptiontypes "github.com/sentinel-official/explorer/types/subscription"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "02-data"
)

var (
	fromHeight int64
	toHeight   int64
	rpcAddress string
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.Int64Var(&fromHeight, "from-height", 901_801, "")
	flag.Int64Var(&toHeight, "to-height", 1_272_000, "")
	flag.StringVar(&rpcAddress, "rpc-address", "http://127.0.0.1:26657", "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func run(db *mongo.Database, height int64) (operations []types.DatabaseOperation, err error) {
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
			case "/sentinel.node.v1.MsgRegisterRequest", "/sentinel.node.v1.MsgService/MsgRegister":
				msg, err := nodetypes.NewMsgRegisterRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dNode := models.Node{
					Address:                msg.NodeAddress().String(),
					Provider:               msg.Provider,
					Price:                  msg.Price,
					RemoteURL:              msg.RemoteURL,
					RegisterHeight:         dBlock.Height,
					RegisterTimestamp:      dBlock.Time,
					RegisterTxHash:         dTxs[tIndex].Hash,
					Bandwidth:              nil,
					Handshake:              nil,
					IntervalSetSessions:    0,
					IntervalUpdateSessions: 0,
					IntervalUpdateStatus:   0,
					Location:               nil,
					Moniker:                "",
					Peers:                  0,
					QOS:                    nil,
					Status:                 hubtypes.StatusInactive.String(),
					Type:                   0,
					Version:                "",
					StatusHeight:           dBlock.Height,
					StatusTimestamp:        dBlock.Time,
					StatusTxHash:           dTxs[tIndex].Hash,
					ReachStatus:            nil,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.NodeInsertOne(ctx, db, &dNode); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.node.v1.MsgUpdateRequest", "/sentinel.node.v1.MsgService/MsgUpdate":
				msg, err := nodetypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"address": msg.From,
				}

				updateSet := bson.M{}
				if msg.Provider != "" {
					updateSet["provider"], updateSet["price"] = msg.Provider, nil
				}
				if msg.Price != nil && len(msg.Price) > 0 {
					updateSet["price"], updateSet["provider"] = msg.Price, ""
				}
				if msg.RemoteURL != "" {
					updateSet["remote_url"] = msg.RemoteURL
				}

				update := bson.M{
					"$set": updateSet,
				}
				projection := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
					if _, err := database.NodeFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.node.v1.MsgSetStatusRequest", "/sentinel.node.v1.MsgService/MsgSetStatus":
				msg, err := nodetypes.NewMsgSetStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"address": msg.From,
				}
				update := bson.M{
					"$set": bson.M{
						"status":           msg.Status,
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Hash,
					},
				}
				projection := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
					if _, err := database.NodeFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.plan.v1.MsgAddRequest", "/sentinel.plan.v1.MsgService/MsgAdd":
				return nil, fmt.Errorf("implement me")
			case "/sentinel.plan.v1.MsgSetStatusRequest", "/sentinel.plan.v1.MsgService/MsgSetStatus":
				return nil, fmt.Errorf("implement me")
			case "/sentinel.plan.v1.MsgAddNodeRequest", "/sentinel.plan.v1.MsgService/MsgAddNode":
				return nil, fmt.Errorf("implement me")
			case "/sentinel.plan.v1.MsgRemoveNodeRequest", "/sentinel.plan.v1.MsgService/MsgRemoveNode":
				return nil, fmt.Errorf("implement me")
			case "/sentinel.provider.v1.MsgRegisterRequest", "/sentinel.provider.v1.MsgService/MsgRegister":
				msg, err := providertypes.NewMsgRegisterRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dProvider := models.Provider{
					Address:           msg.ProvAddress().String(),
					Name:              msg.Name,
					Identity:          msg.Identity,
					Website:           msg.Website,
					Description:       msg.Description,
					RegisterHeight:    dBlock.Height,
					RegisterTimestamp: dBlock.Time,
					RegisterTxHash:    dTxs[tIndex].Hash,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.ProviderInsertOne(ctx, db, &dProvider); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.provider.v1.MsgUpdateRequest", "/sentinel.provider.v1.MsgService/MsgUpdate":
				msg, err := providertypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"address": msg.From,
				}

				updateSet := bson.M{}
				if msg.Name != "" {
					updateSet["name"] = msg.Name
				}
				if msg.Identity != "" {
					updateSet["identity"] = msg.Identity
				}
				if msg.Website != "" {
					updateSet["website"] = msg.Website
				}
				if msg.Description != "" {
					updateSet["description"] = msg.Description
				}

				update := bson.M{
					"$set": updateSet,
				}
				projection := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
					if _, err := database.ProviderFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.session.v1.MsgStartRequest", "/sentinel.session.v1.MsgService/MsgStart":
				msg, err := sessiontypes.NewMsgStartRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				eventStartSession, err := sessiontypes.NewEventStartSessionFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dSession := models.Session{
					ID:              eventStartSession.ID,
					Subscription:    msg.ID,
					Address:         msg.From,
					Node:            msg.Node,
					Duration:        0,
					Bandwidth:       nil,
					StartHeight:     dBlock.Height,
					StartTimestamp:  dBlock.Time,
					StartTxHash:     dTxs[tIndex].Hash,
					EndHeight:       0,
					EndTimestamp:    time.Time{},
					EndTxHash:       "",
					Payment:         nil,
					Rating:          0,
					Status:          hubtypes.StatusActive.String(),
					StatusHeight:    dBlock.Height,
					StatusTimestamp: dBlock.Time,
					StatusTxHash:    dTxs[tIndex].Hash,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.SessionInsertOne(ctx, db, &dSession); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.session.v1.MsgUpdateRequest", "/sentinel.session.v1.MsgService/MsgUpdate":
				msg, err := sessiontypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"id": msg.ID,
				}
				update := bson.M{
					"$set": bson.M{
						"duration":  msg.Duration,
						"bandwidth": msg.Bandwidth,
					},
				}
				projection := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
					if _, err := database.SessionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.session.v1.MsgEndRequest", "/sentinel.session.v1.MsgService/MsgEnd":
				msg, err := sessiontypes.NewMsgEndRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"id": msg.ID,
				}
				update := bson.M{
					"$set": bson.M{
						"rating":           msg.Rating,
						"status":           hubtypes.StatusInactivePending.String(),
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Index,
					},
				}
				projection := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
					if _, err := database.SessionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgSubscribeToNodeRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToNode":
				eventSubscribeToNode, err := subscriptiontypes.NewEventSubscribeToNodeFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dSubscription := models.Subscription{
					ID:              eventSubscribeToNode.ID,
					Owner:           eventSubscribeToNode.Owner,
					Node:            eventSubscribeToNode.Node,
					Price:           eventSubscribeToNode.Price,
					Deposit:         eventSubscribeToNode.Deposit,
					Plan:            0,
					Denom:           "",
					Expiry:          time.Time{},
					Payment:         nil,
					Free:            0,
					StartHeight:     dBlock.Height,
					StartTimestamp:  dBlock.Time,
					StartTxHash:     dTxs[tIndex].Hash,
					EndHeight:       0,
					EndTimestamp:    time.Time{},
					EndTxHash:       "",
					Status:          hubtypes.StatusActive.String(),
					StatusHeight:    dBlock.Height,
					StatusTimestamp: dBlock.Time,
					StatusTxHash:    dTxs[tIndex].Hash,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.SubscriptionInsertOne(ctx, db, &dSubscription); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgSubscribeToPlanRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToPlan":
				eventSubscribeToPlan, err := subscriptiontypes.NewEventSubscribeToPlanFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dSubscription := models.Subscription{
					ID:              eventSubscribeToPlan.ID,
					Owner:           eventSubscribeToPlan.Owner,
					Node:            "",
					Price:           nil,
					Deposit:         nil,
					Plan:            eventSubscribeToPlan.Plan,
					Denom:           eventSubscribeToPlan.Denom,
					Expiry:          eventSubscribeToPlan.Expiry,
					Payment:         eventSubscribeToPlan.Payment,
					Free:            0,
					StartHeight:     dBlock.Height,
					StartTimestamp:  dBlock.Time,
					StartTxHash:     dTxs[tIndex].Hash,
					EndHeight:       0,
					EndTimestamp:    time.Time{},
					EndTxHash:       "",
					Status:          hubtypes.StatusActive.String(),
					StatusHeight:    dBlock.Height,
					StatusTimestamp: dBlock.Time,
					StatusTxHash:    dTxs[tIndex].Hash,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.SubscriptionInsertOne(ctx, db, &dSubscription); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgCancelRequest", "/sentinel.subscription.v1.MsgService/MsgCancel":
				msg, err := subscriptiontypes.NewMsgCancelRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"id": msg.ID,
				}
				update := bson.M{
					"$set": bson.M{
						"status":           hubtypes.StatusInactivePending.String(),
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Index,
					},
				}
				projection := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
					if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgAddQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgAddQuota":
				return nil, fmt.Errorf("impement me")
			case "/sentinel.subscription.v1.MsgUpdateQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgUpdateQuota":
				return nil, fmt.Errorf("impement me")
			default:

			}
		}
	}

	log.Println("EndBlockEventsLen", dBlock.Height, len(dBlock.EndBlockEvents))
	for eIndex := 0; eIndex < len(dBlock.EndBlockEvents); eIndex++ {
		log.Println("Type", eIndex, dBlock.EndBlockEvents[eIndex].Type)
		switch dBlock.EndBlockEvents[eIndex].Type {
		case "sentinel.node.v1.EventSetNodeStatus":
			event, err := nodetypes.NewEventSetNodeStatus(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter := bson.M{
				"address": event.Address,
			}
			update := bson.M{
				"$set": bson.M{
					"status":           event.Status,
					"status_height":    dBlock.Height,
					"status_timestamp": dBlock.Time,
					"status_tx_hash":   "",
				},
			}
			projection := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true)
				if _, err := database.NodeFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
					return err
				}

				return nil
			})
		default:

		}
	}

	return operations, nil
}

func main() {
	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.TODO(), nil); err != nil {
		log.Fatalln(err)
	}

	filter := bson.M{
		"app_name": appName,
	}

	dSyncStatus, err := database.SyncStatusFindOne(context.TODO(), db, filter)
	if err != nil {
		log.Fatalln(err)
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

		operations, err := run(db, height)
		if err != nil {
			log.Fatalln(err)
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

				log.Println("OperationsLen", len(operations))
				for i := 0; i < len(operations); i++ {
					if err := operations[i](ctx); err != nil {
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
			log.Fatalln(err)
		}
	}
}
