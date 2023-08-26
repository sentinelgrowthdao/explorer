package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/querier"
	"github.com/sentinel-official/explorer/types"
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

func run(db *mongo.Database, q *querier.Querier, height int64) (operations []types.DatabaseOperation, err error) {
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
				from := dTxs[tIndex].Messages[mIndex].Data["from"].(string)
				provider := dTxs[tIndex].Messages[mIndex].Data["provider"].(string)
				remoteURL := dTxs[tIndex].Messages[mIndex].Data["remote_url"].(string)

				buf, err := json.Marshal(dTxs[tIndex].Messages[mIndex].Data["price"])
				if err != nil {
					return nil, err
				}

				var price sdk.Coins
				if err := json.Unmarshal(buf, &price); err != nil {
					return nil, err
				}

				fromAddr := utils.MustAccAddressFromBech32(from)
				nodeAddr := hubtypes.NodeAddress(fromAddr.Bytes())

				dNode := types.Node{
					Address:                nodeAddr.String(),
					Provider:               provider,
					Price:                  types.NewCoins(price),
					RemoteURL:              remoteURL,
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
				from := dTxs[tIndex].Messages[mIndex].Data["from"].(string)
				provider := dTxs[tIndex].Messages[mIndex].Data["provider"].(string)
				remoteURL := dTxs[tIndex].Messages[mIndex].Data["remote_url"].(string)

				buf, err := json.Marshal(dTxs[tIndex].Messages[mIndex].Data["price"])
				if err != nil {
					return nil, err
				}

				var price sdk.Coins
				if err := json.Unmarshal(buf, &price); err != nil {
					return nil, err
				}

				filter := bson.M{
					"address": from,
				}
				update := bson.M{
					"$set": bson.M{
						"provider":   provider,
						"price":      types.NewCoins(price),
						"remote_url": remoteURL,
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
			case "/sentinel.node.v1.MsgSetStatusRequest", "/sentinel.node.v1.MsgService/MsgSetStatus":
				from := dTxs[tIndex].Messages[mIndex].Data["from"].(string)
				status := dTxs[tIndex].Messages[mIndex].Data["status"].(string)

				filter := bson.M{
					"address": from,
				}
				update := bson.M{
					"$set": bson.M{
						"status":           status,
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
				from := dTxs[tIndex].Messages[mIndex].Data["from"].(string)
				name := dTxs[tIndex].Messages[mIndex].Data["name"].(string)
				identity := dTxs[tIndex].Messages[mIndex].Data["identity"].(string)
				website := dTxs[tIndex].Messages[mIndex].Data["website"].(string)
				description := dTxs[tIndex].Messages[mIndex].Data["description"].(string)

				fromAddr := utils.MustAccAddressFromBech32(from)
				provAddr := hubtypes.ProvAddress(fromAddr.Bytes())

				dProvider := types.Provider{
					Address:           provAddr.String(),
					Name:              name,
					Identity:          identity,
					Website:           website,
					Description:       description,
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
				from := dTxs[tIndex].Messages[mIndex].Data["from"].(string)
				name := dTxs[tIndex].Messages[mIndex].Data["name"].(string)
				identity := dTxs[tIndex].Messages[mIndex].Data["identity"].(string)
				website := dTxs[tIndex].Messages[mIndex].Data["website"].(string)
				description := dTxs[tIndex].Messages[mIndex].Data["description"].(string)

				filter := bson.M{
					"address": from,
				}
				updateSet := bson.M{}
				if name != "" {
					updateSet["name"] = name
				}
				if identity != "" {
					updateSet["identity"] = identity
				}
				if website != "" {
					updateSet["website"] = website
				}
				if description != "" {
					updateSet["description"] = description
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
				from := dTxs[tIndex].Messages[mIndex].Data["from"].(string)
				node := dTxs[tIndex].Messages[mIndex].Data["node"].(string)

				subscription, err := strconv.ParseUint(dTxs[tIndex].Messages[mIndex].Data["id"].(string), 10, 64)
				if err != nil {
					return nil, err
				}

				eStartSession, err := txResultLog[mIndex].Events.Get("sentinel.session.v1.EventStartSession")
				if err != nil {
					return nil, err
				}

				id, err := strconv.ParseUint(eStartSession.Attributes["id"], 10, 64)
				if err != nil {
					return nil, err
				}

				dSession := types.Session{
					ID:              id,
					Subscription:    subscription,
					Address:         from,
					Node:            node,
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
				id, err := strconv.ParseUint(dTxs[tIndex].Messages[mIndex].Data["proof"].(bson.M)["id"].(string), 10, 64)
				if err != nil {
					return nil, err
				}

				duration, err := time.ParseDuration(dTxs[tIndex].Messages[mIndex].Data["proof"].(bson.M)["duration"].(string))
				if err != nil {
					return nil, err
				}

				buf, err := json.Marshal(dTxs[tIndex].Messages[mIndex].Data["proof"].(bson.M)["bandwidth"])
				if err != nil {
					return nil, err
				}

				var bandwidth hubtypes.Bandwidth
				if err := json.Unmarshal(buf, &bandwidth); err != nil {
					return nil, err
				}

				filter := bson.M{
					"id": id,
				}
				update := bson.M{
					"$set": bson.M{
						"duration":  duration.Nanoseconds(),
						"bandwidth": types.NewBandwidth(&bandwidth),
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
				id, err := strconv.ParseUint(dTxs[tIndex].Messages[mIndex].Data["id"].(string), 10, 64)
				if err != nil {
					return nil, err
				}

				rating, err := strconv.ParseUint(dTxs[tIndex].Messages[mIndex].Data["rating"].(string), 10, 64)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"id": id,
				}
				update := bson.M{
					"$set": bson.M{
						"rating":           rating,
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
				eSubscribeToNode, err := txResultLog[mIndex].Events.Get("sentinel.subscription.v1.EventSubscribeToNode")
				if err != nil {
					return nil, err
				}

				id, err := strconv.ParseUint(eSubscribeToNode.Attributes["id"], 10, 64)
				if err != nil {
					return nil, err
				}

				qSubscription, err := q.QuerySubscription(id, dBlock.Height)
				if err != nil {
					return nil, err
				}

				dSubscription := types.Subscription{
					ID:              qSubscription.Id,
					Owner:           qSubscription.Owner,
					Node:            qSubscription.Node,
					Price:           types.NewCoin(&qSubscription.Price),
					Deposit:         types.NewCoin(&qSubscription.Deposit),
					Plan:            qSubscription.Plan,
					Denom:           qSubscription.Denom,
					Expiry:          qSubscription.Expiry,
					Payment:         nil,
					Free:            qSubscription.Free.Int64(),
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
				eSubscribeToPlan, err := txResultLog[mIndex].Events.Get("sentinel.subscription.v1.EventSubscribeToPlan")
				if err != nil {
					return nil, err
				}

				id, err := strconv.ParseUint(eSubscribeToPlan.Attributes["id"], 10, 64)
				if err != nil {
					return nil, err
				}

				qSubscription, err := q.QuerySubscription(id, dBlock.Height)
				if err != nil {
					return nil, err
				}

				dSubscription := types.Subscription{
					ID:              qSubscription.Id,
					Owner:           qSubscription.Owner,
					Node:            qSubscription.Node,
					Price:           types.NewCoin(&qSubscription.Price),
					Deposit:         types.NewCoin(&qSubscription.Deposit),
					Plan:            qSubscription.Plan,
					Denom:           qSubscription.Denom,
					Expiry:          qSubscription.Expiry,
					Payment:         nil,
					Free:            qSubscription.Free.Int64(),
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
				id, err := strconv.ParseUint(dTxs[tIndex].Messages[mIndex].Data["id"].(string), 10, 64)
				if err != nil {
					return nil, err
				}

				filter := bson.M{
					"id": id,
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

	return operations, nil
}

func main() {
	encCfg := hub.MakeEncodingConfig()

	q, err := querier.NewQuerier(&encCfg, rpcAddress, "/websocket")
	if err != nil {
		log.Fatalln(err)
	}

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
		dSyncStatus = &types.SyncStatus{
			AppName:   appName,
			Height:    fromHeight - 1,
			Timestamp: time.Time{},
		}
	}

	height := dSyncStatus.Height + 1
	for height < toHeight {
		now := time.Now()
		log.Println("Height", height)

		operations, err := run(db, q, height)
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
