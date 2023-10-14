package operations

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewPlanCreateOperation(
	db *mongo.Database,
	v *models.Plan,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.PlanInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewPlanUpdateStatusOperation(
	db *mongo.Database,
	id uint64, status string, height int64, timestamp time.Time, txHash string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}
		update := bson.M{
			"$set": bson.M{
				"status":           status,
				"status_height":    height,
				"status_timestamp": timestamp,
				"status_tx_hash":   txHash,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.PlanFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}

func NewPlanAddNodeOperation(
	db *mongo.Database,
	id uint64, address string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}
		update := bson.M{
			"$addToSet": bson.M{
				"node_addresses": address,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.PlanFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}

func NewPlanRemoveNodeOperation(
	db *mongo.Database,
	id uint64, address string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id": id,
		}
		update := bson.M{
			"$pull": bson.M{
				"node_addresses": address,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.PlanFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
