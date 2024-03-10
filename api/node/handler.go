package node

import (
	"context"
	"fmt"
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
			"_id":             0,
			"addr":            1,
			"gigabyte_prices": 1,
			"handshake_dns":   1,
			"hourly_prices":   1,
			"internet_speed":  1,
			"location":        1,
			"moniker":         1,
			"peers":           1,
			"type":            1,
			"version":         1,
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
			"addr": req.URI.NodeAddr,
		}
		projection := bson.M{
			"_id":             0,
			"addr":            1,
			"gigabyte_prices": 1,
			"handshake_dns":   1,
			"hourly_prices":   1,
			"internet_speed":  1,
			"location":        1,
			"moniker":         1,
			"peers":           1,
			"type":            1,
			"version":         1,
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

func HandlerGetNodeStatistics(db *mongo.Database, excludeAddrs []string) gin.HandlerFunc {
	requestHandlers := map[string]func(db *mongo.Database, req *RequestGetNodeStatistics) ([]bson.M, error){
		"": handleHistorical,
		types.StatisticMethodCurrentSessionAddressCount: handleCurrentSessionAddressCount(excludeAddrs),
	}

	return func(c *gin.Context) {
		req, err := NewRequestGetNodeStatistics(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		hFunc, ok := requestHandlers[req.Query.Method]
		if !ok {
			err := fmt.Errorf("unknown method %s", req.Query.Method)
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		result, err := hFunc(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(result))
	}
}

func handleHistorical(db *mongo.Database, req *RequestGetNodeStatistics) ([]bson.M, error) {
	filter := bson.M{
		"addr":      req.URI.NodeAddr,
		"timeframe": req.Query.Timeframe,
		"timestamp": bson.M{
			"$gte": req.Query.FromTimestamp,
			"$lt":  req.Query.ToTimestamp,
		},
	}
	projection := bson.M{
		"active_session":      1,
		"active_subscription": 1,
		"addr":                1,
		"bytes_earning":       1,
		"hours_earning":       1,
		"_id":                 0,
		"session_address":     1,
		"session_bandwidth":   1,
		"timeframe":           1,
		"timestamp":           1,
	}
	opts := options.Find().
		SetProjection(projection).
		SetSort(req.Sort).
		SetSkip(req.Query.Skip).
		SetLimit(req.Query.Limit)

	result, err := database.NodeStatisticFind(context.TODO(), db, filter, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handleCurrentSessionAddressCount(excludeAddrs []string) func(*mongo.Database, *RequestGetNodeStatistics) ([]bson.M, error) {
	return func(db *mongo.Database, req *RequestGetNodeStatistics) ([]bson.M, error) {
		filter := bson.M{
			"acc_addr": bson.M{
				"$nin": excludeAddrs,
			},
			"node_addr": req.URI.NodeAddr,
		}
		if req.Query.Status != "" {
			filter["status"] = req.Query.Status
		}

		items, err := database.SessionDistinct(context.TODO(), db, "acc_addr", filter)
		if err != nil {
			return nil, err
		}

		return []bson.M{
			{
				"_id":   nil,
				"value": len(items),
			},
		}, err
	}
}
