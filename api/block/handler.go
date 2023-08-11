package block

import (
	"context"
	"math"
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
		var req RequestGetBlocks
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		if req.ToHeight == 0 {
			req.ToHeight = math.MaxInt64
		}
		if req.Limit == 0 || req.Limit > 100 {
			req.Limit = 10
		}

		filter := bson.M{
			"height": bson.M{
				"$gte": req.FromHeight,
				"$lt":  req.ToHeight,
			},
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
			SetSkip(req.Skip).
			SetLimit(req.Limit)

		items, err := database.BlockFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetBlock(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetBlock
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"height": req.Height,
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
