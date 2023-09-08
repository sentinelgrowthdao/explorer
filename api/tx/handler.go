package tx

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

func HandlerGetTxs(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetTxs(c)
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
		if req.URI.Height != 0 {
			filter["height"] = req.URI.Height
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
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.TxFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetTx(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetTx(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"hash": req.URI.Hash,
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
