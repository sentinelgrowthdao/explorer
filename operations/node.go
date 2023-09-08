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

func NewNodeRegisterOperation(
	db *mongo.Database,
	v *models.Node,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.NodeInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeUpdateDetailsOperation(
	db *mongo.Database,
	address, provider string, price types.Coins, remoteURL string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"address": address,
		}

		updateSet := bson.M{}
		if provider != "" {
			updateSet["provider"], updateSet["price"] = provider, nil
		}
		if price != nil && len(price) > 0 {
			updateSet["price"], updateSet["provider"] = price, ""
		}
		if remoteURL != "" {
			updateSet["remote_url"] = remoteURL
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

		if _, err := database.NodeFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}

func NewNodeUpdateStatusOperation(
	db *mongo.Database,
	address, status string, height int64, timestamp time.Time, txHash string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"address": address,
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

		if _, err := database.NodeFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
