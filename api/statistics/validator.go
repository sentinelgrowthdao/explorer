package statistics

import (
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

var (
	validators = map[string]func(req *RequestGetStatistics) error{
		types.StatisticMethodAverageActiveNodeCount:              nil,
		types.StatisticMethodAverageActiveSessionCount:           nil,
		types.StatisticMethodAverageActiveSubscriptionCount:      nil,
		types.StatisticMethodAverageEndSessionCount:              nil,
		types.StatisticMethodAverageEndSubscriptionCount:         nil,
		types.StatisticMethodAverageRegisterNodeCount:            nil,
		types.StatisticMethodAverageSessionPayment:               nil,
		types.StatisticMethodAverageSessionStakingReward:         nil,
		types.StatisticMethodAverageStartSessionCount:            nil,
		types.StatisticMethodAverageStartSubscriptionCount:       nil,
		types.StatisticMethodAverageSubscriptionDeposit:          nil,
		types.StatisticMethodAverageSubscriptionPayment:          nil,
		types.StatisticMethodAverageSubscriptionStakingReward:    nil,
		types.StatisticMethodCurrentNodeCount:                    nil,
		types.StatisticMethodCurrentSessionCount:                 nil,
		types.StatisticMethodCurrentSubscriptionCount:            nil,
		types.StatisticMethodHistoricalActiveNodeCount:           validateHistoricalActiveNodesCount,
		types.StatisticMethodHistoricalActiveSessionCount:        validateHistoricalActiveSessionsCount,
		types.StatisticMethodHistoricalActiveSubscriptionCount:   validateHistoricalActiveSubscriptionsCount,
		types.StatisticMethodHistoricalEndSessionCount:           validateHistoricalEndSessionsCount,
		types.StatisticMethodHistoricalEndSubscriptionCount:      validateHistoricalEndSubscriptionsCount,
		types.StatisticMethodHistoricalRegisterNodeCount:         validateHistoricalRegisterNodesCount,
		types.StatisticMethodHistoricalSessionBandwidth:          validateHistoricalSessionBandwidths,
		types.StatisticMethodHistoricalSessionDuration:           validateHistoricalSessionDurations,
		types.StatisticMethodHistoricalSessionPayment:            validateHistoricalSessionPayments,
		types.StatisticMethodHistoricalSessionStakingReward:      validateHistoricalSessionStakingRewards,
		types.StatisticMethodHistoricalStartSessionCount:         validateHistoricalStartSessionsCount,
		types.StatisticMethodHistoricalStartSubscriptionCount:    validateHistoricalStartSubscriptionsCount,
		types.StatisticMethodHistoricalSubscriptionDeposit:       validateHistoricalSubscriptionDeposits,
		types.StatisticMethodHistoricalSubscriptionPayment:       validateHistoricalSubscriptionPayments,
		types.StatisticMethodHistoricalSubscriptionStakingReward: validateHistoricalSubscriptionStakingRewards,
		types.StatisticMethodTotalSessionBandwidth:               nil,
		types.StatisticMethodTotalSessionDuration:                nil,
		types.StatisticMethodTotalSessionPayment:                 nil,
		types.StatisticMethodTotalSessionStakingReward:           nil,
		types.StatisticMethodTotalSubscriptionDeposit:            nil,
		types.StatisticMethodTotalSubscriptionPayment:            nil,
		types.StatisticMethodTotalSubscriptionStakingReward:      nil,
	}
)

func validateHistoricalActiveNodesCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalActiveSessionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalActiveSubscriptionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionBandwidths(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalEndSessionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalEndSubscriptionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalRegisterNodesCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionDurations(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionPayments(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSessionStakingRewards(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalStartSessionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalStartSubscriptionsCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
		"-value",
		"value",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSubscriptionDeposits(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSubscriptionPayments(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalSubscriptionStakingRewards(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}
