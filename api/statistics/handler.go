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
		types.StatisticMethodAverageActiveNodeCount:              handleAverageActiveNodesCount,
		types.StatisticMethodAverageActiveSessionCount:           handleAverageActiveSessionsCount,
		types.StatisticMethodAverageActiveSubscriptionCount:      handleAverageActiveSubscriptionsCount,
		types.StatisticMethodAverageEndSessionCount:              handleAverageEndSessionsCount,
		types.StatisticMethodAverageEndSubscriptionCount:         handleAverageEndSubscriptionsCount,
		types.StatisticMethodAverageRegisterNodeCount:            handleAverageRegisterNodesCount,
		types.StatisticMethodAverageSessionPayment:               handleAverageSessionPayments,
		types.StatisticMethodAverageSessionStakingReward:         handleAverageSessionStakingRewards,
		types.StatisticMethodAverageStartSessionCount:            handleAverageStartSessionsCount,
		types.StatisticMethodAverageStartSubscriptionCount:       handleAverageStartSubscriptionsCount,
		types.StatisticMethodAverageSubscriptionDeposit:          handleAverageSubscriptionDeposits,
		types.StatisticMethodAverageSubscriptionPayment:          handleAverageSubscriptionPayments,
		types.StatisticMethodAverageSubscriptionStakingReward:    handleAverageSubscriptionStakingRewards,
		types.StatisticMethodCurrentNodeCount:                    handleCurrentNodesCount,
		types.StatisticMethodCurrentSessionCount:                 handleCurrentSessionsCount,
		types.StatisticMethodCurrentSubscriptionCount:            handleCurrentSubscriptionsCount,
		types.StatisticMethodHistoricalActiveNodeCount:           handleHistoricalActiveNodesCount,
		types.StatisticMethodHistoricalActiveSessionCount:        handleHistoricalActiveSessionsCount,
		types.StatisticMethodHistoricalActiveSubscriptionCount:   handleHistoricalActiveSubscriptionsCount,
		types.StatisticMethodHistoricalEndSessionCount:           handleHistoricalEndSessionsCount,
		types.StatisticMethodHistoricalEndSubscriptionCount:      handleHistoricalEndSubscriptionsCount,
		types.StatisticMethodHistoricalRegisterNodeCount:         handleHistoricalRegisterNodesCount,
		types.StatisticMethodHistoricalSessionBandwidth:          handleHistoricalSessionBandwidths,
		types.StatisticMethodHistoricalSessionDuration:           handleHistoricalSessionDurations,
		types.StatisticMethodHistoricalSessionPayment:            handleHistoricalSessionPayments,
		types.StatisticMethodHistoricalSessionStakingReward:      handleHistoricalSessionStakingRewards,
		types.StatisticMethodHistoricalStartSessionCount:         handleHistoricalStartSessionsCount,
		types.StatisticMethodHistoricalStartSubscriptionCount:    handleHistoricalStartSubscriptionsCount,
		types.StatisticMethodHistoricalSubscriptionDeposit:       handleHistoricalSubscriptionDeposits,
		types.StatisticMethodHistoricalSubscriptionPayment:       handleHistoricalSubscriptionPayments,
		types.StatisticMethodHistoricalSubscriptionStakingReward: handleHistoricalSubscriptionStakingRewards,
		types.StatisticMethodTotalSessionBandwidth:               handleTotalSessionBandwidth,
		types.StatisticMethodTotalSessionDuration:                handleTotalSessionDuration,
		types.StatisticMethodTotalSessionPayment:                 handleTotalSessionPayments,
		types.StatisticMethodTotalSessionStakingReward:           handleTotalSessionStakingRewards,
		types.StatisticMethodTotalSubscriptionDeposit:            handleTotalSubscriptionDeposits,
		types.StatisticMethodTotalSubscriptionPayment:            handleTotalSubscriptionPayments,
		types.StatisticMethodTotalSubscriptionStakingReward:      handleTotalSubscriptionStakingRewards,
	}
)

