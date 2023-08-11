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
	hubtypes "github.com/sentinel-official/hub/types"
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
	flag.Int64Var(&height, "from-height", 9_348_475, "")
	flag.Int64Var(&toHeight, "to-height", math.MaxInt64, "")
	// flag.StringVar(&rpcAddress, "rpc-address", "http://188.34.144.2:26657", "")
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
	q, err := querier.NewQuerier(rpcAddress, "/websocket")
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
						case "/sentinel.node.v1.MsgRegisterRequest", "/sentinel.node.v1.MsgService/MsgRegister":
							var msg nodemessages.MsgRegisterRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							dNode := &types.Node{
								Address:                explorerutils.MustHexFromBech32AccAddress(msg.From),
								ProviderAddress:        explorerutils.MustHexFromBech32ProvAddress(msg.Provider),
								Price:                  msg.Price,
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
						case "/sentinel.node.v1.MsgUpdateRequest", "/sentinel.node.v1.MsgService/MsgUpdate":
							var msg nodemessages.MsgUpdateRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"address": explorerutils.MustHexFromBech32NodeAddress(msg.From),
							}

							updateSet := bson.M{}
							if msg.Provider != "" {
								updateSet["provider_address"] = explorerutils.MustHexFromBech32ProvAddress(msg.Provider)
								updateSet["price"] = nil
							}
							if msg.Price != nil && len(msg.Price) > 0 {
								updateSet["provider_address"] = ""
								updateSet["price"] = msg.Price
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
						case "/sentinel.node.v1.MsgSetStatusRequest", "/sentinel.node.v1.MsgService/MsgSetStatus":
							var msg nodemessages.MsgSetStatusRequest
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
						case "/sentinel.subscription.v1.MsgSubscribeToNodeRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToNode",
							"/sentinel.subscription.v1.MsgSubscribeToPlanRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToPlan":
							eSubscribe, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.subscription.v1.EventSubscribe")
							if err != nil {
								return err
							}

							id, err := strconv.ParseUint(eSubscribe.Attributes["id"], 10, 64)
							if err != nil {
								return err
							}

							qSubscription, err := q.QuerySubscription(id, dBlock.Height)
							if err != nil {
								return err
							}

							hexAccAddr := explorerutils.MustHexFromBech32AccAddress(qSubscription.Owner)

							dSubscription := &types.Subscription{
								ID:              qSubscription.Id,
								Address:         hexAccAddr,
								FreeBytes:       qSubscription.Free.Int64(),
								NodeAddress:     explorerutils.MustHexFromBech32NodeAddress(qSubscription.Node),
								Price:           commontypes.NewCoinFromRaw(&qSubscription.Price),
								Deposit:         commontypes.NewCoinFromRaw(&qSubscription.Deposit),
								PlanID:          qSubscription.Plan,
								StakingReward:   nil,
								Payment:         nil,
								Expiry:          qSubscription.Expiry,
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

							var rawCoin sdk.Coin
							if err := json.Unmarshal([]byte(eSubscribe.Attributes["amount"]), &rawCoin); err != nil {
								return err
							}

							amount := sdk.NewCoins(rawCoin)

							if dSubscription.PlanID == 0 {
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
									Coins:     commontypes.NewCoinsFromRaw(amount),
									Subtract:  false,
									Height:    dBlock.Height,
									Timestamp: dBlock.Time,
									TxHash:    dTxs[tIndex].Hash,
								}

								if err := database.DepositEventSave(sctx, db, dDepositEvent); err != nil {
									return err
								}
							} else {
								eStakingReward, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.subscription.v1.EventStakingReward")
								if err != nil {
									return err
								}

								var rawCoin sdk.Coin
								if err := json.Unmarshal([]byte(eStakingReward.Attributes["amount"]), &rawCoin); err != nil {
									return err
								}

								dSubscription.StakingReward = commontypes.NewCoinFromRaw(&rawCoin)
								dSubscription.Payment = commontypes.NewCoinFromRaw(&amount[0])
							}

							if err := database.SubscriptionSave(sctx, db, dSubscription); err != nil {
								return err
							}

							qSubscriptionQuota, err := q.QuerySubscriptionQuota(qSubscription.Id, qSubscription.GetOwner(), dBlock.Height)
							if err != nil {
								return err
							}

							dSubscriptionQuota := &types.SubscriptionQuota{
								ID:             qSubscription.Id,
								Address:        hexAccAddr,
								ConsumedBytes:  qSubscriptionQuota.Consumed.Int64(),
								AllocatedBytes: qSubscriptionQuota.Allocated.Int64(),
							}

							if err := database.SubscriptionQuotaSave(sctx, db, dSubscriptionQuota); err != nil {
								return err
							}

							dSubscriptionQuotaEvent := &types.SubscriptionQuotaEvent{
								ID:             qSubscription.Id,
								Address:        hexAccAddr,
								ConsumedBytes:  qSubscriptionQuota.Consumed.Int64(),
								AllocatedBytes: qSubscriptionQuota.Allocated.Int64(),
								Height:         dBlock.Height,
								Timestamp:      dBlock.Time,
								TxHash:         dTxs[tIndex].Hash,
							}

							if err := database.SubscriptionQuotaEventSave(sctx, db, dSubscriptionQuotaEvent); err != nil {
								return err
							}
						case "/sentinel.subscription.v1.MsgCancelRequest", "/sentinel.subscription.v1.MsgService/MsgCancel":
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
						case "/sentinel.subscription.v1.MsgAddQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgAddQuota":
							var msg subscriptionmessages.MsgAddQuotaRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							hexAccAddr := explorerutils.MustHexFromBech32AccAddress(msg.Address)

							dSubscriptionQuota := &types.SubscriptionQuota{
								ID:             msg.ID,
								Address:        hexAccAddr,
								ConsumedBytes:  0,
								AllocatedBytes: msg.Bytes,
							}

							if err := database.SubscriptionQuotaSave(sctx, db, dSubscriptionQuota); err != nil {
								return err
							}

							dSubscriptionQuotaEvent := &types.SubscriptionQuotaEvent{
								ID:             msg.ID,
								Address:        hexAccAddr,
								ConsumedBytes:  0,
								AllocatedBytes: msg.Bytes,
								Height:         dBlock.Height,
								Timestamp:      dBlock.Time,
								TxHash:         dTxs[tIndex].Hash,
							}

							if err := database.SubscriptionQuotaEventSave(sctx, db, dSubscriptionQuotaEvent); err != nil {
								return err
							}
						case "/sentinel.subscription.v1.MsgUpdateQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgUpdateQuota":
							var msg subscriptionmessages.MsgUpdateQuotaRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							accAddr, err := sdk.AccAddressFromBech32(msg.Address)
							if err != nil {
								return err
							}

							qSubscriptionQuota, err := q.QuerySubscriptionQuota(msg.ID, accAddr, dBlock.Height)
							if err != nil {
								return err
							}

							hexAccAddr := explorerutils.MustHexFromBech32AccAddress(msg.Address)

							filter = bson.M{
								"id":      msg.ID,
								"address": hexAccAddr,
							}
							update := bson.M{
								"$set": bson.M{
									"consumed_bytes":  qSubscriptionQuota.Consumed.Int64(),
									"allocated_bytes": qSubscriptionQuota.Allocated.Int64(),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.SubscriptionQuotaFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}

							dSubscriptionQuotaEvent := &types.SubscriptionQuotaEvent{
								ID:             msg.ID,
								Address:        hexAccAddr,
								ConsumedBytes:  qSubscriptionQuota.Consumed.Int64(),
								AllocatedBytes: qSubscriptionQuota.Allocated.Int64(),
								Height:         dBlock.Height,
								Timestamp:      dBlock.Time,
								TxHash:         dTxs[tIndex].Hash,
							}

							if err := database.SubscriptionQuotaEventSave(sctx, db, dSubscriptionQuotaEvent); err != nil {
								return err
							}
						case "/sentinel.session.v1.MsgStartRequest", "/sentinel.session.v1.MsgService/MsgStart":
							var msg sessionmessages.MsgStartRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							eStart, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.session.v1.EventStart")
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
								NodeAddress:     explorerutils.MustHexFromBech32NodeAddress(msg.Node),
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
						case "/sentinel.session.v1.MsgUpdateRequest", "/sentinel.session.v1.MsgService/MsgUpdate":
							var msg sessionmessages.MsgUpdateRequest
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
						case "/sentinel.session.v1.MsgEndRequest", "/sentinel.session.v1.MsgService/MsgEnd":
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
						case "/sentinel.plan.v1.MsgAddRequest", "/sentinel.plan.v1.MsgService/MsgAdd":
							var msg planmessages.MsgAddRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							eAdd, err := getStringEvent(resultLogs[mIndex].Events, "sentinel.plan.v1.EventAdd")
							if err != nil {
								return err
							}

							id, err := strconv.ParseUint(eAdd.Attributes["id"], 10, 64)
							if err != nil {
								return err
							}

							dPlan := &types.Plan{
								ID:              id,
								ProviderAddress: explorerutils.MustHexFromBech32ProvAddress(msg.From),
								Price:           msg.Price,
								Validity:        msg.Validity.Nanoseconds(),
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

							if err := database.PlanSave(sctx, db, dPlan); err != nil {
								return err
							}
						case "/sentinel.plan.v1.MsgSetStatusRequest", "/sentinel.plan.v1.MsgService/MsgSetStatus":
							var msg planmessages.MsgSetStatusRequest
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
						case "/sentinel.plan.v1.MsgAddNodeRequest", "/sentinel.plan.v1.MsgService/MsgAddNode":
							var msg planmessages.MsgAddNodeRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$push": bson.M{
									"node_addresses": explorerutils.MustHexFromBech32NodeAddress(msg.Address),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.PlanFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.plan.v1.MsgRemoveNodeRequest", "/sentinel.plan.v1.MsgService/MsgRemoveNode":
							var msg planmessages.MsgRemoveNodeRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							filter = bson.M{
								"id": msg.ID,
							}
							update := bson.M{
								"$pull": bson.M{
									"node_addresses": explorerutils.MustHexFromBech32NodeAddress(msg.Address),
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.PlanFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
							if err != nil {
								return err
							}
						case "/sentinel.provider.v1.MsgRegisterRequest", "/sentinel.provider.v1.MsgService/MsgRegister":
							var msg providermessages.MsgRegisterRequest
							if err := json.Unmarshal(buf, &msg); err != nil {
								return err
							}

							dProvider := &types.Provider{
								Address:       explorerutils.MustHexFromBech32AccAddress(msg.From),
								Name:          msg.Name,
								Identity:      msg.Identity,
								Website:       msg.Website,
								Description:   msg.Description,
								JoinHeight:    dBlock.Height,
								JoinTimestamp: dBlock.Time,
								JoinTxHash:    dTxs[tIndex].Hash,
							}

							if err := database.ProviderSave(sctx, db, dProvider); err != nil {
								return err
							}
						case "/sentinel.provider.v1.MsgUpdateRequest", "/sentinel.provider.v1.MsgService/MsgUpdate":
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
							update := bson.M{
								"$set": bson.M{
									"name":        qProvider.Name,
									"identity":    qProvider.Identity,
									"website":     qProvider.Website,
									"description": qProvider.Description,
								},
							}
							projection = bson.M{
								"_id": 1,
							}

							_, err = database.ProviderFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
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
					case "sentinel.node.v1.EventUpdate":
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
								"price": commontypes.NewCoinsFromRaw(qNode.Price),
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.NodeFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}
					case "sentinel.node.v1.EventSetStatus":
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
					case "sentinel.session.v1.EventStakingReward", "sentinel.session.v1.EventPay":
						var rawCoin sdk.Coin
						if err := json.Unmarshal([]byte(dBlock.EndBlockEvents[eIndex].Attributes["amount"]), &rawCoin); err != nil {
							return err
						}

						id, err := strconv.ParseInt(dBlock.EndBlockEvents[eIndex].Attributes["id"], 10, 64)
						if err != nil {
							return err
						}

						coin := commontypes.NewCoinFromRaw(&rawCoin)

						filter = bson.M{
							"id": id,
						}

						updateSet := bson.M{
							"payment": coin,
						}

						if dBlock.EndBlockEvents[eIndex].Type == "sentinel.session.v1.EventStakingReward" {
							updateSet = bson.M{
								"staking_reward": coin,
							}
						}

						update := bson.M{
							"$set": updateSet,
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

						dDepositEvent := &types.DepositEvent{
							Address:   dSession.Address,
							Coins:     commontypes.Coins{coin},
							Subtract:  true,
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

						var rawCoins sdk.Coins
						if err := json.Unmarshal([]byte(dBlock.EndBlockEvents[eIndex].Attributes["coins"]), &rawCoins); err != nil {
							return err
						}

						coins := commontypes.NewCoinsFromRaw(rawCoins)

						dDepositEvent := &types.DepositEvent{
							Address:   hexAccAddr,
							Coins:     coins,
							Subtract:  true,
							Height:    dBlock.Height,
							Timestamp: dBlock.Time,
							TxHash:    "",
						}

						if err := database.DepositEventSave(sctx, db, dDepositEvent); err != nil {
							return err
						}
					case "sentinel.session.v1.EventSetStatus":
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
							"subscription_id": 1,
							"address":         1,
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

						qSubscriptionQuota, err := q.QuerySubscriptionQuota(dSession.SubscriptionID, accAddr, dBlock.Height)
						if err != nil {
							continue // TODO: cannot return error; subscription deletion before session deletion
						}

						filter = bson.M{
							"id":      dSession.SubscriptionID,
							"address": dSession.Address,
						}
						update = bson.M{
							"$set": bson.M{
								"consumed_bytes":  qSubscriptionQuota.Consumed.Int64(),
								"allocated_bytes": qSubscriptionQuota.Allocated.Int64(),
							},
						}
						projection = bson.M{
							"_id": 1,
						}

						_, err = database.SubscriptionQuotaFindOneAndUpdate(sctx, db, filter, update, options.FindOneAndUpdate().SetProjection(projection).SetUpsert(true))
						if err != nil {
							return err
						}

						dSubscriptionQuotaEvent := &types.SubscriptionQuotaEvent{
							ID:             dSession.SubscriptionID,
							Address:        dSession.Address,
							ConsumedBytes:  qSubscriptionQuota.Consumed.Int64(),
							AllocatedBytes: qSubscriptionQuota.Allocated.Int64(),
							Height:         dBlock.Height,
							Timestamp:      dBlock.Time,
							TxHash:         "",
						}

						if err := database.SubscriptionQuotaEventSave(sctx, db, dSubscriptionQuotaEvent); err != nil {
							return err
						}
					case "sentinel.subscription.v1.EventSetStatus":
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
