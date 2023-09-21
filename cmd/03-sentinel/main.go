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
	"github.com/sentinel-official/explorer/operations"
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
	appName = "03-sentinel"
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
	flag.Int64Var(&toHeight, "to-height", 5_125_000, "")
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

				ops = append(
					ops,
					operations.NewNodeRegisterOperation(
						db, &dNode,
					),
				)
			case "/sentinel.node.v1.MsgUpdateRequest", "/sentinel.node.v1.MsgService/MsgUpdate":
				msg, err := nodetypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

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

				ops = append(
					ops,
					operations.NewNodeUpdateDetailsOperation(
						db, msg.From, msg.Provider, msg.Price, msg.RemoteURL,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
			case "/sentinel.node.v1.MsgSetStatusRequest", "/sentinel.node.v1.MsgService/MsgSetStatus":
				msg, err := nodetypes.NewMsgSetStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:        types.EventTypeNodeUpdateStatus,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					NodeAddress: msg.From,
					Status:      msg.Status,
				}

				ops = append(
					ops,
					operations.NewNodeUpdateStatusOperation(
						db, msg.From, msg.Status, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)

			case "/sentinel.plan.v1.MsgAddRequest", "/sentinel.plan.v1.MsgService/MsgAdd":
				msg, err := plantypes.NewMsgAddRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				eventAddPlan, err := plantypes.NewEventAddPlanFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				dPlan := models.Plan{
					ID:              eventAddPlan.ID,
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

				ops = append(
					ops,
					operations.NewPlanCreateOperation(
						db, &dPlan,
					),
				)
			case "/sentinel.plan.v1.MsgSetStatusRequest", "/sentinel.plan.v1.MsgService/MsgSetStatus":
				msg, err := plantypes.NewMsgSetStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:      types.EventTypePlanUpdateStatus,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					PlanID:    msg.ID,
					Status:    msg.Status,
				}

				ops = append(
					ops,
					operations.NewPlanUpdateStatusOperation(
						db, msg.ID, msg.Status, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
			case "/sentinel.plan.v1.MsgAddNodeRequest", "/sentinel.plan.v1.MsgService/MsgAddNode":
				msg, err := plantypes.NewMsgAddNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:        types.EventTypePlanAddNode,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					PlanID:      msg.ID,
					NodeAddress: msg.Address,
				}

				ops = append(
					ops,
					operations.NewPlanAddNodeOperation(
						db, msg.ID, msg.Address,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
			case "/sentinel.plan.v1.MsgRemoveNodeRequest", "/sentinel.plan.v1.MsgService/MsgRemoveNode":
				msg, err := plantypes.NewMsgRemoveNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:        types.EventTypePlanRemoveNode,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					PlanID:      msg.ID,
					NodeAddress: msg.Address,
				}

				ops = append(
					ops,
					operations.NewPlanRemoveNodeOperation(
						db, msg.ID, msg.Address,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
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

				ops = append(
					ops,
					operations.NewProviderRegisterOperation(
						db, &dProvider,
					),
				)
			case "/sentinel.provider.v1.MsgUpdateRequest", "/sentinel.provider.v1.MsgService/MsgUpdate":
				msg, err := providertypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

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

				ops = append(
					ops,
					operations.NewProviderUpdateOperation(
						db, msg.From, msg.Name, msg.Identity, msg.Website, msg.Description,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
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

				ops = append(
					ops,
					operations.NewSessionStartOperation(
						db, &dSession,
					),
				)
			case "/sentinel.session.v1.MsgUpdateRequest", "/sentinel.session.v1.MsgService/MsgUpdate":
				msg, err := sessiontypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:      types.EventTypeSessionUpdateDetails,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					SessionID: msg.ID,
					Bandwidth: msg.Bandwidth,
					Duration:  msg.Duration,
				}

				ops = append(
					ops,
					operations.NewSessionUpdateDetailsOperation(
						db, msg.ID, msg.Bandwidth, msg.Duration, nil, -1,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
			case "/sentinel.session.v1.MsgEndRequest", "/sentinel.session.v1.MsgService/MsgEnd":
				msg, err := sessiontypes.NewMsgEndRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:      types.EventTypeSessionUpdateStatus,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					SessionID: msg.ID,
					Status:    hubtypes.StatusInactivePending.String(),
				}

				ops = append(
					ops,
					operations.NewSessionUpdateDetailsOperation(
						db, msg.ID, nil, -1, nil, msg.Rating,
					),
					operations.NewSessionUpdateStatusOperation(
						db, msg.ID, hubtypes.StatusInactivePending.String(), dBlock.Height, dBlock.Time, dTxs[tIndex].Hash,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
			case "/sentinel.subscription.v1.MsgSubscribeToNodeRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToNode":
				eventSubscribeToNode, err := subscriptiontypes.NewEventSubscribeToNodeFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				eventAddQuota, err := subscriptiontypes.NewEventAddQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				eventAddDeposit, err := deposittypes.NewEventAddFromEvents(txResultLog[mIndex].Events)
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

				dSubscriptionQuota := models.SubscriptionQuota{
					ID:        eventAddQuota.ID,
					Address:   eventAddQuota.Address,
					Allocated: eventAddQuota.Allocated,
					Consumed:  eventAddQuota.Consumed,
				}

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

				dEvent2 := models.Event{
					Type:       types.EventTypeDepositAdd,
					Height:     dBlock.Height,
					Timestamp:  dBlock.Time,
					TxHash:     dTxs[tIndex].Hash,
					AccAddress: eventAddDeposit.Address,
					Coins:      eventAddDeposit.Coins,
				}

				ops = append(
					ops,
					operations.NewSubscriptionCreateOperation(
						db, &dSubscription,
					),
					operations.NewSubscriptionQuotaAddOperation(
						db, &dSubscriptionQuota,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
					operations.NewDepositUpdateOperation(
						db, eventAddDeposit.Address, eventAddDeposit.Current, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash,
					),
					operations.NewEventSaveOperation(
						db, &dEvent2,
					),
				)
			case "/sentinel.subscription.v1.MsgSubscribeToPlanRequest", "/sentinel.subscription.v1.MsgService/MsgSubscribeToPlan":
				eventSubscribeToPlan, err := subscriptiontypes.NewEventSubscribeToPlanFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

				eventAddQuota, err := subscriptiontypes.NewEventAddQuotaFromEvents(txResultLog[mIndex].Events)
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

				dSubscriptionQuota := models.SubscriptionQuota{
					ID:        eventAddQuota.ID,
					Address:   eventAddQuota.Address,
					Allocated: eventAddQuota.Allocated,
					Consumed:  eventAddQuota.Consumed,
				}

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

				ops = append(
					ops,
					operations.NewSubscriptionCreateOperation(
						db, &dSubscription,
					),
					operations.NewSubscriptionQuotaAddOperation(
						db, &dSubscriptionQuota,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
			case "/sentinel.subscription.v1.MsgCancelRequest", "/sentinel.subscription.v1.MsgService/MsgCancel":
				msg, err := subscriptiontypes.NewMsgCancelRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateStatus,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: msg.ID,
					Status:         hubtypes.StatusInactivePending.String(),
				}

				ops = append(
					ops,
					operations.NewSubscriptionUpdateStatusOperation(
						db, msg.ID, hubtypes.StatusInactivePending.String(), dBlock.Height, dBlock.Time, dTxs[tIndex].Hash,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
				)
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

				dEvent2 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAddQuota.ID,
					Free:           eventAddQuota.Free,
				}

				ops = append(
					ops,
					operations.NewSubscriptionQuotaAddOperation(
						db, &dSubscriptionQuota,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
					operations.NewSubscriptionUpdateDetailsOperation(
						db, eventAddQuota.ID, eventAddQuota.Free, nil,
					),
					operations.NewEventSaveOperation(
						db, &dEvent2,
					),
				)
			case "/sentinel.subscription.v1.MsgUpdateQuotaRequest", "/sentinel.subscription.v1.MsgService/MsgUpdateQuota":
				eventUpdateQuota, err := subscriptiontypes.NewEventUpdateQuotaFromEvents(txResultLog[mIndex].Events)
				if err != nil {
					return nil, err
				}

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

				dEvent2 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventUpdateQuota.ID,
					Free:           eventUpdateQuota.Free,
				}

				ops = append(
					ops,
					operations.NewSubscriptionQuotaUpdateOperation(
						db, eventUpdateQuota.ID, eventUpdateQuota.Address, eventUpdateQuota.Allocated, eventUpdateQuota.Consumed,
					),
					operations.NewEventSaveOperation(
						db, &dEvent1,
					),
					operations.NewSubscriptionUpdateDetailsOperation(
						db, eventUpdateQuota.ID, eventUpdateQuota.Free, nil,
					),
					operations.NewEventSaveOperation(
						db, &dEvent2,
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
		case "sentinel.deposit.v1.EventSubtract":
			event, err := deposittypes.NewEventSubtract(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:       types.EventTypeDepositSubtract,
				Height:     dBlock.Height,
				Timestamp:  dBlock.Time,
				TxHash:     "",
				AccAddress: event.Address,
				Coins:      event.Coins,
			}

			ops = append(
				ops,
				operations.NewDepositUpdateOperation(
					db, event.Address, event.Current, dBlock.Height, dBlock.Time, "",
				),
				operations.NewEventSaveOperation(
					db, &dEvent1,
				),
			)
		case "sentinel.node.v1.EventSetNodeStatus":
			event, err := nodetypes.NewEventSetNodeStatus(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:        types.EventTypeNodeUpdateStatus,
				Height:      dBlock.Height,
				Timestamp:   dBlock.Time,
				TxHash:      "",
				NodeAddress: event.Address,
				Status:      event.Status,
			}

			ops = append(
				ops,
				operations.NewNodeUpdateStatusOperation(
					db, event.Address, event.Status, dBlock.Height, dBlock.Time, "",
				),
				operations.NewEventSaveOperation(
					db, &dEvent1,
				),
			)
		case "sentinel.session.v1.EventEndSession":
			event, err := sessiontypes.NewEventEndSession(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:      types.EventTypeSessionUpdateStatus,
				Height:    dBlock.Height,
				Timestamp: dBlock.Time,
				TxHash:    "",
				SessionID: event.ID,
				Status:    event.Status,
			}

			ops = append(
				ops,
				operations.NewSessionUpdateStatusOperation(
					db, event.ID, event.Status, dBlock.Height, dBlock.Time, "",
				),
				operations.NewEventSaveOperation(
					db, &dEvent1,
				),
			)
		case "sentinel.session.v1.EventPay":
			event, err := sessiontypes.NewEventPay(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewSessionUpdateDetailsOperation(
					db, event.ID, nil, -1, event.Payment, -1,
				),
			)
		case "sentinel.subscription.v1.EventCancelSubscription":
			event, err := subscriptiontypes.NewEventCancelSubscription(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:           types.EventTypeSubscriptionUpdateStatus,
				Height:         dBlock.Height,
				Timestamp:      dBlock.Time,
				TxHash:         "",
				SubscriptionID: event.ID,
				Status:         event.Status,
			}

			ops = append(
				ops,
				operations.NewSubscriptionUpdateStatusOperation(
					db, event.ID, event.Status, dBlock.Height, dBlock.Time, "",
				),
				operations.NewEventSaveOperation(
					db, &dEvent1,
				),
			)
		case "sentinel.subscription.v1.EventRefund":
			event, err := subscriptiontypes.NewEventRefund(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewSubscriptionUpdateDetailsOperation(
					db, event.ID, nil, event.Refund,
				),
			)
		case "sentinel.subscription.v1.EventUpdateQuota":
			event, err := subscriptiontypes.NewEventUpdateQuota(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

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

			ops = append(
				ops,
				operations.NewSubscriptionQuotaUpdateOperation(
					db, event.ID, event.Address, event.Allocated, event.Consumed,
				),
				operations.NewEventSaveOperation(
					db, &dEvent1,
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

				log.Println("OperationsLen", len(ops))
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
