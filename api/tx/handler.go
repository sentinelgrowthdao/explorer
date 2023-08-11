package tx

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

func HandlerGetTxs(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetTxs
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}
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

		filter := bson.M{}
		if req.Height > 0 {
			filter = bson.M{
				"height": req.Height,
			}
		} else {
			filter = bson.M{
				"height": bson.M{
					"$gte": req.FromHeight,
					"$lt":  req.ToHeight,
				},
			}
		}

		projection := bson.M{
			"hash":                 1,
			"height":               1,
			"index":                1,
			"signer_infos.address": 1,
			"fee":                  1,
			"gas_limit":            1,
			"payer":                1,
			"result.gas_wanted":    1,
			"result.gas_used":      1,
		}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Skip).
			SetLimit(req.Limit)

		items, err := database.TxFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetTx(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetTx
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"hash": req.Hash,
		}
		projection := bson.M{
			"hash":                 1,
			"height":               1,
			"index":                1,
			"signer_infos.address": 1,
			"fee":                  1,
			"gas_limit":            1,
			"payer":                1,
			"memo":                 1,
			"result.code":          1,
			"result.codespace":     1,
			"result.gas_wanted":    1,
			"result.gas_used":      1,
		}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.TxFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}
