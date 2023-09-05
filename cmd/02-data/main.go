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
	deposittypes "github.com/sentinel-official/explorer/types/deposit"
	nodetypes "github.com/sentinel-official/explorer/types/node"
	plantypes "github.com/sentinel-official/explorer/types/plan"
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

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "app_name", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err := database.SyncStatusIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "height", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.BlockIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "height", Value: 1},
				bson.E{Key: "result.code", Value: 1},
			},
			Options: options.Index().
				SetPartialFilterExpression(
					bson.M{
						"result.code": 0,
					},
				),
		},
	}

	_, err = database.TxIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "address", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.DepositIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "address", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.NodeIndexesCreateMany(ctx, db, indexes)
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

	_, err = database.PlanIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "address", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.ProviderIndexesCreateMany(ctx, db, indexes)
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

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "id", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.SubscriptionIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "id", Value: 1},
				bson.E{Key: "address", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.SubscriptionQuotaIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	return nil
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
					Address:           msg.NodeAddress().String(),
					Provider:          msg.Provider,
					Price:             msg.Price,
					RemoteURL:         msg.RemoteURL,
					RegisterHeight:    dBlock.Height,
					RegisterTimestamp: dBlock.Time,
					RegisterTxHash:    dTxs[tIndex].Hash,
					Status:            hubtypes.StatusInactive.String(),
					StatusHeight:      dBlock.Height,
					StatusTimestamp:   dBlock.Time,
					StatusTxHash:      dTxs[tIndex].Hash,
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

				filter1 := bson.M{
					"address": msg.From,
				}

				updateSet1 := bson.M{}
				if msg.Provider != "" {
					updateSet1["provider"], updateSet1["price"] = msg.Provider, nil
				}
				if msg.Price != nil && len(msg.Price) > 0 {
					updateSet1["price"], updateSet1["provider"] = msg.Price, ""
				}
				if msg.RemoteURL != "" {
					updateSet1["remote_url"] = msg.RemoteURL
				}

				update1 := bson.M{
					"$set": updateSet1,
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.NodeFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:        types.EventTypeNodeUpdateDetails,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					NodeAddress: msg.From,
					ProvAddress: msg.Provider,
					Price:       msg.Price,
					RemoteURL:   msg.RemoteURL,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.node.v1.MsgSetStatusRequest", "/sentinel.node.v1.MsgService/MsgSetStatus":
				msg, err := nodetypes.NewMsgSetStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"address": msg.From,
				}
				update1 := bson.M{
					"$set": bson.M{
						"status":           msg.Status,
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Hash,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.NodeFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:        types.EventTypeNodeUpdateStatus,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					NodeAddress: msg.From,
					Status:      msg.Status,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.plan.v1.MsgAddRequest", "/sentinel.plan.v1.MsgService/MsgAdd":
				msg, err := plantypes.NewMsgAddRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				eventAdd, err := plantypes.NewEventAddFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dPlan := models.Plan{
					ID:              eventAdd.ID,
					ProviderAddress: msg.From,
					Price:           msg.Price,
					Validity:        msg.Validity,
					Bytes:           msg.Bytes,
					NodeAddresses:   []string{},
					AddHeight:       dBlock.Height,
					AddTimestamp:    dBlock.Time,
					AddTxHash:       dTxs[tIndex].Hash,
					Status:          hubtypes.StatusInactive.String(),
					StatusHeight:    dBlock.Height,
					StatusTimestamp: dBlock.Time,
					StatusTxHash:    dTxs[tIndex].Hash,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.PlanInsertOne(ctx, db, &dPlan); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.plan.v1.MsgSetStatusRequest", "/sentinel.plan.v1.MsgService/MsgSetStatus":
				msg, err := plantypes.NewMsgSetStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"id": msg.ID,
				}
				update1 := bson.M{
					"$set": bson.M{
						"status":           msg.Status,
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Index,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.PlanFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:      types.EventTypePlanUpdateStatus,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					PlanID:    msg.ID,
					Status:    msg.Status,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.plan.v1.MsgAddNodeRequest", "/sentinel.plan.v1.MsgService/MsgAddNode":
				msg, err := plantypes.NewMsgAddNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"id": msg.ID,
				}
				update1 := bson.M{
					"$push": bson.M{
						"node_addresses": msg.Address,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.PlanFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:        types.EventTypePlanAddNode,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					PlanID:      msg.ID,
					NodeAddress: msg.Address,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.plan.v1.MsgRemoveNodeRequest", "/sentinel.plan.v1.MsgService/MsgRemoveNode":
				msg, err := plantypes.NewMsgRemoveNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"id": msg.ID,
				}
				update1 := bson.M{
					"$pull": bson.M{
						"node_addresses": msg.Address,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.PlanFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:        types.EventTypePlanRemoveNode,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					PlanID:      msg.ID,
					NodeAddress: msg.Address,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
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

				filter1 := bson.M{
					"address": msg.From,
				}

				updateSet1 := bson.M{}
				if msg.Name != "" {
					updateSet1["name"] = msg.Name
				}
				if msg.Identity != "" {
					updateSet1["identity"] = msg.Identity
				}
				if msg.Website != "" {
					updateSet1["website"] = msg.Website
				}
				if msg.Description != "" {
					updateSet1["description"] = msg.Description
				}

				update1 := bson.M{
					"$set": updateSet1,
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.ProviderFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:        types.EventTypeProviderUpdateDetails,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					ProvAddress: msg.From,
					Name:        msg.Name,
					Identity:    msg.Identity,
					Website:     msg.Website,
					Description: msg.Description,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
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
					Bandwidth:       nil,
					Duration:        0,
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

				filter1 := bson.M{
					"id": msg.ID,
				}
				update1 := bson.M{
					"$set": bson.M{
						"bandwidth": msg.Bandwidth,
						"duration":  msg.Duration,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.SessionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:      types.EventTypeSessionUpdateDetails,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					SessionID: msg.ID,
					Bandwidth: msg.Bandwidth,
					Duration:  msg.Duration,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.session.v1.MsgEndRequest", "/sentinel.session.v1.MsgService/MsgEnd":
				msg, err := sessiontypes.NewMsgEndRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"id": msg.ID,
				}
				update1 := bson.M{
					"$set": bson.M{
						"rating":           msg.Rating,
						"status":           hubtypes.StatusInactivePending.String(),
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Index,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.SessionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:      types.EventTypeSessionUpdateStatus,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					SessionID: msg.ID,
					Status:    hubtypes.StatusInactivePending.String(),
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
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
					Free:            0,
					Node:            eventSubscribeToNode.Node,
					Price:           eventSubscribeToNode.Price,
					Deposit:         eventSubscribeToNode.Deposit,
					Refund:          nil,
					Plan:            0,
					Denom:           "",
					Expiry:          time.Time{},
					Payment:         nil,
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

				eventAddQuota, err := subscriptiontypes.NewEventAddQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dSubscriptionQuota := models.SubscriptionQuota{
					ID:        eventAddQuota.ID,
					Address:   eventAddQuota.Address,
					Allocated: eventAddQuota.Allocated,
					Consumed:  eventAddQuota.Consumed,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.SubscriptionQuotaInsertOne(ctx, db, &dSubscriptionQuota); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionQuotaUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAddQuota.ID,
					AccAddress:     eventAddQuota.Address,
					Allocated:      eventAddQuota.Allocated,
					Consumed:       eventAddQuota.Consumed,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})

				eventAdd, err := deposittypes.NewEventAddFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"address": eventAdd.Address,
				}
				update1 := bson.M{
					"$set": bson.M{
						"coins":     eventAdd.Current,
						"height":    dBlock.Height,
						"timestamp": dBlock.Time,
						"tx_hash":   dTxs[tIndex].Hash,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.DepositFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent2 := models.Event{
					Type:       types.EventTypeDepositAdd,
					Height:     dBlock.Height,
					Timestamp:  dBlock.Time,
					TxHash:     dTxs[tIndex].Hash,
					AccAddress: eventAdd.Address,
					Coins:      eventAdd.Coins,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent2); err != nil {
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
					Free:            0,
					Node:            "",
					Price:           nil,
					Deposit:         nil,
					Refund:          nil,
					Plan:            eventSubscribeToPlan.Plan,
					Denom:           eventSubscribeToPlan.Payment.Denom,
					Expiry:          eventSubscribeToPlan.Expiry,
					Payment:         eventSubscribeToPlan.Payment,
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

				eventAddQuota, err := subscriptiontypes.NewEventAddQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dSubscriptionQuota := models.SubscriptionQuota{
					ID:        eventAddQuota.ID,
					Address:   eventAddQuota.Address,
					Allocated: eventAddQuota.Allocated,
					Consumed:  eventAddQuota.Consumed,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.SubscriptionQuotaInsertOne(ctx, db, &dSubscriptionQuota); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionQuotaUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAddQuota.ID,
					AccAddress:     eventAddQuota.Address,
					Allocated:      eventAddQuota.Allocated,
					Consumed:       eventAddQuota.Consumed,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgCancelRequest", "/sentinel.subscription.v1.MsgService/MsgCancel":
				msg, err := subscriptiontypes.NewMsgCancelRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"id": msg.ID,
				}
				update1 := bson.M{
					"$set": bson.M{
						"status":           hubtypes.StatusInactivePending.String(),
						"status_height":    dBlock.Height,
						"status_timestamp": dBlock.Time,
						"status_tx_hash":   dTxs[tIndex].Index,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateStatus,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: msg.ID,
					Status:         hubtypes.StatusInactivePending.String(),
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgAddQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgAddQuota":
				eventAddQuota, err := subscriptiontypes.NewEventAddQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dSubscriptionQuota := models.SubscriptionQuota{
					ID:        eventAddQuota.ID,
					Address:   eventAddQuota.Address,
					Consumed:  eventAddQuota.Consumed,
					Allocated: eventAddQuota.Allocated,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.SubscriptionQuotaInsertOne(ctx, db, &dSubscriptionQuota); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionQuotaUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAddQuota.ID,
					AccAddress:     eventAddQuota.Address,
					Allocated:      eventAddQuota.Allocated,
					Consumed:       eventAddQuota.Consumed,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})

				filter1 := bson.M{
					"id": eventAddQuota.ID,
				}
				update1 := bson.M{
					"$set": bson.M{
						"free": eventAddQuota.Free,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent2 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAddQuota.ID,
					Free:           eventAddQuota.Free,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent2); err != nil {
						return err
					}

					return nil
				})
			case "/sentinel.subscription.v1.MsgUpdateQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgUpdateQuota":
				eventUpdateQuota, err := subscriptiontypes.NewEventUpdateQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				filter1 := bson.M{
					"id":      eventUpdateQuota.ID,
					"address": eventUpdateQuota.Address,
				}
				update1 := bson.M{
					"$set": bson.M{
						"allocated": eventUpdateQuota.Allocated,
					},
				}
				projection1 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
					if _, err := database.SubscriptionQuotaFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionQuotaUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventUpdateQuota.ID,
					AccAddress:     eventUpdateQuota.Address,
					Allocated:      eventUpdateQuota.Allocated,
					Consumed:       eventUpdateQuota.Consumed,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
						return err
					}

					return nil
				})

				filter2 := bson.M{
					"id": eventUpdateQuota.ID,
				}
				update2 := bson.M{
					"$set": bson.M{
						"free": eventUpdateQuota.Free,
					},
				}
				projection2 := bson.M{
					"_id": 1,
				}

				operations = append(operations, func(ctx mongo.SessionContext) error {
					opts := options.FindOneAndUpdate().SetProjection(projection2).SetUpsert(true)
					if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter2, update2, opts); err != nil {
						return err
					}

					return nil
				})

				dEvent2 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventUpdateQuota.ID,
					Free:           eventUpdateQuota.Free,
				}
				operations = append(operations, func(ctx mongo.SessionContext) error {
					if _, err := database.EventInsertOne(ctx, db, &dEvent2); err != nil {
						return err
					}

					return nil
				})
			default:

			}
		}
	}

	log.Println("EndBlockEventsLen", dBlock.Height, len(dBlock.EndBlockEvents))
	for eIndex := 0; eIndex < len(dBlock.EndBlockEvents); eIndex++ {
		log.Println("Type", eIndex, dBlock.EndBlockEvents[eIndex].Type)
		switch dBlock.EndBlockEvents[eIndex].Type {
		case "sentinel.deposit.v1.EventSubtract":
			event, err := deposittypes.NewEventSubtract(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"address": event.Address,
			}
			update1 := bson.M{
				"$set": bson.M{
					"coins":     event.Current,
					"height":    dBlock.Height,
					"timestamp": dBlock.Time,
					"tx_hash":   "",
				},
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.DepositFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})

			dEvent1 := models.Event{
				Type:       types.EventTypeDepositSubtract,
				Height:     dBlock.Height,
				Timestamp:  dBlock.Time,
				TxHash:     "",
				AccAddress: event.Address,
				Coins:      event.Coins,
			}
			operations = append(operations, func(ctx mongo.SessionContext) error {
				if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
					return err
				}

				return nil
			})
		case "sentinel.node.v1.EventSetNodeStatus":
			event, err := nodetypes.NewEventSetNodeStatus(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"address": event.Address,
			}
			update1 := bson.M{
				"$set": bson.M{
					"status":           event.Status,
					"status_height":    dBlock.Height,
					"status_timestamp": dBlock.Time,
					"status_tx_hash":   "",
				},
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.NodeFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})

			dEvent1 := models.Event{
				Type:        types.EventTypeNodeUpdateStatus,
				Height:      dBlock.Height,
				Timestamp:   dBlock.Time,
				TxHash:      "",
				NodeAddress: event.Address,
				Status:      event.Status,
			}
			operations = append(operations, func(ctx mongo.SessionContext) error {
				if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
					return err
				}

				return nil
			})
		case "sentinel.session.v1.EventEndSession":
			event, err := sessiontypes.NewEventEndSession(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"id": event.ID,
			}

			updateSet1 := bson.M{
				"rating":           0,
				"status":           event.Status,
				"status_height":    dBlock.Height,
				"status_timestamp": dBlock.Time,
				"status_tx_hash":   "",
			}
			if event.Status == hubtypes.StatusInactive.String() {
				updateSet1["end_height"] = dBlock.Height
				updateSet1["end_timestamp"] = dBlock.Time
				updateSet1["end_tx_hash"] = ""
			}

			update1 := bson.M{
				"$set": updateSet1,
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.SessionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})

			dEvent1 := models.Event{
				Type:      types.EventTypeSessionUpdateStatus,
				Height:    dBlock.Height,
				Timestamp: dBlock.Time,
				TxHash:    "",
				SessionID: event.ID,
				Status:    event.Status,
			}
			operations = append(operations, func(ctx mongo.SessionContext) error {
				if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
					return err
				}

				return nil
			})
		case "sentinel.session.v1.EventPay":
			event, err := sessiontypes.NewEventPay(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"id": event.ID,
			}

			update1 := bson.M{
				"$set": bson.M{
					"payment": event.Payment,
				},
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.SessionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})
		case "sentinel.subscription.v1.EventCancelSubscription":
			event, err := subscriptiontypes.NewEventCancelSubscription(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"id": event.ID,
			}

			updateSet1 := bson.M{
				"status":           event.Status,
				"status_height":    dBlock.Height,
				"status_timestamp": dBlock.Time,
				"status_tx_hash":   "",
			}
			if event.Status == hubtypes.StatusInactive.String() {
				updateSet1["end_height"] = dBlock.Height
				updateSet1["end_timestamp"] = dBlock.Time
				updateSet1["end_tx_hash"] = ""
			}

			update1 := bson.M{
				"$set": updateSet1,
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})

			dEvent1 := models.Event{
				Type:           types.EventTypeSubscriptionUpdateStatus,
				Height:         dBlock.Height,
				Timestamp:      dBlock.Time,
				TxHash:         "",
				SubscriptionID: event.ID,
				Status:         event.Status,
			}
			operations = append(operations, func(ctx mongo.SessionContext) error {
				if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
					return err
				}

				return nil
			})
		case "sentinel.subscription.v1.EventRefund":
			event, err := subscriptiontypes.NewEventRefund(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"id": event.ID,
			}

			update1 := bson.M{
				"$set": bson.M{
					"refund": event.Refund,
				},
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})
		case "sentinel.subscription.v1.EventUpdateQuota":
			event, err := subscriptiontypes.NewEventUpdateQuota(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			filter1 := bson.M{
				"id":      event.ID,
				"address": event.Address,
			}
			update1 := bson.M{
				"$set": bson.M{
					"consumed": event.Consumed,
				},
			}
			projection1 := bson.M{
				"_id": 1,
			}

			operations = append(operations, func(ctx mongo.SessionContext) error {
				opts := options.FindOneAndUpdate().SetProjection(projection1).SetUpsert(true)
				if _, err := database.SubscriptionQuotaFindOneAndUpdate(ctx, db, filter1, update1, opts); err != nil {
					return err
				}

				return nil
			})

			dEvent1 := models.Event{
				Type:           types.EventTypeSubscriptionQuotaUpdateDetails,
				Height:         dBlock.Height,
				Timestamp:      dBlock.Time,
				TxHash:         "",
				SubscriptionID: event.ID,
				AccAddress:     event.Address,
				Allocated:      event.Allocated,
				Consumed:       event.Consumed,
			}
			operations = append(operations, func(ctx mongo.SessionContext) error {
				if _, err := database.EventInsertOne(ctx, db, &dEvent1); err != nil {
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

		operations, err := run(db, height)
		if err != nil {
			log.Panicln(err)
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
			log.Panicln(err)
		}
	}
}
