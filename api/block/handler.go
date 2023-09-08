package block

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

func HandlerGetBlocks(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetBlocks(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{}
		if req.Query.FromHeight != 0 {
			filter["height"] = bson.M{
				"$gte": req.Query.FromHeight,
			}
		}
		if req.Query.ToHeight != 0 {
			filter["height"] = bson.M{
				"$lte": req.Query.ToHeight,
			}
		}
		if req.Query.FromHeight != 0 && req.Query.ToHeight != 0 {
			filter["height"] = bson.M{
				"$gte": req.Query.FromHeight,
				"$lte": req.Query.ToHeight,
			}
		}
		projection := bson.M{
			"height":           1,
			"time":             1,
			"proposer_address": 1,
			"num_txs":          1,
			"duration":         1,
		}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.BlockFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetBlock(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetBlock(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"height": req.URI.Height,
		}
		projection := bson.M{
			"height":           1,
			"time":             1,
			"proposer_address": 1,
			"num_txs":          1,
			"duration":         1,
		}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.BlockFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}
