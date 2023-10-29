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
		types.StatisticMethodAverageActiveNodeCount:            handleAverageActiveNodesCount,
		types.StatisticMethodAverageActiveSessionCount:         handleAverageActiveSessionsCount,
		types.StatisticMethodAverageActiveSubscriptionCount:    handleAverageActiveSubscriptionsCount,
		types.StatisticMethodAverageEndSessionCount:            handleAverageEndSessionsCount,
		types.StatisticMethodAverageEndSubscriptionCount:       handleAverageEndSubscriptionsCount,
		types.StatisticMethodAverageRegisterNodeCount:          handleAverageRegisterNodesCount,
		types.StatisticMethodAverageBytesPayment:               handleAverageBytesPayments,
		types.StatisticMethodAverageBytesStakingReward:         handleAverageBytesStakingRewards,
		types.StatisticMethodAverageStartSessionCount:          handleAverageStartSessionsCount,
		types.StatisticMethodAverageStartSubscriptionCount:     handleAverageStartSubscriptionsCount,
		types.StatisticMethodAverageSubscriptionDeposit:        handleAverageSubscriptionDeposits,
		types.StatisticMethodAveragePlanPayment:                handleAveragePlanPayments,
		types.StatisticMethodAveragePlanStakingReward:          handleAveragePlanStakingRewards,
		types.StatisticMethodCurrentNodeCount:                  handleCurrentNodesCount,
		types.StatisticMethodCurrentSessionCount:               handleCurrentSessionsCount,
		types.StatisticMethodCurrentSubscriptionCount:          handleCurrentSubscriptionsCount,
		types.StatisticMethodHistoricalActiveNodeCount:         handleHistoricalActiveNodesCount,
		types.StatisticMethodHistoricalActiveSessionCount:      handleHistoricalActiveSessionsCount,
		types.StatisticMethodHistoricalActiveSubscriptionCount: handleHistoricalActiveSubscriptionsCount,
		types.StatisticMethodHistoricalEndSessionCount:         handleHistoricalEndSessionsCount,
		types.StatisticMethodHistoricalEndSubscriptionCount:    handleHistoricalEndSubscriptionsCount,
		types.StatisticMethodHistoricalRegisterNodeCount:       handleHistoricalRegisterNodesCount,
		types.StatisticMethodHistoricalSessionBytes:            handleHistoricalSessionBytess,
		types.StatisticMethodHistoricalSessionDuration:         handleHistoricalSessionDurations,
		types.StatisticMethodHistoricalBytesPayment:            handleHistoricalBytesPayments,
		types.StatisticMethodHistoricalBytesStakingReward:      handleHistoricalBytesStakingRewards,
		types.StatisticMethodHistoricalStartSessionCount:       handleHistoricalStartSessionsCount,
		types.StatisticMethodHistoricalStartSubscriptionCount:  handleHistoricalStartSubscriptionsCount,
		types.StatisticMethodHistoricalSubscriptionDeposit:     handleHistoricalSubscriptionDeposits,
		types.StatisticMethodHistoricalPlanPayment:             handleHistoricalPlanPayments,
		types.StatisticMethodHistoricalPlanStakingReward:       handleHistoricalPlanStakingRewards,
		types.StatisticMethodTotalSessionBytes:                 handleTotalSessionBytes,
		types.StatisticMethodTotalSessionDuration:              handleTotalSessionDuration,
		types.StatisticMethodTotalBytesPayment:                 handleTotalBytesPayments,
		types.StatisticMethodTotalBytesStakingReward:           handleTotalBytesStakingRewards,
		types.StatisticMethodTotalSubscriptionDeposit:          handleTotalSubscriptionDeposits,
		types.StatisticMethodTotalPlanPayment:                  handleTotalPlanPayments,
		types.StatisticMethodTotalPlanStakingReward:            handleTotalPlanStakingRewards,
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

func handleTotal(db *mongo.Database, t, unwind, _id, value interface{}, req *RequestGetStatistics) ([]bson.M, error) {
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
	}

	if unwind != nil {
		pipeline = append(
			pipeline,
			bson.M{
				"$unwind": unwind,
			},
		)
	}

	pipeline = append(
		pipeline,
		[]bson.M{
			{
				"$group": bson.M{
					"_id": _id,
					"value": bson.M{
						"$sum": bson.M{
							"$toLong": value,
						},
					},
				},
			},
			{
				"$sort": bson.D{
					bson.E{Key: "_id", Value: 1},
				},
			},
		}...,
	)

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

func handleAverageBytesPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeBytesPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageBytesStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeBytesStakingReward, "$value", "$value.denom", "$value.amount", req)
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

func handleAveragePlanPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypePlanPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAveragePlanStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypePlanStakingReward, "$value", "$value.denom", "$value.amount", req)
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

func handleHistoricalSessionBytess(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionBytes, req)
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

func handleHistoricalBytesPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeBytesPayment, req)
}

func handleHistoricalBytesStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeBytesStakingReward, req)
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

func handleHistoricalPlanPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypePlanPayment, req)
}

func handleHistoricalPlanStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypePlanStakingReward, req)
}

func handleTotalSessionBytes(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":      types.StatisticTypeSessionBytes,
				"timeframe": req.Query.Timeframe,
				"timestamp": bson.M{
					"$gte": req.Query.FromTimestamp,
					"$lt":  req.Query.ToTimestamp,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"download": bson.M{
					"$sum": bson.M{
						"$toLong": "$value.download",
					},
				},
				"upload": bson.M{
					"$sum": bson.M{
						"$toLong": "$value.upload",
					},
				},
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
	return handleTotal(db, types.StatisticTypeSessionDuration, nil, nil, "$value", req)
}

func handleTotalBytesPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeBytesPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalBytesStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeBytesStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionDeposits(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalPlanPayments(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypePlanPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalPlanStakingRewards(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypePlanStakingReward, "$value", "$value.denom", "$value.amount", req)
}
