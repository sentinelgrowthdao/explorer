package operations

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

func NewNodeStatisticUpdateSessionStartCount(
	db *mongo.Database,
	address string, timestamp time.Time, v int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"address":   address,
			"timestamp": timestamp,
		}
		update := bson.M{
			"$inc": bson.M{
				"session_start_count": v,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeStatisticUpdateSessionEndCount(
	db *mongo.Database,
	timestamp time.Time, id uint64, v int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}
		projection := bson.M{
			"_id":  0,
			"node": 1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		item, err := database.SessionFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}

		filter = bson.M{
			"address":   item.Node,
			"timestamp": timestamp,
		}
		update := bson.M{
			"$inc": bson.M{
				"session_end_count": v,
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeStatisticUpdateSubscriptionStartCount(
	db *mongo.Database,
	address string, timestamp time.Time, v int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"address":   address,
			"timestamp": timestamp,
		}
		update := bson.M{
			"$inc": bson.M{
				"subscription_start_count": v,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeStatisticUpdateSubscriptionBytes(
	db *mongo.Database,
	address string, timestamp time.Time, v sdk.Int,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"address":   address,
			"timestamp": timestamp,
		}
		projection := bson.M{
			"_id":                0,
			"subscription_bytes": 1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		item, err := database.NodeStatisticFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}
		if item != nil {
			v = v.Add(utils.MustIntFromString(item.SubscriptionBytes))
		}

		update := bson.M{
			"$set": bson.M{
				"subscription_bytes": v.String(),
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeStatisticUpdateEarningsForBytes(
	db *mongo.Database,
	timestamp time.Time, id uint64, v *types.Coin,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}
		projection := bson.M{
			"_id":  0,
			"node": 1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		session, err := database.SessionFindOne(ctx, db, filter)
		if err != nil {
			return err
		}

		filter = bson.M{
			"address":   session.Node,
			"timestamp": timestamp,
		}
		projection = bson.M{
			"_id":                0,
			"earnings_for_bytes": 1,
		}
		findOneOpts = options.FindOne().
			SetProjection(projection)

		item, err := database.NodeStatisticFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}
		if item == nil {
			item = models.NewNodeStatistic()
		}
		if item.EarningsForBytes == nil {
			item.EarningsForBytes = types.NewCoins(nil)
		}

		update := bson.M{
			"$set": bson.M{
				"earnings_for_bytes": item.EarningsForBytes.Add(v),
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeStatisticUpdateSubscriptionEndCount(
	db *mongo.Database,
	timestamp time.Time, id uint64, v int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}
		projection := bson.M{
			"_id":  0,
			"node": 1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		item, err := database.SubscriptionFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}

		filter = bson.M{
			"address":   item.Node,
			"timestamp": timestamp,
		}
		update := bson.M{
			"$inc": bson.M{
				"subscription_end_count": v,
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeStatisticUpdateSessionDetails(
	db *mongo.Database,
	address string, timestamp time.Time, id uint64, bandwidth *types.Bandwidth, duration int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"address":   address,
			"timestamp": timestamp,
		}
		projection := bson.M{
			"_id":               0,
			"session_bandwidth": 1,
			"session_duration":  1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		item, err := database.NodeStatisticFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}
		if item == nil {
			item = models.NewNodeStatistic()
		}
		if item.SessionBandwidth == nil {
			item.SessionBandwidth = types.NewBandwidth(nil)
		}

		filter = bson.M{
			"session_id": id,
			"timestamp": bson.M{
				"$lt": timestamp,
			},
			"type": types.EventTypeSessionUpdateDetails,
		}
		projection = bson.M{
			"_id":       0,
			"bandwidth": 1,
			"duration":  1,
			"timestamp": 1,
		}
		sort := bson.D{
			bson.E{Key: "timestamp", Value: -1},
		}
		findOpts := options.Find().
			SetLimit(1).
			SetProjection(projection).
			SetSort(sort)

		items, err := database.EventFind(ctx, db, filter, findOpts)
		if err != nil {
			return err
		}
		if len(items) > 0 {
			item.SessionBandwidth = item.SessionBandwidth.Sub(items[0].Bandwidth)
			item.SessionDuration = item.SessionDuration - items[0].Duration
		}

		item.SessionBandwidth = item.SessionBandwidth.Add(bandwidth)
		item.SessionDuration = item.SessionDuration + duration

		filter = bson.M{
			"address":   address,
			"timestamp": timestamp,
		}
		update := bson.M{
			"$set": bson.M{
				"session_bandwidth": item.SessionBandwidth,
				"session_duration":  item.SessionDuration,
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		findOneAndUpdateOpts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.NodeStatisticFindOneAndUpdate(ctx, db, filter, update, findOneAndUpdateOpts); err != nil {
			return err
		}

		return nil
	}
}
