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
		req, err := NewRequestGetSubscriptions(c)
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
			filter["plan_id"] = req.URI.ID
		}
		if req.URI.NodeAddr != "" {
			filter["node_address"] = req.URI.NodeAddr
		}

		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.SubscriptionFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetSubscription(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetSubscription(c)
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

		item, err := database.SubscriptionFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetSubscriptionEvents(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetSubscriptionEvents(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"type": bson.M{
				"$in": bson.A{
					types.EventTypeSubscriptionUpdateDetails,
					types.EventTypeSubscriptionUpdateStatus,
				},
			},
			"subscription_id": req.URI.ID,
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

func HandlerGetAllocations(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetAllocations(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"id": req.URI.ID,
		}
		projection := bson.M{}
		opts := options.Find().
			SetProjection(projection).
			SetSkip(req.Query.Skip).
			SetLimit(req.Query.Limit)

		items, err := database.SubscriptionAllocationFind(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(items))
	}
}

func HandlerGetAllocation(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetAllocation(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"id":      req.URI.ID,
			"address": req.URI.AccAddr,
		}
		projection := bson.M{}
		opts := options.FindOne().
			SetProjection(projection)

		item, err := database.SubscriptionAllocationFindOne(context.TODO(), db, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(item))
	}
}

func HandlerGetAllocationEvents(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetAllocationEvents(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		filter := bson.M{
			"type": bson.M{
				"$in": bson.A{
					types.EventTypeSubscriptionAllocationUpdateDetails,
				},
			},
			"subscription_id": req.URI.ID,
			"acc_address":     req.URI.AccAddr,
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
