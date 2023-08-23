package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/flock"
	hubapp "github.com/sentinel-official/hub/app"
	hubtypes "github.com/sentinel-official/hub/types"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/querier"
	"github.com/sentinel-official/explorer/types"
	commontypes "github.com/sentinel-official/explorer/types/common"
	nodemessages "github.com/sentinel-official/explorer/types/messages/node"
	planmessages "github.com/sentinel-official/explorer/types/messages/plan"
	providermessages "github.com/sentinel-official/explorer/types/messages/provider"
	sessionmessages "github.com/sentinel-official/explorer/types/messages/session"
	subscriptionmessages "github.com/sentinel-official/explorer/types/messages/subscription"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

const (
	appName = "02-events"
)

var (
	height     int64
	toHeight   int64
	rpcAddress string
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.Int64Var(&height, "from-height", 12_310_005, "")
	flag.Int64Var(&toHeight, "to-height", math.MaxInt64, "")
	flag.StringVar(&rpcAddress, "rpc-address", "http://10.104.0.6:26657", "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func getStringEvent(events types.StringEvents, s string) (*types.StringEvent, error) {
	for _, event := range events {
		if event.Type == s {
			return event, nil
		}
	}

	return nil, fmt.Errorf("event %s does not exist", s)
}

func main() {
	encCfg := hubapp.DefaultEncodingConfig()

	q, err := querier.NewQuerier(&encCfg, rpcAddress, "/websocket")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.PrepareDatabase(context.Background(), appName, dbAddress, dbUsername, dbPassword, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	filter := bson.M{
		"app_name": appName,
	}

	dSyncStatus, err := database.SyncStatusFindOne(context.Background(), db, filter)
	if err != nil {
		log.Fatalln(err)
	}
	if dSyncStatus == nil {
		dSyncStatus = &types.SyncStatus{
			AppName:   appName,
			Height:    height - 1,
			Timestamp: time.Time{},
		}
	}

	height = dSyncStatus.Height + 1

	collNodes := flock.New(filepath.Join(os.TempDir(), "mongodb.sentinelhub-2.nodes.lock"))

	for height < toHeight {
		if err := collNodes.Lock(); err != nil {
			log.Fatalln(err)
		}

		log.Println("Height", height)

		now := time.Now()
		err = db.Client().UseSession(
			context.Background(),
			func(sctx mongo.SessionContext) error {
				err = sctx.StartTransaction(
					options.Transaction().
						SetReadConcern(readconcern.Snapshot()).
						SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
				)
				if err != nil {
					return err
				}

				abort := true
				defer func() {
					if abort {
						_ = sctx.AbortTransaction(sctx)
					}
				}()

				filter = bson.M{
					"height": height,
				}
				projection := bson.M{
					"height":           1,
					"time":             1,
					"num_txs":          1,
					"end_block_events": 1,
				}

				dBlock, err := database.BlockFindOne(sctx, db, filter, options.FindOne().SetProjection(projection))
				if err != nil {
					return err
				}
				if dBlock == nil {
					log.Println("Sleeping...")
					time.Sleep(5 * time.Second)
					return nil
				}

				filter = bson.M{
					"height":      height,
					"result.code": 0,
				}
				projection = bson.M{
					"hash":     1,
					"messages": 1,
				}

				dTxs, err := database.TxFindAll(sctx, db, filter)
				if err != nil {
					return err
				}

				log.Println("TxsLen", dBlock.Height, len(dTxs))
				for tIndex := 0; tIndex < len(dTxs); tIndex++ {
					log.Println("MessagesLen", dBlock.Height, tIndex, len(dTxs[tIndex].Messages))
					resultLogs := types.NewTxResultABCIMessageLogsFromRaw(dTxs[tIndex].Result.Logs)
					for mIndex := 0; mIndex < len(dTxs[tIndex].Messages); mIndex++ {
						log.Println("MessageType", dBlock.Height, tIndex, mIndex, dTxs[tIndex].Messages[mIndex].Type)
						buf, err := json.Marshal(dTxs[tIndex].Messages[mIndex].Data)
						if err != nil {
							return err
						}

						log.Println(string(buf))

						switch dTxs[tIndex].Messages[mIndex].Type {
						case "/sentinel.node.v2.MsgRegisterRequest", "/sentinel.node.v2.MsgService/MsgRegister":
							var msg nodemessages.MsgRegisterRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							dNode := &types.Node{
								Address:                explorerutils.MustHexFromBech32AccAddress(msg.From),
								GigabytePrices:         msg.GigabytePrices,
								HourlyPrices:           msg.HourlyPrices,
								RemoteURL:              msg.RemoteURL,
								JoinHeight:             dBlock.Height,
								JoinTimestamp:          dBlock.Time,
								JoinTxHash:             dTxs[tIndex].Hash,
								Bandwidth:              nil,
								Handshake:              nil,
								IntervalSetSessions:    0,
								IntervalUpdateSessions: 0,
								IntervalUpdateStatus:   0,
								Location:               nil,
								Moniker:                "",
								Peers:                  0,
								QOS:                    nil,
								Type:                   0,
								Version:                "",
								Status:                 hubtypes.StatusInactive.String(),
								StatusHeight:           dBlock.Height,
								StatusTimestamp:        dBlock.Time,
								StatusTxHash:           dTxs[tIndex].Hash,
								ReachStatus:            nil,
							}

							if err := database.NodeSave(sctx, db, dNode); err != nil {
								return err
							}
						case "/sentinel.node.v2.MsgUpdateDetailsRequest", "/sentinel.node.v2.MsgService/MsgUpdateDetails":
							var msg nodemessages.MsgUpdateDetailsRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"address": explorerutils.MustHexFromBech32NodeAddress(msg.From),
							}

							updateSet := bson.M{}
							if msg.GigabytePrices != nil && len(msg.GigabytePrices) > 0 {
								updateSet["gigabyte_prices"] = msg.GigabytePrices
							}
							if msg.HourlyPrices != nil && len(msg.HourlyPrices) > 0 {
								updateSet["hourly_prices"] = msg.HourlyPrices
							}
							if msg.RemoteURL != "" {
								updateSet["remote_url"] = msg.RemoteURL
							}

							update := bson.M{
								"$set": updateSet,
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.NodeFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.node.v2.MsgUpdateStatusRequest", "/sentinel.node.v2.MsgService/MsgUpdateStatus":
							var msg nodemessages.MsgUpdateStatusRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							nodeAddr := explorerutils.MustHexFromBech32NodeAddress(msg.From)

							filter = bson.M{
								"address": nodeAddr,
							}
							update := bson.M{
								"$set": bson.M{
									"status":           msg.Status,
									"status_height":    dBlock.Height,
									"status_timestamp": dBlock.Time,
									"status_tx_hash":   dTxs[tIndex].Hash,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.NodeFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							dNodeEvent := &types.NodeEvent{
								Address:   nodeAddr,
								Status:    msg.Status,
								Height:    dBlock.Height,
								Timestamp: dBlock.Time,
								TxHash:    dTxs[tIndex].Hash,
							}

							if err := database.NodeEventSave(sctx, db, dNodeEvent); err != nil {
								return err
							}
						case "/sentinel.node.v2.MsgSubscribeRequest", "/sentinel.node.v2.MsgService/MsgSubscribe":
							eCreateSubscription, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.node.v2.EventCreateSubscription")
							if err != nil {
								return err
							}

							id, err := strconv.ParseUint(eCreateSubscription.Attributes["id"], 10, 64)
							if err != nil {
								return err
							}

							qSubscriptionI, err := q.QuerySubscription(id, dBlock.Height)
							if err != nil {
								return err
							}

							qSubscription, ok := qSubscriptionI.(*subscriptiontypes.NodeSubscription)
							if !ok {
								return fmt.Errorf("invalid subscription type %s", qSubscriptionI.Type())
							}

							hexAccAddr := explorerutils.MustHexFromBech32AccAddress(qSubscription.Address)
							hexNodeAddr := explorerutils.MustHexFromBech32NodeAddress(qSubscription.NodeAddress)

							dSubscription := &types.Subscription{
								ID:              id,
								Address:         hexAccAddr,
								InactiveAt:      qSubscription.InactiveAt,
								NodeAddress:     hexNodeAddr,
								Gigabytes:       qSubscription.Gigabytes,
								Hours:           qSubscription.Hours,
								Deposit:         commontypes.NewCoinFromRaw(&qSubscription.Deposit),
								PlanID:          0,
								Denom:           "",
								StakingReward:   nil,
								Payment:         nil,
								StartHeight:     dBlock.Height,
								StartTimestamp:  dBlock.Time,
								StartTxHash:     dTxs[tIndex].Hash,
								EndHeight:       0,
								EndTimestamp:    time.Time{},
								EndTxHash:       "",
								Status:          qSubscription.Status.String(),
								StatusHeight:    dBlock.Height,
								StatusTimestamp: dBlock.Time,
								StatusTxHash:    dTxs[tIndex].Hash,
							}

							if err = database.SubscriptionSave(sctx, db, dSubscription); err != nil {
								return err
							}

							accAddr, err := sdk.AccAddressFromHex(hexAccAddr)
							if err != nil {
								return err
							}

							qDeposit, err := q.QueryDeposit(accAddr, dBlock.Height)
							if err != nil {
								return err
							}

							filter = bson.M{
								"address": hexAccAddr,
							}
							update := bson.M{
								"$set": bson.M{
									"coins":     commontypes.NewCoinsFromRaw(qDeposit.Coins),
									"height":    dBlock.Height,
									"timestamp": dBlock.Time,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.DepositFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							dDepositEvent := &types.DepositEvent{
								Address:   hexAccAddr,
								Coins:     commontypes.NewCoinsFromRaw(sdk.NewCoins(qSubscription.Deposit)),
								Action:    "addition",
								Height:    dBlock.Height,
								Timestamp: dBlock.Time,
								TxHash:    dTxs[tIndex].Hash,
							}

							if err = database.DepositEventSave(sctx, db, dDepositEvent); err != nil {
								return err
							}

							if qSubscription.Gigabytes != 0 {
								qAllocation, err := q.QueryAllocation(qSubscription.ID, accAddr, dBlock.Height)
								if err != nil {
									return err
								}

								dAllocation := &types.Allocation{
									ID:            qSubscription.ID,
									Address:       hexAccAddr,
									UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
									GrantedBytes:  qAllocation.GrantedBytes.Int64(),
								}

								if err := database.AllocationSave(sctx, db, dAllocation); err != nil {
									return err
								}

								dAllocationEvent := &types.AllocationEvent{
									ID:            qSubscription.ID,
									Address:       hexAccAddr,
									UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
									GrantedBytes:  qAllocation.GrantedBytes.Int64(),
									Height:        dBlock.Height,
									Timestamp:     dBlock.Time,
									TxHash:        dTxs[tIndex].Hash,
								}

								if err := database.AllocationEventSave(sctx, db, dAllocationEvent); err != nil {
									return err
								}
							}
						case "/sentinel.plan.v2.MsgCreateRequest", "/sentinel.plan.v2.MsgService/MsgCreate":
							var msg planmessages.MsgCreateRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							eCreate, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.plan.v2.EventCreate")
							if err != nil {
								return err
							}

							id, err := strconv.ParseUint(eCreate.Attributes["id"], 10, 64)
							if err != nil {
								return err
							}

							dPlan := &types.Plan{
								ID:              id,
								ProviderAddress: explorerutils.MustHexFromBech32ProvAddress(msg.From),
								Prices:          msg.Prices,
								Duration:        msg.Duration.Nanoseconds(),
								Gigabytes:       msg.Gigabytes,
								NodeAddresses:   []string{},
								CreateHeight:    dBlock.Height,
								CreateTimestamp: dBlock.Time,
								CreateTxHash:    dTxs[tIndex].Hash,
								Status:          hubtypes.StatusInactive.String(),
								StatusHeight:    dBlock.Height,
								StatusTimestamp: dBlock.Time,
								StatusTxHash:    dTxs[tIndex].Hash,
							}

							if err := database.PlanSave(sctx, db, dPlan); err != nil {
								return err
							}
						case "/sentinel.plan.v2.MsgUpdateStatusRequest", "/sentinel.plan.v2.MsgService/MsgUpdateStatus":
							var msg planmessages.MsgUpdateStatusRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$set": bson.M{
									"status":           msg.Status,
									"status_height":    dBlock.Height,
									"status_timestamp": dBlock.Time,
									"status_tx_hash":   dTxs[tIndex].Hash,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.PlanFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.plan.v2.MsgLinkNodeRequest", "/sentinel.plan.v2.MsgService/MsgLinkNode":
							var msg planmessages.MsgLinkNodeRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$push": bson.M{
									"node_addresses": explorerutils.MustHexFromBech32NodeAddress(msg.NodeAddress),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.PlanFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.plan.v2.MsgUnlinkNodeRequest", "/sentinel.plan.v2.MsgService/MsgUnlinkNode":
							var msg planmessages.MsgUnlinkNodeRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$pull": bson.M{
									"node_addresses": explorerutils.MustHexFromBech32NodeAddress(msg.NodeAddress),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.PlanFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.plan.v2.MsgSubscribeRequest", "/sentinel.plan.v2.MsgService/MsgSubscribe":
							eCreateSubscription, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.plan.v2.EventCreateSubscription")
							if err != nil {
								return err
							}

							id, err := strconv.ParseUint(eCreateSubscription.Attributes["id"], 10, 64)
							if err != nil {
								return err
							}

							qSubscriptionI, err := q.QuerySubscription(id, dBlock.Height)
							if err != nil {
								return err
							}

							qSubscription, ok := qSubscriptionI.(*subscriptiontypes.PlanSubscription)
							if !ok {
								return fmt.Errorf("invalid subscription type %s", qSubscriptionI.Type())
							}

							hexAccAddr := explorerutils.MustHexFromBech32AccAddress(qSubscription.Address)

							dSubscription := &types.Subscription{
								ID:              qSubscription.ID,
								Address:         hexAccAddr,
								InactiveAt:      qSubscription.InactiveAt,
								NodeAddress:     "",
								Gigabytes:       0,
								Hours:           0,
								Deposit:         nil,
								PlanID:          qSubscription.PlanID,
								Denom:           qSubscription.Denom,
								StakingReward:   nil,
								Payment:         nil,
								StartHeight:     dBlock.Height,
								StartTimestamp:  dBlock.Time,
								StartTxHash:     dTxs[tIndex].Hash,
								EndHeight:       0,
								EndTimestamp:    time.Time{},
								EndTxHash:       "",
								Status:          qSubscription.Status.String(),
								StatusHeight:    dBlock.Height,
								StatusTimestamp: dBlock.Time,
								StatusTxHash:    dTxs[tIndex].Hash,
							}

							ePayForPlan, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.subscription.v2.EventPayForPlan")
							if err != nil {
								return err
							}

							payment, err := sdk.ParseCoinNormalized(ePayForPlan.Attributes["payment"])
							if err != nil {
								return err
							}

							stakingReward, err := sdk.ParseCoinNormalized(ePayForPlan.Attributes["staking_reward"])
							if err != nil {
								return err
							}

							dSubscription.StakingReward = commontypes.NewCoinFromRaw(&stakingReward)
							dSubscription.Payment = commontypes.NewCoinFromRaw(&payment)

							if err := database.SubscriptionSave(sctx, db, dSubscription); err != nil {
								return err
							}

							accAddr, err := sdk.AccAddressFromHex(hexAccAddr)
							if err != nil {
								return err
							}

							qAllocation, err := q.QueryAllocation(qSubscription.ID, accAddr, dBlock.Height)
							if err != nil {
								return err
							}

							dAllocation := &types.Allocation{
								ID:            qSubscription.ID,
								Address:       hexAccAddr,
								UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
								GrantedBytes:  qAllocation.GrantedBytes.Int64(),
							}

							if err := database.AllocationSave(sctx, db, dAllocation); err != nil {
								return err
							}

							dAllocationEvent := &types.AllocationEvent{
								ID:            qSubscription.ID,
								Address:       hexAccAddr,
								UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
								GrantedBytes:  qAllocation.GrantedBytes.Int64(),
								Height:        dBlock.Height,
								Timestamp:     dBlock.Time,
								TxHash:        dTxs[tIndex].Hash,
							}

							if err := database.AllocationEventSave(sctx, db, dAllocationEvent); err != nil {
								return err
							}
						case "/sentinel.provider.v2.MsgRegisterRequest", "/sentinel.provider.v2.MsgService/MsgRegister":
							var msg providermessages.MsgRegisterRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							dProvider := &types.Provider{
								Address:         explorerutils.MustHexFromBech32AccAddress(msg.From),
								Name:            msg.Name,
								Identity:        msg.Identity,
								Website:         msg.Website,
								Description:     msg.Description,
								Status:          hubtypes.StatusInactive.String(),
								StatusHeight:    dBlock.Height,
								StatusTimestamp: dBlock.Time,
								StatusTxHash:    dTxs[tIndex].Hash,
								JoinHeight:      dBlock.Height,
								JoinTimestamp:   dBlock.Time,
								JoinTxHash:      dTxs[tIndex].Hash,
							}

							if err := database.ProviderSave(sctx, db, dProvider); err != nil {
								return err
							}
						case "/sentinel.provider.v2.MsgUpdateRequest", "/sentinel.provider.v2.MsgService/MsgUpdate":
							var msg providermessages.MsgUpdateRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							provAddr, err := hubtypes.ProvAddressFromBech32(msg.From)
							if err != nil {
								return err
							}

							qProvider, err := q.QueryProvider(provAddr, dBlock.Height)
							if err != nil {
								return err
							}

							filter = bson.M{
								"address": explorerutils.MustHexFromBech32ProvAddress(msg.From),
							}

							updateSet := bson.M{
								"name":        qProvider.Name,
								"identity":    qProvider.Identity,
								"website":     qProvider.Website,
								"description": qProvider.Description,
								"status":      qProvider.Status,
							}

							if msg.Status != hubtypes.StatusUnspecified.String() {
								updateSet["status_height"] = dBlock.Height
								updateSet["status_timestamp"] = dBlock.Time
								updateSet["status_tx_hash"] = dTxs[tIndex].Hash
							}

							update := bson.M{
								"$set": updateSet,
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.ProviderFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.subscription.v2.MsgCancelRequest", "/sentinel.subscription.v2.MsgService/MsgCancel":
							var msg subscriptionmessages.MsgCancelRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$set": bson.M{
									"end_height":       dBlock.Height,
									"end_timestamp":    dBlock.Time,
									"end_tx_hash":      dTxs[tIndex].Hash,
									"status":           hubtypes.StatusInactivePending.String(),
									"status_height":    dBlock.Height,
									"status_timestamp": dBlock.Time,
									"status_tx_hash":   dTxs[tIndex].Hash,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.SubscriptionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							filter = bson.M{
								"subscription_id": msg.ID,
							}

							_, err = database.SessionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.subscription.v2.MsgAllocateRequest", "/sentinel.subscription.v2.MsgService/MsgAllocate":
							var msg subscriptionmessages.MsgAllocateRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							accAddr, err := sdk.AccAddressFromBech32(msg.From)
							if err != nil {
								return err
							}

							qAllocation, err := q.QueryAllocation(msg.ID, accAddr, dBlock.Height)
							if err != nil {
								return err
							}

							hexAccAddr := explorerutils.MustHexFromBech32AccAddress(msg.From)

							filter = bson.M{
								"id":      msg.ID,
								"address": hexAccAddr,
							}
							update := bson.M{
								"$set": bson.M{
									"granted_bytes":  qAllocation.GrantedBytes.Int64(),
									"utilised_bytes": qAllocation.UtilisedBytes.Int64(),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.AllocationFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							dAllocationEvent := &types.AllocationEvent{
								ID:            msg.ID,
								Address:       hexAccAddr,
								GrantedBytes:  qAllocation.GrantedBytes.Int64(),
								UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
								Height:        dBlock.Height,
								Timestamp:     dBlock.Time,
								TxHash:        dTxs[tIndex].Hash,
							}

							if err := database.AllocationEventSave(sctx, db, dAllocationEvent); err != nil {
								return err
							}

							accAddr, err = sdk.AccAddressFromBech32(msg.Address)
							if err != nil {
								return err
							}

							qAllocation, err = q.QueryAllocation(msg.ID, accAddr, dBlock.Height)
							if err != nil {
								return err
							}

							hexAccAddr = explorerutils.MustHexFromBech32AccAddress(msg.Address)

							filter = bson.M{
								"id":      msg.ID,
								"address": hexAccAddr,
							}
							update = bson.M{
								"$set": bson.M{
									"granted_bytes":  qAllocation.GrantedBytes.Int64(),
									"utilised_bytes": qAllocation.UtilisedBytes.Int64(),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.AllocationFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							dAllocationEvent = &types.AllocationEvent{
								ID:            msg.ID,
								Address:       hexAccAddr,
								GrantedBytes:  qAllocation.GrantedBytes.Int64(),
								UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
								Height:        dBlock.Height,
								Timestamp:     dBlock.Time,
								TxHash:        dTxs[tIndex].Hash,
							}

							if err := database.AllocationEventSave(sctx, db, dAllocationEvent); err != nil {
								return err
							}
						case "/sentinel.session.v2.MsgStartRequest", "/sentinel.session.v2.MsgService/MsgStart":
							var msg sessionmessages.MsgStartRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							eStart, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.session.v2.EventStart")
							if err != nil {
								return err
							}

							id, err := strconv.ParseUint(eStart.Attributes["id"], 10, 64)
							if err != nil {
								return err
							}

							dSession := &types.Session{
								ID:              id,
								SubscriptionID:  msg.ID,
								Address:         explorerutils.MustHexFromBech32AccAddress(msg.From),
								NodeAddress:     explorerutils.MustHexFromBech32NodeAddress(msg.Address),
								Duration:        0,
								Bandwidth:       nil,
								StartHeight:     dBlock.Height,
								StartTimestamp:  dBlock.Time,
								StartTxHash:     dTxs[tIndex].Hash,
								EndHeight:       0,
								EndTimestamp:    time.Time{},
								EndTxHash:       "",
								StakingReward:   nil,
								Payment:         nil,
								Rating:          0,
								Status:          hubtypes.StatusActive.String(),
								StatusHeight:    dBlock.Height,
								StatusTimestamp: dBlock.Time,
								StatusTxHash:    dTxs[tIndex].Hash,
							}

							if err := database.SessionSave(sctx, db, dSession); err != nil {
								return err
							}

							dSessionEvent := &types.SessionEvent{
								ID:        dSession.ID,
								Bandwidth: dSession.Bandwidth,
								Duration:  dSession.Duration,
								Signature: "",
								Height:    dBlock.Height,
								Timestamp: dBlock.Time,
								TxHash:    dTxs[tIndex].Hash,
							}

							if err := database.SessionEventSave(sctx, db, dSessionEvent); err != nil {
								return err
							}
						case "/sentinel.session.v2.MsgUpdateDetailsRequest", "/sentinel.session.v2.MsgService/MsgUpdateDetails":
							var msg sessionmessages.MsgUpdateDetailsRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$set": bson.M{
									"duration":  msg.Duration.Nanoseconds(),
									"bandwidth": msg.Bandwidth,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.SessionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							dSessionEvent := &types.SessionEvent{
								ID:        msg.ID,
								Bandwidth: msg.Bandwidth,
								Duration:  msg.Duration.Nanoseconds(),
								Signature: msg.Signature,
								Height:    dBlock.Height,
								Timestamp: dBlock.Time,
								TxHash:    dTxs[tIndex].Hash,
							}

							if err := database.SessionEventSave(sctx, db, dSessionEvent); err != nil {
								return err
							}
						case "/sentinel.session.v2.MsgEndRequest", "/sentinel.session.v2.MsgService/MsgEnd":
							var msg sessionmessages.MsgEndRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$set": bson.M{
									"end_height":       dBlock.Height,
									"end_timestamp":    dBlock.Time,
									"end_tx_hash":      dTxs[tIndex].Hash,
									"rating":           msg.Rating,
									"status":           hubtypes.StatusInactivePending.String(),
									"status_height":    dBlock.Height,
									"status_timestamp": dBlock.Time,
									"status_tx_hash":   dTxs[tIndex].Hash,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.SessionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						default:

						}
					}
				}

				log.Println("EndBlockEventsLen", dBlock.Height, len(dBlock.EndBlockEvents))
				for eIndex := 0; eIndex < len(dBlock.EndBlockEvents); eIndex++ {
					log.Println("EndBlockEventType", eIndex, dBlock.EndBlockEvents[eIndex].Type)
					switch dBlock.EndBlockEvents[eIndex].Type {
					case "sentinel.node.v2.EventUpdateDetails":
						nodeAddr, err := hubtypes.NodeAddressFromBech32(dBlock.EndBlockEvents[eIndex].Attributes["address"])
						if err != nil {
							return err
						}

						qNode, err := q.QueryNode(nodeAddr, dBlock.Height)
						if err != nil {
							return err
						}

						filter = bson.M{
							"address": explorerutils.MustHexFromBech32NodeAddress(qNode.Address),
						}
						update := bson.M{
							"$set": bson.M{
								"gigabyte_prices": commontypes.NewCoinsFromRaw(qNode.GigabytePrices),
								"hourly_prices":   commontypes.NewCoinsFromRaw(qNode.HourlyPrices),
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.NodeFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}
					case "sentinel.node.v2.EventUpdateStatus":
						hexNodeAddr := explorerutils.MustHexFromBech32NodeAddress(dBlock.EndBlockEvents[eIndex].Attributes["address"])

						filter = bson.M{
							"address": hexNodeAddr,
						}
						update := bson.M{
							"$set": bson.M{
								"status":           dBlock.EndBlockEvents[eIndex].Attributes["status"],
								"status_height":    dBlock.Height,
								"status_timestamp": dBlock.Time,
								"status_tx_hash":   "",
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.NodeFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}

						dNodeEvent := &types.NodeEvent{
							Address:   hexNodeAddr,
							Status:    dBlock.EndBlockEvents[eIndex].Attributes["status"],
							Height:    dBlock.Height,
							Timestamp: dBlock.Time,
						}

						if err := database.NodeEventSave(sctx, db, dNodeEvent); err != nil {
							return err
						}
					case "sentinel.subscription.v2.EventPayForSession":
						payment, err := sdk.ParseCoinNormalized(dBlock.EndBlockEvents[eIndex].Attributes["payment"])
						if err != nil {
							return err
						}

						stakingReward, err := sdk.ParseCoinNormalized(dBlock.EndBlockEvents[eIndex].Attributes["staking_reward"])
						if err != nil {
							return err
						}

						id, err := strconv.ParseInt(dBlock.EndBlockEvents[eIndex].Attributes["session_id"], 10, 64)
						if err != nil {
							return err
						}

						filter = bson.M{
							"id": id,
						}

						update := bson.M{
							"$set": bson.M{
								"payment":        commontypes.NewCoinFromRaw(&payment),
								"staking_reward": commontypes.NewCoinFromRaw(&stakingReward),
							},
						}
						projection = bson.M{
							"address": 1,
						}

						dSession, err := database.SessionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}
						if dSession == nil || dSession.Address == "" {
							continue
						}

						accAddr, err := sdk.AccAddressFromHex(dSession.Address)
						if err != nil {
							return err
						}

						qDeposit, err := q.QueryDeposit(accAddr, dBlock.Height)
						if err != nil {
							return err
						}

						filter = bson.M{
							"address": dSession.Address,
						}
						update = bson.M{
							"$set": bson.M{
								"coins":     commontypes.NewCoinsFromRaw(qDeposit.Coins),
								"height":    dBlock.Height,
								"timestamp": dBlock.Time,
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.DepositFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}

						coins := sdk.NewCoins(stakingReward.Add(payment))

						dDepositEvent := &types.DepositEvent{
							Address:   dSession.Address,
							Coins:     commontypes.NewCoinsFromRaw(coins),
							Action:    "subtract",
							Height:    dBlock.Height,
							Timestamp: dBlock.Time,
							TxHash:    "",
						}

						if err := database.DepositEventSave(sctx, db, dDepositEvent); err != nil {
							return err
						}
					case "sentinel.deposit.v1.EventAdd":
					case "sentinel.deposit.v1.EventSubtract":
						accAddr, err := sdk.AccAddressFromBech32(dBlock.EndBlockEvents[eIndex].Attributes["address"])
						if err != nil {
							return err
						}

						qDeposit, err := q.QueryDeposit(accAddr, dBlock.Height)
						if err != nil {
							return err
						}

						hexAccAddr := explorerutils.MustHexFromBech32AccAddress(qDeposit.Address)

						filter = bson.M{
							"address": hexAccAddr,
						}
						update := bson.M{
							"$set": bson.M{
								"coins":     commontypes.NewCoinsFromRaw(qDeposit.Coins),
								"height":    dBlock.Height,
								"timestamp": dBlock.Time,
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.DepositFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}

						rawCoins, err := sdk.ParseCoinsNormalized(dBlock.EndBlockEvents[eIndex].Attributes["coins"])
						if err != nil {
							return err
						}

						coins := commontypes.NewCoinsFromRaw(rawCoins)
						dDepositEvent := &types.DepositEvent{
							Address:   hexAccAddr,
							Coins:     coins,
							Action:    "subtract",
							Height:    dBlock.Height,
							Timestamp: dBlock.Time,
							TxHash:    "",
						}

						if err := database.DepositEventSave(sctx, db, dDepositEvent); err != nil {
							return err
						}
					case "sentinel.session.v2.EventUpdateStatus":
						id, err := strconv.ParseUint(dBlock.EndBlockEvents[eIndex].Attributes["id"], 10, 64)
						if err != nil {
							return err
						}

						filter = bson.M{
							"id": id,
						}

						updateSet := bson.M{
							"status":           dBlock.EndBlockEvents[eIndex].Attributes["status"],
							"status_height":    dBlock.Height,
							"status_timestamp": dBlock.Time,
							"status_tx_hash":   "",
						}

						if dBlock.EndBlockEvents[eIndex].Attributes["status"] == hubtypes.StatusInactivePending.String() {
							updateSet["end_height"] = dBlock.Height
							updateSet["end_timestamp"] = dBlock.Time
							updateSet["end_tx_hash"] = ""
							updateSet["rating"] = 0
						}

						update := bson.M{
							"$set": updateSet,
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.SessionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}
					case "sentinel.subscription.v2.EventAllocate":
						id, err := strconv.ParseUint(dBlock.EndBlockEvents[eIndex].Attributes["id"], 10, 64)
						if err != nil {
							return err
						}

						accAddr, err := sdk.AccAddressFromBech32(dBlock.EndBlockEvents[eIndex].Attributes["address"])
						if err != nil {
							return err
						}

						hexAccAddr := explorerutils.MustHexFromBech32AccAddress(accAddr.String())

						qAllocation, err := q.QueryAllocation(id, accAddr, dBlock.Height)
						if err != nil {
							return err
						}

						filter = bson.M{
							"id":      id,
							"address": hexAccAddr,
						}
						update := bson.M{
							"$set": bson.M{
								"granted_bytes":  qAllocation.GrantedBytes.Int64(),
								"utilised_bytes": qAllocation.UtilisedBytes.Int64(),
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.AllocationFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}

						dAllocationEvent := &types.AllocationEvent{
							ID:            id,
							Address:       hexAccAddr,
							GrantedBytes:  qAllocation.GrantedBytes.Int64(),
							UtilisedBytes: qAllocation.UtilisedBytes.Int64(),
							Height:        dBlock.Height,
							Timestamp:     dBlock.Time,
							TxHash:        "",
						}

						if err := database.AllocationEventSave(sctx, db, dAllocationEvent); err != nil {
							return err
						}
					case "sentinel.subscription.v2.EventUpdateStatus":
						id, err := strconv.ParseUint(dBlock.EndBlockEvents[eIndex].Attributes["id"], 10, 64)
						if err != nil {
							return err
						}

						filter = bson.M{
							"id": id,
						}

						updateSet := bson.M{
							"status":           dBlock.EndBlockEvents[eIndex].Attributes["status"],
							"status_height":    dBlock.Height,
							"status_timestamp": dBlock.Time,
							"status_tx_hash":   "",
						}

						if dBlock.EndBlockEvents[eIndex].Attributes["status"] == hubtypes.StatusInactivePending.String() {
							updateSet["end_height"] = dBlock.Height
							updateSet["end_timestamp"] = dBlock.Time
							updateSet["end_tx_hash"] = ""
						}

						update := bson.M{
							"$set": updateSet,
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.SubscriptionFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}
					default:
					}
				}

				filter = bson.M{
					"app_name": appName,
				}
				update := bson.M{
					"$set": bson.M{
						"height": height,
					},
				}
				projection = bson.M{
					"_id": 1,
				}

				_, err = database.SyncStatusFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
				if err != nil {
					return err
				}

				height++

				abort = false
				return sctx.CommitTransaction(sctx)
			},
		)
		log.Println("Duration", time.Since(now))
		if err != nil {
			log.Fatalln(err)
		}

		if err := collNodes.Unlock(); err != nil {
			log.Fatalln(err)
		}
	}
}
