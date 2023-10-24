package operations

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewDepositAdd(
	db *mongo.Database,
	addr string, coins types.Coins, height int64, timestamp time.Time, txHash string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"addr": addr,
		}
		projection := bson.M{
			"_id":   0,
			"coins": 1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		item, err := database.DepositFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}
		if item == nil {
			item = models.NewDeposit()
		}

		update := bson.M{
			"$set": bson.M{
				"coins":     item.Coins.Add(coins...),
				"height":    height,
				"timestamp": timestamp,
				"tx_hash":   txHash,
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.DepositFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}

func NewDepositSubtract(
	db *mongo.Database,
	addr string, coins types.Coins, height int64, timestamp time.Time, txHash string,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		filter := bson.M{
			"addr": addr,
		}
		projection := bson.M{
			"_id":   0,
			"coins": 1,
		}
		findOneOpts := options.FindOne().
			SetProjection(projection)

		item, err := database.DepositFindOne(ctx, db, filter, findOneOpts)
		if err != nil {
			return err
		}
		if item == nil {
			return fmt.Errorf("nil deposit")
		}

		update := bson.M{
			"$set": bson.M{
				"coins":     item.Coins.Sub(coins...),
				"height":    height,
				"timestamp": timestamp,
				"tx_hash":   txHash,
			},
		}
		projection = bson.M{
			"_id": 1,
		}
		opts := options.FindOneAndUpdate().
			SetProjection(projection).
			SetUpsert(true)

		if _, err := database.DepositFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
			return err
		}

		return nil
	}
}
