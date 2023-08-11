package statistics

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

var (
	requestHandlers = map[string]func(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error){
		types.MethodAverageActiveNodeCount:              handleAverageActiveNodesCount,
		types.MethodAverageActiveSessionCount:           handleAverageActiveSessionsCount,
		types.MethodAverageActiveSubscriptionCount:      handleAverageActiveSubscriptionsCount,
		types.MethodAverageEndSessionCount:              handleAverageEndSessionsCount,
		types.MethodAverageEndSubscriptionCount:         handleAverageEndSubscriptionsCount,
		types.MethodAverageJoinNodeCount:                handleAverageJoinNodesCount,
		types.MethodAverageSessionPayment:               handleAverageSessionPayments,
		types.MethodAverageSessionStakingReward:         handleAverageSessionStakingRewards,
		types.MethodAverageStartSessionCount:            handleAverageStartSessionsCount,
		types.MethodAverageStartSubscriptionCount:       handleAverageStartSubscriptionsCount,
		types.MethodAverageSubscriptionDeposit:          handleAverageSubscriptionDeposits,
		types.MethodAverageSubscriptionPayment:          handleAverageSubscriptionPayments,
		types.MethodAverageSubscriptionStakingReward:    handleAverageSubscriptionStakingRewards,
		types.MethodCurrentNodeCount:                    handleCurrentNodesCount,
		types.MethodCurrentSessionCount:                 handleCurrentSessionsCount,
		types.MethodCurrentSubscriptionCount:            handleCurrentSubscriptionsCount,
		types.MethodHistoricalActiveNodeCount:           handleHistoricalActiveNodesCount,
		types.MethodHistoricalActiveSessionCount:        handleHistoricalActiveSessionsCount,
		types.MethodHistoricalActiveSubscriptionCount:   handleHistoricalActiveSubscriptionsCount,
		types.MethodHistoricalBandwidthConsumption:      handleHistoricalBandwidthConsumptions,
		types.MethodHistoricalEndSessionCount:           handleHistoricalEndSessionsCount,
		types.MethodHistoricalEndSubscriptionCount:      handleHistoricalEndSubscriptionsCount,
		types.MethodHistoricalJoinNodeCount:             handleHistoricalJoinNodesCount,
		types.MethodHistoricalSessionDuration:           handleHistoricalSessionDurations,
		types.MethodHistoricalSessionPayment:            handleHistoricalSessionPayments,
		types.MethodHistoricalSessionStakingReward:      handleHistoricalSessionStakingRewards,
		types.MethodHistoricalStartSessionCount:         handleHistoricalStartSessionsCount,
		types.MethodHistoricalStartSubscriptionCount:    handleHistoricalStartSubscriptionsCount,
		types.MethodHistoricalSubscriptionDeposit:       handleHistoricalSubscriptionDeposits,
		types.MethodHistoricalSubscriptionPayment:       handleHistoricalSubscriptionPayments,
		types.MethodHistoricalSubscriptionStakingReward: handleHistoricalSubscriptionStakingRewards,
		types.MethodTotalBandwidthConsumption:           handleTotalBandwidthConsumption,
		types.MethodTotalSessionDuration:                handleTotalSessionDuration,
		types.MethodTotalSessionPayment:                 handleTotalSessionPayments,
		types.MethodTotalSessionStakingReward:           handleTotalSessionStakingRewards,
		types.MethodTotalSubscriptionDeposit:            handleTotalSubscriptionDeposits,
		types.MethodTotalSubscriptionPayment:            handleTotalSubscriptionPayments,
		types.MethodTotalSubscriptionStakingReward:      handleTotalSubscriptionStakingRewards,
	}
)

func HandlerGetStatistics(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetStatistics(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		handlerFunc, ok := requestHandlers[req.Method]
		if !ok {
			err = fmt.Errorf("unknown method %s", req.Method)
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		result, err := handlerFunc(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewResponseError(2, err.Error()))
			return
		}

		c.JSON(http.StatusOK, types.NewResponseResult(result))
	}
}

func handleAverage(db *mongo.Database, tag, unwind, _id, value string, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"tag":       tag,
				"timeframe": req.Timeframe,
				"timestamp": bson.M{
					"$gte": req.FromTimestamp,
					"$lt":  req.ToTimestamp,
				},
			},
		},
		{
			"$unwind": unwind,
		},
		{
			"$group": bson.M{
				"_id": _id,
				"value": bson.M{
					"$avg": value,
				},
			},
		},
		{
			"$sort": bson.D{
				bson.E{Key: "_id", Value: 1},
			},
		},
	}

	return database.StatisticsAggregate(context.TODO(), db, pipeline)
}

