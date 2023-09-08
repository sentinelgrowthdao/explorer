package operations

import (
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewSessionStartOperation(
	db *mongo.Database,
	v *models.Session,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.SessionInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewSessionUpdateDetailsOperation(
	db *mongo.Database,
	id uint64, bandwidth *types.Bandwidth, duration int64, payment *types.Coin, rating int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}

		updateSet := bson.M{}
		if bandwidth != nil {
			updateSet["bandwidth"] = bandwidth
		}
		if duration != -1 {
			updateSet["duration"] = duration
		}
		if payment != nil {
			updateSet["payment"] = payment
		}
		if rating != -1 {
			updateSet["rating"] = rating
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

		if _, err := database.SessionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}

func NewSessionUpdateStatusOperation(
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

		if _, err := database.SessionFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
