package statistics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type RequestGetStatistics struct {
	Method string `form:"method"`

	FromTimestamp time.Time `form:"from_timestamp"`
	ToTimestamp   time.Time `form:"to_timestamp,default=2050-01-01T00:00:00Z" binding:"gtfield=FromTimestamp"`
	Timeframe     string    `form:"timeframe,default=day" binding:"oneof=day week month year"`
	Status        string    `form:"status" binding:"omitempty,oneof=STATUS_ACTIVE STATUS_INACTIVE STATUS_INACTIVE_PENDING"`
	Sort          bson.D
	SortBy        string `form:"sort_by"`
	Skip          int64  `form:"skip,default=0" binding:"gte=0"`
	Limit         int64  `form:"limit,default=30" binding:"gte=0,lte=100"`
}

func NewRequestGetStatistics(c *gin.Context) (req *RequestGetStatistics, err error) {
	req = &RequestGetStatistics{}
	if err = c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if req.Method == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}

	validatorFunc, ok := requestValidators[req.Method]
	if !ok {
		return nil, fmt.Errorf("unknown method %s", req.Method)
	}

	if err = validatorFunc(req); err != nil {
		return nil, err
	}

	return req, nil
}

var (
	requestValidators = map[string]func(req *RequestGetStatistics) error{
		types.MethodAverageActiveNodeCount:              validateAverageActiveNodesCount,
		types.MethodAverageActiveSessionCount:           validateAverageActiveSessionsCount,
		types.MethodAverageActiveSubscriptionCount:      validateAverageActiveSubscriptionsCount,
		types.MethodAverageEndSessionCount:              validateAverageEndSessionsCount,
		types.MethodAverageEndSubscriptionCount:         validateAverageEndSubscriptionsCount,
		types.MethodAverageJoinNodeCount:                validateAverageJoinNodesCount,
		types.MethodAverageSessionPayment:               validateAverageSessionPayments,
		types.MethodAverageSessionStakingReward:         validateAverageSessionStakingRewards,
		types.MethodAverageStartSessionCount:            validateAverageStartSessionsCount,
		types.MethodAverageStartSubscriptionCount:       validateAverageStartSubscriptionsCount,
		types.MethodAverageSubscriptionDeposit:          validateAverageSubscriptionDeposits,
		types.MethodAverageSubscriptionPayment:          validateAverageSubscriptionPayments,
		types.MethodAverageSubscriptionStakingReward:    validateAverageSubscriptionStakingRewards,
		types.MethodCurrentNodeCount:                    validateCurrentNodesCount,
		types.MethodCurrentSessionCount:                 validateCurrentSessionsCount,
		types.MethodCurrentSubscriptionCount:            validateCurrentSubscriptionsCount,
		types.MethodHistoricalActiveNodeCount:           validateHistoricalActiveNodesCount,
		types.MethodHistoricalActiveSessionCount:        validateHistoricalActiveSessionsCount,
		types.MethodHistoricalActiveSubscriptionCount:   validateHistoricalActiveSubscriptionsCount,
		types.MethodHistoricalBandwidthConsumption:      validateHistoricalBandwidthConsumptions,
		types.MethodHistoricalEndSessionCount:           validateHistoricalEndSessionsCount,
		types.MethodHistoricalEndSubscriptionCount:      validateHistoricalEndSubscriptionsCount,
		types.MethodHistoricalJoinNodeCount:             validateHistoricalJoinNodesCount,
		types.MethodHistoricalSessionDuration:           validateHistoricalSessionDurations,
		types.MethodHistoricalSessionPayment:            validateHistoricalSessionPayments,
		types.MethodHistoricalSessionStakingReward:      validateHistoricalSessionStakingRewards,
		types.MethodHistoricalStartSessionCount:         validateHistoricalStartSessionsCount,
		types.MethodHistoricalStartSubscriptionCount:    validateHistoricalStartSubscriptionsCount,
		types.MethodHistoricalSubscriptionDeposit:       validateHistoricalSubscriptionDeposits,
		types.MethodHistoricalSubscriptionPayment:       validateHistoricalSubscriptionPayments,
		types.MethodHistoricalSubscriptionStakingReward: validateHistoricalSubscriptionStakingRewards,
		types.MethodTotalBandwidthConsumption:           validateTotalBandwidthConsumption,
		types.MethodTotalSessionDuration:                validateTotalSessionDuration,
		types.MethodTotalSessionPayment:                 validateTotalSessionPayments,
		types.MethodTotalSessionStakingReward:           validateTotalSessionStakingRewards,
		types.MethodTotalSubscriptionDeposit:            validateTotalSubscriptionDeposits,
		types.MethodTotalSubscriptionPayment:            validateTotalSubscriptionPayments,
		types.MethodTotalSubscriptionStakingReward:      validateTotalSubscriptionStakingRewards,
	}
)

func validateAverageActiveNodesCount(_ *RequestGetStatistics) error           { return nil }
func validateAverageActiveSessionsCount(_ *RequestGetStatistics) error        { return nil }
func validateAverageActiveSubscriptionsCount(_ *RequestGetStatistics) error   { return nil }
func validateAverageEndSessionsCount(_ *RequestGetStatistics) error           { return nil }
func validateAverageEndSubscriptionsCount(_ *RequestGetStatistics) error      { return nil }
func validateAverageJoinNodesCount(_ *RequestGetStatistics) error             { return nil }
func validateAverageSessionPayments(_ *RequestGetStatistics) error            { return nil }
func validateAverageSessionStakingRewards(_ *RequestGetStatistics) error      { return nil }
func validateAverageStartSessionsCount(_ *RequestGetStatistics) error         { return nil }
func validateAverageStartSubscriptionsCount(_ *RequestGetStatistics) error    { return nil }
func validateAverageSubscriptionDeposits(_ *RequestGetStatistics) error       { return nil }
func validateAverageSubscriptionPayments(_ *RequestGetStatistics) error       { return nil }
func validateAverageSubscriptionStakingRewards(_ *RequestGetStatistics) error { return nil }
func validateCurrentNodesCount(_ *RequestGetStatistics) error                 { return nil }
func validateCurrentSessionsCount(_ *RequestGetStatistics) error              { return nil }
func validateCurrentSubscriptionsCount(_ *RequestGetStatistics) error         { return nil }

func validateHistoricalActiveNodesCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalActiveSessionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalActiveSubscriptionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalBandwidthConsumptions(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalEndSessionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalEndSubscriptionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalJoinNodesCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionDurations(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionPayments(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionStakingRewards(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalStartSessionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalStartSubscriptionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"value",
		"-value",
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSubscriptionDeposits(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSubscriptionPayments(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSubscriptionStakingRewards(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"timestamp",
		"-timestamp",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return err
	}

	return nil
}

func validateTotalBandwidthConsumption(_ *RequestGetStatistics) error       { return nil }
func validateTotalSessionDuration(_ *RequestGetStatistics) error            { return nil }
func validateTotalSessionPayments(_ *RequestGetStatistics) error            { return nil }
func validateTotalSessionStakingRewards(_ *RequestGetStatistics) error      { return nil }
func validateTotalSubscriptionDeposits(_ *RequestGetStatistics) error       { return nil }
func validateTotalSubscriptionPayments(_ *RequestGetStatistics) error       { return nil }
func validateTotalSubscriptionStakingRewards(_ *RequestGetStatistics) error { return nil }