func HandlerGetStatistics(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := NewRequestGetStatistics(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.NewResponseError(1, err.Error()))
			return
		}

		handlerFunc, ok := requestHandlers[req.Query.Method]
		if !ok {
			err := fmt.Errorf("unknown method %s", req.Query.Method)
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

func handleAverage(db *mongo.Database, t, unwind, _id, value string, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":      t,
				"timeframe": req.Query.Timeframe,
				"timestamp": bson.M{
					"$gte": req.Query.FromTimestamp,
					"$lt":  req.Query.ToTimestamp,
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

	return database.StatisticAggregate(context.TODO(), db, pipeline)
}

func handleHistorical(db *mongo.Database, t string, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{
		"type":      t,
		"timeframe": req.Query.Timeframe,
		"timestamp": bson.M{
			"$gte": req.Query.FromTimestamp,
			"$lt":  req.Query.ToTimestamp,
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
		SetSkip(req.Query.Skip).
		SetLimit(req.Query.Limit)

	return database.StatisticFind(context.TODO(), db, filter, opts)
}

func handleTotal(db *mongo.Database, t, unwind string, _id, value string, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":      t,
				"timeframe": req.Query.Timeframe,
				"timestamp": bson.M{
					"$gte": req.Query.FromTimestamp,
					"$lt":  req.Query.ToTimestamp,
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

	return database.StatisticAggregate(context.TODO(), db, pipeline)
}

func handleAverageActiveNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeActiveNode, "$_id", "", "$value", req)
}

func handleAverageActiveSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeActiveSession, "$_id", "", "$value", req)
}

func handleAverageActiveSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeActiveSubscription, "$_id", "", "$value", req)
}

func handleAverageEndSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeEndSession, "$_id", "", "$value", req)
}

func handleAverageEndSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeEndSubscription, "$_id", "", "$value", req)
}

func handleAverageRegisterNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeRegisterNode, "$_id", "", "$value", req)
}

func handleAverageSessionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeSessionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageSessionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeSessionStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageStartSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeStartSession, "$_id", "", "$value", req)
}

func handleAverageStartSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeStartSubscription, "$_id", "", "$value", req)
}

func handleAverageSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageSubscriptionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeSubscriptionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageSubscriptionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeSubscriptionStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleCurrentNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{}
	if req.Query.Status != "" {
		filter["status"] = req.Query.Status
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
	if req.Query.Status != "" {
		filter["status"] = req.Query.Status
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
	if req.Query.Status != "" {
		filter["status"] = req.Query.Status
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
	return handleHistorical(db, types.StatisticTypeActiveNode, req)
}

func handleHistoricalActiveSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeActiveSession, req)
}

func handleHistoricalActiveSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeActiveSubscription, req)
}

func handleHistoricalSessionBandwidths(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionBandwidth, req)
}

func handleHistoricalEndSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeEndSession, req)
}

func handleHistoricalEndSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeEndSubscription, req)
}

func handleHistoricalRegisterNodesCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeRegisterNode, req)
}

func handleHistoricalSessionDurations(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionDuration, req)
}

func handleHistoricalSessionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionPayment, req)
}

func handleHistoricalSessionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionStakingReward, req)
}

func handleHistoricalStartSessionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeStartSession, req)
}

func handleHistoricalStartSubscriptionsCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeStartSubscription, req)
}

func handleHistoricalSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSubscriptionDeposit, req)
}

func handleHistoricalSubscriptionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSubscriptionPayment, req)
}

func handleHistoricalSubscriptionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSubscriptionStakingReward, req)
}

func handleTotalSessionBandwidth(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":      types.StatisticTypeSessionBandwidth,
				"timeframe": req.Query.Timeframe,
				"timestamp": bson.M{
					"$gte": req.Query.FromTimestamp,
					"$lt":  req.Query.ToTimestamp,
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

	return database.StatisticAggregate(context.TODO(), db, pipeline)
}

func handleTotalSessionDuration(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSessionDuration, "$_id", "", "$value", req)
}

func handleTotalSessionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSessionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSessionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSessionStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSubscriptionPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSubscriptionStakingReward, "$value", "$value.denom", "$value.amount", req)
}
