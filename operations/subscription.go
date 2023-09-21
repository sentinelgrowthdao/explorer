package operations

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewSubscriptionCreateOperation(
	db *mongo.Database,
	v *models.Subscription,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.SubscriptionInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewSubscriptionUpdateDetailsOperation(
	db *mongo.Database,
	id uint64, free sdk.Int, refund *types.Coin,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}

		updateSet := bson.M{}
		if !free.IsNil() {
			updateSet["free"] = free
		}
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

func NewSubscriptionUpdateStatusOperation(
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
