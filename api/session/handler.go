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
		var req RequestGetSessions
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		if req.Limit == 0 || req.Limit > 100 {
			req.Limit = 10
		}

		filter := bson.M{}
		if req.Status != "" {
			filter["status"] = req.Status
		}
		if req.AccountAddress != "" {
			filter["address"] = req.AccountAddress
		}
		if req.ID > 0 {
			filter["subscription_id"] = req.ID
		}
		if req.NodeAddress != "" {
			filter["node_address"] = req.NodeAddress
		}

		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Skip).
			SetLimit(req.Limit)

		items, err := database.SessionFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetSession(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetSession
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"id": req.ID,
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
		var req RequestGetSessionEvents
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		if req.Limit == 0 || req.Limit > 100 {
			req.Limit = 10
		}

		filter := bson.M{
			"id": req.ID,
		}
		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Skip).
			SetLimit(req.Limit)

		items, err := database.SessionEventFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}
