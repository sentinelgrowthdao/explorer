package operations

import (
	"fmt"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewSubscriptionCreate(
	db *mongo.Database,
	v *models.Subscription,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if v.NodeAddr != "" {
			filter := bson.M{
				"addr": v.NodeAddr,
			}
			projection := bson.M{
				"_id":             0,
				"gigabyte_prices": 1,
				"hourly_prices":   1,
			}
			opts := options.FindOne().
				SetProjection(projection)

			item, err := database.NodeFindOne(ctx, db, filter, opts)
			if err != nil {
				return err
			}

			if v.Gigabytes != 0 {
				v.Price = item.GigabytePrices.Get(v.Deposit.Denom).Copy()
			}
			if v.Hours != 0 {
				v.Price = item.HourlyPrices.Get(v.Deposit.Denom).Copy()
			}
		}

		if v.PlanID != 0 {
			filter := bson.M{
				"id": v.PlanID,
			}
			projection := bson.M{
				"_id":      0,
				"duration": 1,
				"prices":   1,
			}
			opts := options.FindOne().
				SetProjection(projection)

			item, err := database.PlanFindOne(ctx, db, filter, opts)
			if err != nil {
				return err
			}

			v.Price = item.Prices.Get(v.Payment.Denom).Copy()
			v.InactiveAt = v.StartTimestamp.Add(time.Duration(item.Duration))
		}

		if v.Price == nil {
			return fmt.Errorf("price is nil")
		}

		if _, err := database.SubscriptionInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewSubscriptionUpdateDetails(
	db *mongo.Database,
	id uint64, refund *types.Coin,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}

		updateSet := bson.M{}
		if refund != nil {
			updateSet["refund"] = refund
		}

		update := bson.M{
			"$set": updateSet,
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}

func NewSubscriptionUpdateStatus(
	db *mongo.Database,
	id uint64, status string, height int64, timestamp time.Time, txHash string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}

		updateSet := bson.M{
			"status":           status,
			"status_height":    height,
			"status_timestamp": timestamp,
			"status_tx_hash":   txHash,
		}
		if status == hubtypes.StatusInactive.String() {
			updateSet["end_height"] = height
			updateSet["end_timestamp"] = timestamp
			updateSet["end_tx_hash"] = ""
		}

		update := bson.M{
			"$set": updateSet,
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.SubscriptionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
