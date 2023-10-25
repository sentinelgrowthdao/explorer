package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
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
	appName = "03_sentinelhub"
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
				bson.E{Key: "addr", Value: 1},
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
				bson.E{Key: "addr", Value: 1},
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
				bson.E{Key: "addr", Value: 1},
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
				bson.E{Key: "acc_addr", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.SubscriptionAllocationIndexesCreateMany(ctx, db, indexes)
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
		case "sentinel.deposit.v1.EventSubtract":
			event, err := deposittypes.NewEventSubtract(dBlock.BeginBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:      types.EventTypeDepositSubtract,
				Height:    dBlock.Height,
				Timestamp: dBlock.Time,
				TxHash:    "",
				AccAddr:   event.Address,
				Coins:     event.Coins,
			}

			ops = append(
				ops,
				operations.NewDepositSubtract(db, event.Address, event.Coins, dBlock.Height, dBlock.Time, ""),
				operations.NewEventCreate(db, &dEvent1),
			)
		case "sentinel.subscription.v2.EventPayForPayout":
			event, err := subscriptiontypes.NewEventPayForPayout(dBlock.BeginBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dSubscriptionPayout := models.SubscriptionPayout{
				ID:            event.ID,
				AccAddr:       event.Address,
				NodeAddr:      event.NodeAddress,
				Payment:       event.Payment,
				StakingReward: event.StakingReward,
				Height:        dBlock.Height,
				Timestamp:     dBlock.Time,
				TxHash:        "",
			}

			ops = append(
				ops,
				operations.NewSubscriptionPayoutCreate(db, &dSubscriptionPayout),
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
		log.Println("TxHash", dTxs[tIndex].Hash)
		log.Println("MessagesLen", tIndex, len(dTxs[tIndex].Messages))

		for eIndex, mIndex := -1, 0; mIndex < len(dTxs[tIndex].Messages); mIndex++ {
			log.Println("Type", dTxs[tIndex].Messages[mIndex].Type)

			if strings.Contains(dTxs[tIndex].Messages[mIndex].Type, "MsgExec") {
				msgs := dTxs[tIndex].Messages[mIndex].Data["msgs"].([]bson.M)
				for _, msg := range msgs {
					log.Println("MsgExec @type", msg["@type"].(string))
					if strings.Contains(msg["@type"].(string), "sentinel") {
						return nil, fmt.Errorf("invalid /cosmos.authz.v1beta1.MsgExec")
					}
				}
			}

			switch dTxs[tIndex].Messages[mIndex].Type {
			case "/sentinel.node.v2.MsgRegisterRequest", "/sentinel.node.v2.MsgService/MsgRegister":
				msg, err := nodetypes.NewMsgRegisterRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dNode := models.Node{
					Addr:              msg.NodeAddr().String(),
					GigabytePrices:    msg.GigabytePrices,
					HourlyPrices:      msg.HourlyPrices,
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
					operations.NewNodeRegister(db, &dNode),
				)
			case "/sentinel.node.v2.MsgUpdateDetailsRequest", "/sentinel.node.v2.MsgService/MsgUpdateDetails":
				msg, err := nodetypes.NewMsgUpdateDetailsRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:           types.EventTypeNodeUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					NodeAddr:       msg.From,
					GigabytePrices: msg.GigabytePrices,
					HourlyPrices:   msg.HourlyPrices,
					RemoteURL:      msg.RemoteURL,
				}

				ops = append(
					ops,
					operations.NewNodeUpdateDetails(db, msg.From, msg.GigabytePrices, msg.HourlyPrices, msg.RemoteURL),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.node.v2.MsgUpdateStatusRequest", "/sentinel.node.v2.MsgService/MsgUpdateStatus":
				msg, err := nodetypes.NewMsgUpdateStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:      types.EventTypeNodeUpdateStatus,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					NodeAddr:  msg.From,
					Status:    msg.Status,
				}

				ops = append(
					ops,
					operations.NewNodeUpdateStatus(db, msg.From, msg.Status, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.node.v2.MsgSubscribeRequest", "/sentinel.node.v2.MsgService/MsgSubscribe":
				msg, err := nodetypes.NewMsgSubscribeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				var (
					eventAdd                *deposittypes.EventAdd
					eventAllocate           *subscriptiontypes.EventAllocate
					eventCreateSubscription *nodetypes.EventCreateSubscription
				)

				eIndex, eventAdd, err = deposittypes.NewEventAddFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				if msg.Gigabytes != 0 {
					eIndex, eventAllocate, err = subscriptiontypes.NewEventAllocateFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
					if err != nil {
						return nil, err
					}
				}

				eIndex, eventCreateSubscription, err = nodetypes.NewEventCreateSubscriptionFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				inactiveAt := dBlock.Time.Add(90 * 24 * time.Hour)
				if msg.Hours != 0 {
					inactiveAt = dBlock.Time.Add(time.Duration(msg.Hours) * time.Hour)
				}

				dSubscription := models.Subscription{
					ID:              eventCreateSubscription.ID,
					AccAddr:         msg.From,
					NodeAddr:        msg.NodeAddress,
					Gigabytes:       msg.Gigabytes,
					Hours:           msg.Hours,
					Price:           nil,
					Deposit:         eventAdd.Coins[0],
					Refund:          nil,
					InactiveAt:      inactiveAt,
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

				dEvent1 := models.Event{
					Type:      types.EventTypeDepositAdd,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					AccAddr:   eventAdd.Address,
					Coins:     eventAdd.Coins,
				}

				ops = append(
					ops,
					operations.NewSubscriptionCreate(db, &dSubscription),
					operations.NewDepositAdd(db, eventAdd.Address, eventAdd.Coins, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash),
					operations.NewEventCreate(db, &dEvent1),
				)

				if msg.Gigabytes != 0 {
					dSubscriptionAllocation := models.SubscriptionAllocation{
						ID:            eventAllocate.ID,
						AccAddr:       eventAllocate.Address,
						GrantedBytes:  eventAllocate.GrantedBytes,
						UtilisedBytes: eventAllocate.UtilisedBytes,
					}

					dEvent1 := models.Event{
						Type:           types.EventTypeSubscriptionAllocationUpdateDetails,
						Height:         dBlock.Height,
						Timestamp:      dBlock.Time,
						TxHash:         dTxs[tIndex].Hash,
						SubscriptionID: eventAllocate.ID,
						AccAddr:        eventAllocate.Address,
						GrantedBytes:   eventAllocate.GrantedBytes,
						UtilisedBytes:  eventAllocate.UtilisedBytes,
					}

					ops = append(
						ops,
						operations.NewSubscriptionAllocationCreate(db, &dSubscriptionAllocation),
						operations.NewEventCreate(db, &dEvent1),
					)
				}
			case "/sentinel.plan.v2.MsgCreateRequest", "/sentinel.plan.v2.MsgService/MsgCreate":
				msg, err := plantypes.NewMsgCreateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				var (
					eventCreate *plantypes.EventCreate
				)

				eIndex, eventCreate, err = plantypes.NewEventCreateFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				dPlan := models.Plan{
					ID:              eventCreate.ID,
					ProvAddr:        msg.From,
					Prices:          msg.Prices,
					Duration:        msg.Duration,
					Gigabytes:       msg.Gigabytes,
					NodeAddrs:       []string{},
					CreateHeight:    dBlock.Height,
					CreateTimestamp: dBlock.Time,
					CreateTxHash:    dTxs[tIndex].Hash,
					Status:          hubtypes.StatusInactive.String(),
					StatusHeight:    dBlock.Height,
					StatusTimestamp: dBlock.Time,
					StatusTxHash:    dTxs[tIndex].Hash,
				}

				ops = append(
					ops,
					operations.NewPlanCreate(db, &dPlan),
				)
			case "/sentinel.plan.v2.MsgUpdateStatusRequest", "/sentinel.plan.v2.MsgService/MsgUpdateStatus":
				msg, err := plantypes.NewMsgUpdateStatusRequest(dTxs[tIndex].Messages[mIndex].Data)
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
					operations.NewPlanUpdateStatus(db, msg.ID, msg.Status, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.plan.v2.MsgLinkNodeRequest", "/sentinel.plan.v2.MsgService/MsgLinkNode":
				msg, err := plantypes.NewMsgLinkNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:      types.EventTypePlanLinkNode,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					PlanID:    msg.ID,
					NodeAddr:  msg.NodeAddress,
				}

				ops = append(
					ops,
					operations.NewPlanLinkNode(db, msg.ID, msg.NodeAddress),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.plan.v2.MsgUnlinkNodeRequest", "/sentinel.plan.v2.MsgService/MsgUnlinkNode":
				msg, err := plantypes.NewMsgUnlinkNodeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:      types.EventTypePlanUnlinkNode,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					PlanID:    msg.ID,
					NodeAddr:  msg.NodeAddress,
				}

				ops = append(
					ops,
					operations.NewPlanUnlinkNode(db, msg.ID, msg.NodeAddress),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.plan.v2.MsgSubscribeRequest", "/sentinel.plan.v2.MsgService/MsgSubscribe":
				msg, err := plantypes.NewMsgSubscribeRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				var (
					eventPayForPlan         *subscriptiontypes.EventPayForPlan
					eventAllocate           *subscriptiontypes.EventAllocate
					eventCreateSubscription *plantypes.EventCreateSubscription
				)

				eIndex, eventPayForPlan, err = subscriptiontypes.NewEventPayForPlanFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				eIndex, eventAllocate, err = subscriptiontypes.NewEventAllocateFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				eIndex, eventCreateSubscription, err = plantypes.NewEventCreateSubscriptionFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				dSubscription := models.Subscription{
					ID:              eventCreateSubscription.ID,
					AccAddr:         msg.From,
					PlanID:          msg.ID,
					Price:           nil,
					Payment:         eventPayForPlan.Payment,
					StakingReward:   eventPayForPlan.StakingReward,
					InactiveAt:      time.Time{},
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

				dSubscriptionAllocation := models.SubscriptionAllocation{
					ID:            eventAllocate.ID,
					AccAddr:       eventAllocate.Address,
					GrantedBytes:  eventAllocate.GrantedBytes,
					UtilisedBytes: eventAllocate.UtilisedBytes,
				}

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionAllocationUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAllocate.ID,
					AccAddr:        eventAllocate.Address,
					GrantedBytes:   eventAllocate.GrantedBytes,
					UtilisedBytes:  eventAllocate.UtilisedBytes,
				}

				ops = append(
					ops,
					operations.NewSubscriptionCreate(db, &dSubscription),
					operations.NewSubscriptionAllocationCreate(db, &dSubscriptionAllocation),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.provider.v2.MsgRegisterRequest", "/sentinel.provider.v2.MsgService/MsgRegister":
				msg, err := providertypes.NewMsgRegisterRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dProvider := models.Provider{
					Addr:              msg.ProvAddr().String(),
					Name:              msg.Name,
					Identity:          msg.Identity,
					Website:           msg.Website,
					Description:       msg.Description,
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
					operations.NewProviderRegister(db, &dProvider),
				)
			case "/sentinel.provider.v2.MsgUpdateRequest", "/sentinel.provider.v2.MsgService/MsgUpdate":
				msg, err := providertypes.NewMsgUpdateRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:        types.EventTypeProviderUpdateDetails,
					Height:      dBlock.Height,
					Timestamp:   dBlock.Time,
					TxHash:      dTxs[tIndex].Hash,
					ProvAddr:    msg.From,
					Name:        msg.Name,
					Identity:    msg.Identity,
					Website:     msg.Website,
					Description: msg.Description,
					Status:      msg.Status,
				}

				ops = append(
					ops,
					operations.NewProviderUpdate(db, msg.From, msg.Name, msg.Identity, msg.Website, msg.Description, msg.Status),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.session.v2.MsgStartRequest", "/sentinel.session.v2.MsgService/MsgStart":
				msg, err := sessiontypes.NewMsgStartRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				var (
					eventStart *sessiontypes.EventStart
				)

				eIndex, eventStart, err = sessiontypes.NewEventStartFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				dSession := models.Session{
					ID:              eventStart.ID,
					SubscriptionID:  msg.ID,
					AccAddr:         msg.From,
					NodeAddr:        msg.NodeAddress,
					Bandwidth:       nil,
					Duration:        0,
					Payment:         nil,
					StakingReward:   nil,
					Rating:          0,
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

				ops = append(
					ops,
					operations.NewSessionCreate(db, &dSession),
				)
			case "/sentinel.session.v2.MsgUpdateDetailsRequest", "/sentinel.session.v2.MsgService/MsgUpdate":
				msg, err := sessiontypes.NewMsgUpdateDetailsRequest(dTxs[tIndex].Messages[mIndex].Data)
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
					operations.NewSessionUpdateDetails(db, msg.ID, msg.Bandwidth, msg.Duration, nil, nil, -1),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.session.v2.MsgEndRequest", "/sentinel.session.v2.MsgService/MsgEnd":
				msg, err := sessiontypes.NewMsgEndRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				status := hubtypes.StatusInactivePending.String()
				dEvent1 := models.Event{
					Type:      types.EventTypeSessionUpdateStatus,
					Height:    dBlock.Height,
					Timestamp: dBlock.Time,
					TxHash:    dTxs[tIndex].Hash,
					SessionID: msg.ID,
					Status:    status,
				}

				ops = append(
					ops,
					operations.NewSessionUpdateDetails(db, msg.ID, nil, -1, nil, nil, msg.Rating),
					operations.NewSessionUpdateStatus(db, msg.ID, status, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.subscription.v2.MsgCancelRequest", "/sentinel.subscription.v2.MsgService/MsgCancel":
				msg, err := subscriptiontypes.NewMsgCancelRequest(dTxs[tIndex].Messages[mIndex].Data)
				if err != nil {
					return nil, err
				}

				status := hubtypes.StatusInactivePending.String()
				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionUpdateStatus,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: msg.ID,
					Status:         status,
				}

				ops = append(
					ops,
					operations.NewSubscriptionUpdateStatus(db, msg.ID, status, dBlock.Height, dBlock.Time, dTxs[tIndex].Hash),
					operations.NewEventCreate(db, &dEvent1),
				)
			case "/sentinel.subscription.v2.MsgAllocateRequest", "/sentinel.subscription.v2.MsgService/MsgAllocate":
				var (
					eventAllocate1 *subscriptiontypes.EventAllocate
					eventAllocate2 *subscriptiontypes.EventAllocate
				)

				eIndex, eventAllocate1, err = subscriptiontypes.NewEventAllocateFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				dEvent1 := models.Event{
					Type:           types.EventTypeSubscriptionAllocationUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAllocate1.ID,
					AccAddr:        eventAllocate1.Address,
					GrantedBytes:   eventAllocate1.GrantedBytes,
					UtilisedBytes:  eventAllocate1.UtilisedBytes,
				}

				eIndex, eventAllocate2, err = subscriptiontypes.NewEventAllocateFromEvents(dTxs[tIndex].Result.Events[eIndex+1:])
				if err != nil {
					return nil, err
				}

				dEvent2 := models.Event{
					Type:           types.EventTypeSubscriptionAllocationUpdateDetails,
					Height:         dBlock.Height,
					Timestamp:      dBlock.Time,
					TxHash:         dTxs[tIndex].Hash,
					SubscriptionID: eventAllocate1.ID,
					AccAddr:        eventAllocate1.Address,
					GrantedBytes:   eventAllocate1.GrantedBytes,
					UtilisedBytes:  eventAllocate1.UtilisedBytes,
				}

				ops = append(
					ops,
					operations.NewSubscriptionAllocationUpdate(db, eventAllocate1.ID, eventAllocate1.Address, eventAllocate1.GrantedBytes, eventAllocate1.UtilisedBytes),
					operations.NewEventCreate(db, &dEvent1),
					operations.NewSubscriptionAllocationUpdate(db, eventAllocate2.ID, eventAllocate2.Address, eventAllocate2.GrantedBytes, eventAllocate2.UtilisedBytes),
					operations.NewEventCreate(db, &dEvent2),
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
				Type:      types.EventTypeDepositSubtract,
				Height:    dBlock.Height,
				Timestamp: dBlock.Time,
				TxHash:    "",
				AccAddr:   event.Address,
				Coins:     event.Coins,
			}

			ops = append(
				ops,
				operations.NewDepositSubtract(db, event.Address, event.Coins, dBlock.Height, dBlock.Time, ""),
				operations.NewEventCreate(db, &dEvent1),
			)
		case "sentinel.node.v2.EventUpdateDetails":
			event, err := nodetypes.NewEventUpdateDetails(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:           types.EventTypeNodeUpdateDetails,
				Height:         dBlock.Height,
				Timestamp:      dBlock.Time,
				TxHash:         "",
				NodeAddr:       event.Address,
				GigabytePrices: event.GigabytePrices,
				HourlyPrices:   event.HourlyPrices,
				RemoteURL:      event.RemoteURL,
			}

			ops = append(
				ops,
				operations.NewNodeUpdateDetails(db, event.Address, event.GigabytePrices, event.HourlyPrices, event.RemoteURL),
				operations.NewEventCreate(db, &dEvent1),
			)
		case "sentinel.node.v2.EventUpdateStatus":
			event, err := nodetypes.NewEventUpdateStatus(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:      types.EventTypeNodeUpdateStatus,
				Height:    dBlock.Height,
				Timestamp: dBlock.Time,
				TxHash:    "",
				NodeAddr:  event.Address,
				Status:    event.Status,
			}

			ops = append(
				ops,
				operations.NewNodeUpdateStatus(db, event.Address, event.Status, dBlock.Height, dBlock.Time, ""),
				operations.NewEventCreate(db, &dEvent1),
			)
		case "sentinel.session.v2.EventUpdateStatus":
			event, err := sessiontypes.NewEventUpdateStatus(dBlock.EndBlockEvents[eIndex])
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
				operations.NewSessionUpdateStatus(db, event.ID, event.Status, dBlock.Height, dBlock.Time, ""),
				operations.NewEventCreate(db, &dEvent1),
			)
		case "sentinel.subscription.v2.EventPayForSession":
			event, err := subscriptiontypes.NewEventPayForSession(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewSessionUpdateDetails(db, event.ID, nil, -1, event.Payment, event.StakingReward, -1),
			)
		case "sentinel.subscription.v2.EventUpdateStatus":
			event, err := subscriptiontypes.NewEventUpdateStatus(dBlock.EndBlockEvents[eIndex])
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
				operations.NewSubscriptionUpdateStatus(db, event.ID, event.Status, dBlock.Height, dBlock.Time, ""),
				operations.NewEventCreate(db, &dEvent1),
			)
		case "sentinel.subscription.v2.EventRefund":
			event, err := subscriptiontypes.NewEventRefund(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			ops = append(
				ops,
				operations.NewSubscriptionUpdateDetails(db, event.ID, event.Amount),
			)
		case "sentinel.subscription.v2.EventAllocate":
			event, err := subscriptiontypes.NewEventAllocate(dBlock.EndBlockEvents[eIndex])
			if err != nil {
				return nil, err
			}

			dEvent1 := models.Event{
				Type:           types.EventTypeSubscriptionAllocationUpdateDetails,
				Height:         dBlock.Height,
				Timestamp:      dBlock.Time,
				TxHash:         "",
				SubscriptionID: event.ID,
				AccAddr:        event.Address,
				GrantedBytes:   event.GrantedBytes,
				UtilisedBytes:  event.UtilisedBytes,
			}

			ops = append(
				ops,
				operations.NewSubscriptionAllocationUpdate(db, event.ID, event.Address, event.GrantedBytes, event.UtilisedBytes),
				operations.NewEventCreate(db, &dEvent1),
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
