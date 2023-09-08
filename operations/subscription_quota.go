package operations

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewSubscriptionQuotaAddOperation(
	db *mongo.Database,
	v *models.SubscriptionQuota,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.SubscriptionQuotaInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewSubscriptionQuotaUpdateOperation(
	db *mongo.Database,
	id uint64, address string, allocated, consumed int64,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"id":      id,
			"address": address,
		}
		update := bson.M{
			"$set": bson.M{
				"allocated": allocated,
				"consumed":  consumed,
			},
		}
		projection := bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.SubscriptionQuotaFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
