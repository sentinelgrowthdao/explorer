package operations

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewSubscriptionPayoutCreate(
	db *mongo.Database,
	v *models.SubscriptionPayout,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.SubscriptionPayoutInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}
