package node

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

func HandlerGetNodes(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetNodes(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{}
		if req.Query.Status != "" {
			filter["status"] = req.Query.Status
		}

		projection := bson.M{
			"_id":            0,
			"addr":           1,
			"handshake_dns":  1,
			"internet_speed": 1,
			"location":       1,
			"moniker":        1,
			"type":           1,
			"version":        1,
		}
		opts := options.Find().
			SetProjection(projection).
			SetSort(req.Sort).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.NodeFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetNode(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetNode(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"address": req.URI.NodeAddr,
		}
		projection := bson.M{
			"_id":            0,
			"addr":           1,
			"handshake_dns":  1,
			"internet_speed": 1,
			"location":       1,
			"moniker":        1,
			"type":           1,
			"version":        1,
		}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.NodeFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetNodeEvents(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetNodeEvents(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"type": bson.M{
				"$in": bson.A{
					types.EventTypeNodeUpdateDetails,
					types.EventTypeNodeUpdateStatus,
				},
			},
			"node_address": req.URI.NodeAddr,
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
