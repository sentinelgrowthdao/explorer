package operations

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewProviderRegister(
	db *mongo.Database,
	v *models.Provider,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.ProviderInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewProviderUpdate(
	db *mongo.Database,
	addr, name, identity, website, description, status string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"addr": addr,
		}

		updateSet := bson.M{}
		if name != "" {
			updateSet["name"] = name
		}
		if identity != "" {
			updateSet["identity"] = identity
		}
		if website != "" {
			updateSet["website"] = website
		}
		if description != "" {
			updateSet["description"] = description
		}
		if status != "" {
			updateSet["status"] = status
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

		if _, err := database.ProviderFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
