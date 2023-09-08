package operations

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/models"
	"github.com/sentinel-official/explorer/types"
)

func NewEventSaveOperation(
	db *mongo.Database,
	v *models.Event,
) types.DatabaseOperation {
	return func(ctx mongo.SessionContext) error {
		if _, err := database.EventInsertOne(ctx, db, v); err != nil {
			return err
		}

		return nil
	}
}
