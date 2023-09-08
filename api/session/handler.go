package session

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

func HandlerGetSessions(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetSessions(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{}
		if req.Query.Status != "" {
			filter["status"] = req.Query.Status
		}
		if req.URI.AccAddr != "" {
			filter["address"] = req.URI.AccAddr
		}
		if req.URI.ID != 0 {
			filter["subscription_id"] = req.URI.ID
		}
		if req.URI.NodeAddr != "" {
			filter["node_address"] = req.URI.NodeAddr
		}

		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.SessionFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetSession(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetSession(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"id": req.URI.ID,
		}
		projection := bson.M{}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.SessionFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetSessionEvents(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetSessionEvents(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"type": bson.M{
				"$in": bson.A{
					types.EventTypeSessionUpdateDetails,
					types.EventTypeSessionUpdateStatus,
				},
			},
			"session_id": req.URI.ID,
		}
		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
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
