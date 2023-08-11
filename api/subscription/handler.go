package subscription

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

func HandlerGetSubscriptions(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetSubscriptions
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
		if req.NodeAddress != "" {
			filter["node_address"] = req.NodeAddress
		}
		if req.ID > 0 {
			filter["plan_id"] = req.ID
		}

		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Skip).
			SetLimit(req.Limit)

		items, err := database.SubscriptionFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetSubscription(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetSubscription
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

		item, err := database.SubscriptionFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetSubscriptionQuotas(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetSubscriptionQuotas
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

		items, err := database.SubscriptionQuotaFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetSubscriptionQuota(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetSubscriptionQuota
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"id":      req.ID,
			"address": req.AccountAddress,
		}
		projection := bson.M{}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.SubscriptionQuotaFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetSubscriptionQuotaEvents(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestGetSubscriptionQuotaEvents
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		if req.Limit == 0 || req.Limit > 100 {
			req.Limit = 10
		}

		filter := bson.M{
			"id":      req.ID,
			"address": req.AccountAddress,
		}
		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Skip).
			SetLimit(req.Limit)

		items, err := database.SubscriptionQuotaEventFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}
