package deposit

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
)

func HandlerGetDeposits(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetDeposits(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{}
		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSort(req.Sort).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.DepositFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetDeposit(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetDeposit(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"address": req.URI.AccAddr,
		}
		projection := bson.M{}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.DepositFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetDepositEvents(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetDepositEvents(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"type": bson.M{
				"$in": bson.A{
					types.EventTypeDepositAdd,
					types.EventTypeDepositSubtract,
				},
			},
			"acc_address": req.URI.AccAddr,
		}
		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSort(req.Sort).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.EventFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}
