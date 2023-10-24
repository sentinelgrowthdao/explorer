package operations

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewSubscriptionAllocationCreate(
	db *mongo.Database,
	v *models.SubscriptionAllocation,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.SubscriptionAllocationInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewSubscriptionAllocationUpdate(
	db *mongo.Database,
	id uint64, addr string, grantedBytes, utilisedBytes string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id":   id,
			"addr": addr,
		}
		update := bson.M{
			"$set": bson.M{
				"granted_bytes":  grantedBytes,
				"utilised_bytes": utilisedBytes,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.SubscriptionAllocationFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