func handleHistorical(db *mongo.Database, tag string, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{
		"tag":       tag,
		"timeframe": req.Timeframe,
		"timestamp": bson.M{
			"$gte": req.FromTimestamp,
			"$lt":  req.ToTimestamp,
		},
	}
	projection := bson.M{
		"_id":       0,
		"timestamp": 1,
		"value":     1,
	}
	opts := options.Find().
		SetProjection(projection).
		SetSort(req.Sort).
		SetSkip(req.Skip).
		SetLimit(req.Limit)

	return database.StatisticsFindAll(context.TODO(), db, filter, opts)
}

func handleTotal(db *mongo.Database, tag, unwind string, _id, value string, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"tag":       tag,
				"timeframe": req.Timeframe,
				"timestamp": bson.M{
					"$gte": req.FromTimestamp,
					"$lt":  req.ToTimestamp,
				},
			},
		},
		{
			"$unwind": unwind,
		},
		{
			"$group": bson.M{
				"_id": _id,
				"value": bson.M{
					"$sum": value,
				},
			},
		},
		{
			"$sort": bson.D{
				bson.E{Key: "_id", Value: 1},
			},
		},
	}

	return database.StatisticsAggregate(context.TODO(), db, pipeline)
}

func handleAverageActiveNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagActiveNode, "$_id", "", "$value", req)
}

func handleAverageActiveSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagActiveSession, "$_id", "", "$value", req)
}

func handleAverageActiveSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagActiveSubscription, "$_id", "", "$value", req)
}

func handleAverageEndSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagEndSession, "$_id", "", "$value", req)
}

func handleAverageEndSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagEndSubscription, "$_id", "", "$value", req)
}

func handleAverageJoinNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagJoinNode, "$_id", "", "$value", req)
}

func handleAverageSessionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagSessionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageSessionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagSessionStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageStartSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagStartSession, "$_id", "", "$value", req)
}

func handleAverageStartSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagStartSubscription, "$_id", "", "$value", req)
}

func handleAverageSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageSubscriptionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagSubscriptionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageSubscriptionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.TagSubscriptionStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleCurrentNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	count, err := database.NodeCountDocuments(context.TODO(), db, filter)
	if err != nil {
		return nil, err
	}

	return []bson.M{
		{
			"_id":   "",
			"value": count,
		},
	}, err
}

func handleCurrentSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	count, err := database.SessionCountDocuments(context.TODO(), db, filter)
	if err != nil {
		return nil, err
	}

	return []bson.M{
		{
			"_id":   "",
			"value": count,
		},
	}, err
}

func handleCurrentSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	count, err := database.SubscriptionCountDocuments(context.TODO(), db, filter)
	if err != nil {
		return nil, err
	}

	return []bson.M{
		{
			"_id":   "",
			"value": count,
		},
	}, err
}

func handleHistoricalActiveNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagActiveNode, req)
}

func handleHistoricalActiveSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagActiveSession, req)
}

func handleHistoricalActiveSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagActiveSubscription, req)
}

func handleHistoricalBandwidthConsumptions(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagBandwidthConsumption, req)
}

func handleHistoricalEndSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagEndSession, req)
}

func handleHistoricalEndSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagEndSubscription, req)
}

func handleHistoricalJoinNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagJoinNode, req)
}

func handleHistoricalSessionDurations(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagSessionDuration, req)
}

func handleHistoricalSessionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagSessionPayment, req)
}

func handleHistoricalSessionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagSessionStakingReward, req)
}

func handleHistoricalStartSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagStartSession, req)
}

func handleHistoricalStartSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagStartSubscription, req)
}

func handleHistoricalSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagSubscriptionDeposit, req)
}

func handleHistoricalSubscriptionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagSubscriptionPayment, req)
}

func handleHistoricalSubscriptionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.TagSubscriptionStakingReward, req)
}

func handleTotalBandwidthConsumption(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"tag":       types.TagBandwidthConsumption,
				"timeframe": req.Timeframe,
				"timestamp": bson.M{
					"$gte": req.FromTimestamp,
					"$lt":  req.ToTimestamp,
				},
			},
		},
		{
			"$unwind": "$_id",
		},
		{
			"$group": bson.M{
				"_id": "",
				"download": bson.M{
					"$sum": "$value.download",
				},
				"upload": bson.M{
					"$sum": "$value.upload",
				},
			},
		},
		{
			"$sort": bson.D{
				bson.E{Key: "_id", Value: 1},
			},
		},
		{
			"$project": bson.M{
				"_id": "$_id",
				"value": bson.M{
					"download": "$download",
					"upload":   "$upload",
				},
			},
		},
	}

	return database.StatisticsAggregate(context.TODO(), db, pipeline)
}

func handleTotalSessionDuration(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.TagSessionDuration, "$_id", "", "$value", req)
}

func handleTotalSessionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.TagSessionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSessionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.TagSessionStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.TagSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.TagSubscriptionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.TagSubscriptionStakingReward, "$value", "$value.denom", "$value.amount", req)
}
