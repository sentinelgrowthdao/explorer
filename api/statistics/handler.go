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
		types.StatisticMethodAverageActiveNodeCount:            handleAverageActiveNodeCount,
		types.StatisticMethodAverageActiveSessionCount:         handleAverageActiveSessionCount,
		types.StatisticMethodAverageActiveSubscriptionCount:    handleAverageActiveSubscriptionCount,
		types.StatisticMethodAverageBytesPayment:               handleAverageBytesPayment,
		types.StatisticMethodAverageBytesStakingReward:         handleAverageBytesStakingReward,
		types.StatisticMethodAverageEndSessionCount:            handleAverageEndSessionCount,
		types.StatisticMethodAverageEndSubscriptionCount:       handleAverageEndSubscriptionCount,
		types.StatisticMethodAveragePlanPayment:                handleAveragePlanPayment,
		types.StatisticMethodAveragePlanStakingReward:          handleAveragePlanStakingReward,
		types.StatisticMethodAverageRegisterNodeCount:          handleAverageRegisterNodeCount,
		types.StatisticMethodAverageStartSessionCount:          handleAverageStartSessionCount,
		types.StatisticMethodAverageStartSubscriptionCount:     handleAverageStartSubscriptionCount,
		types.StatisticMethodAverageSubscriptionDeposit:        handleAverageSubscriptionDeposit,
		types.StatisticMethodCurrentNodeCount:                  handleCurrentNodeCount,
		types.StatisticMethodCurrentSessionAddressCount:        handleCurrentSessionAddressCount,
		types.StatisticMethodCurrentSessionCount:               handleCurrentSessionCount,
		types.StatisticMethodCurrentSessionNodeCount:           handleCurrentSessionNodeCount,
		types.StatisticMethodCurrentSubscriptionCount:          handleCurrentSubscriptionCount,
		types.StatisticMethodHistoricalActiveNodeCount:         handleHistoricalActiveNodeCount,
		types.StatisticMethodHistoricalActiveSessionCount:      handleHistoricalActiveSessionCount,
		types.StatisticMethodHistoricalActiveSubscriptionCount: handleHistoricalActiveSubscriptionCount,
		types.StatisticMethodHistoricalBytesPayment:            handleHistoricalBytesPayment,
		types.StatisticMethodHistoricalBytesStakingReward:      handleHistoricalBytesStakingReward,
		types.StatisticMethodHistoricalEndSessionCount:         handleHistoricalEndSessionCount,
		types.StatisticMethodHistoricalEndSubscriptionCount:    handleHistoricalEndSubscriptionCount,
		types.StatisticMethodHistoricalHoursPayment:            handleHistoricalHoursPayment,
		types.StatisticMethodHistoricalHoursStakingReward:      handleHistoricalHoursStakingReward,
		types.StatisticMethodHistoricalPlanPayment:             handleHistoricalPlanPayment,
		types.StatisticMethodHistoricalPlanStakingReward:       handleHistoricalPlanStakingReward,
		types.StatisticMethodHistoricalRegisterNodeCount:       handleHistoricalRegisterNodeCount,
		types.StatisticMethodHistoricalSessionAddressCount:     handleHistoricalSessionAddressCount,
		types.StatisticMethodHistoricalSessionBytes:            handleHistoricalSessionBytes,
		types.StatisticMethodHistoricalSessionDuration:         handleHistoricalSessionDuration,
		types.StatisticMethodHistoricalSessionNodeCount:        handleHistoricalSessionNodeCount,
		types.StatisticMethodHistoricalStartSessionCount:       handleHistoricalStartSessionCount,
		types.StatisticMethodHistoricalStartSubscriptionCount:  handleHistoricalStartSubscriptionCount,
		types.StatisticMethodHistoricalSubscriptionDeposit:     handleHistoricalSubscriptionDeposit,
		types.StatisticMethodTotalBytesPayment:                 handleTotalBytesPayment,
		types.StatisticMethodTotalBytesStakingReward:           handleTotalBytesStakingReward,
		types.StatisticMethodTotalHoursPayment:                 handleTotalHoursPayment,
		types.StatisticMethodTotalHoursStakingReward:           handleTotalHoursStakingReward,
		types.StatisticMethodTotalPlanPayment:                  handleTotalPlanPayment,
		types.StatisticMethodTotalPlanStakingReward:            handleTotalPlanStakingReward,
		types.StatisticMethodTotalSessionBytes:                 handleTotalSessionBytes,
		types.StatisticMethodTotalSessionDuration:              handleTotalSessionDuration,
		types.StatisticMethodTotalSubscriptionDeposit:          handleTotalSubscriptionDeposit,
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

func handleAverageActiveNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeActiveNode, "$_id", "", "$value", req)
}

func handleAverageActiveSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeActiveSession, "$_id", "", "$value", req)
}

func handleAverageActiveSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeActiveSubscription, "$_id", "", "$value", req)
}

func handleAverageEndSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeEndSession, "$_id", "", "$value", req)
}

func handleAverageEndSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeEndSubscription, "$_id", "", "$value", req)
}

func handleAverageRegisterNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeRegisterNode, "$_id", "", "$value", req)
}

func handleAverageBytesPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeBytesPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageBytesStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeBytesStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleAverageStartSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeStartSession, "$_id", "", "$value", req)
}

func handleAverageStartSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeStartSubscription, "$_id", "", "$value", req)
}

func handleAverageSubscriptionDeposit(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypeSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleAveragePlanPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypePlanPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleAveragePlanStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleAverage(db, types.StatisticTypePlanStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleCurrentNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
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
			"_id":   nil,
			"value": count,
		},
	}, err
}

func handleCurrentSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
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
			"_id":   nil,
			"value": count,
		},
	}, err
}

func handleCurrentSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
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
			"_id":   nil,
			"value": count,
		},
	}, err
}

func handleCurrentSessionAddressCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{}
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

func handleCurrentSessionNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	filter := bson.M{}
	if req.Query.Status != "" {
		filter["status"] = req.Query.Status
	}

	items, err := database.SessionDistinct(context.TODO(), db, "node_addr", filter)
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

func handleHistoricalActiveNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeActiveNode, req)
}

func handleHistoricalActiveSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeActiveSession, req)
}

func handleHistoricalSessionAddressCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionAddress, req)
}

func handleHistoricalActiveSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeActiveSubscription, req)
}

func handleHistoricalSessionBytes(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionBytes, req)
}

func handleHistoricalEndSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeEndSession, req)
}

func handleHistoricalEndSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeEndSubscription, req)
}

func handleHistoricalRegisterNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeRegisterNode, req)
}

func handleHistoricalSessionDuration(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionDuration, req)
}

func handleHistoricalBytesPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeBytesPayment, req)
}

func handleHistoricalHoursStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeHoursStakingReward, req)
}

func handleHistoricalBytesStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeBytesStakingReward, req)
}

func handleHistoricalStartSessionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeStartSession, req)
}

func handleHistoricalSessionNodeCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSessionNode, req)
}

func handleHistoricalStartSubscriptionCount(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeStartSubscription, req)
}

func handleHistoricalSubscriptionDeposit(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeSubscriptionDeposit, req)
}

func handleHistoricalPlanPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypePlanPayment, req)
}

func handleHistoricalHoursPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleHistorical(db, types.StatisticTypeHoursPayment, req)
}

func handleHistoricalPlanStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
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

func handleTotalBytesPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeBytesPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalHoursPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeHoursPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalBytesStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeBytesStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalHoursStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeHoursStakingReward, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalSubscriptionDeposit(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypeSubscriptionDeposit, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalPlanPayment(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypePlanPayment, "$value", "$value.denom", "$value.amount", req)
}

func handleTotalPlanStakingReward(db *mongo.Database, req *RequestGetStatistics) ([]bson.M, error) {
	return handleTotal(db, types.StatisticTypePlanStakingReward, "$value", "$value.denom", "$value.amount", req)
}
