package statistics

import (
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

var (
	validators = map[string]func(req *RequestGetStatistics) error{
		types.StatisticMethodHistoricalActiveNodeCount:         validateHistoricalActiveNodesCount,
		types.StatisticMethodHistoricalActiveSessionCount:      validateHistoricalActiveSessionsCount,
		types.StatisticMethodHistoricalActiveSubscriptionCount: validateHistoricalActiveSubscriptionsCount,
		types.StatisticMethodHistoricalEndSessionCount:         validateHistoricalEndSessionsCount,
		types.StatisticMethodHistoricalEndSubscriptionCount:    validateHistoricalEndSubscriptionsCount,
		types.StatisticMethodHistoricalRegisterNodeCount:       validateHistoricalRegisterNodesCount,
		types.StatisticMethodHistoricalSessionBytes:            validateHistoricalSessionBytess,
		types.StatisticMethodHistoricalSessionDuration:         validateHistoricalSessionDurations,
		types.StatisticMethodHistoricalActiveAddressCount:      validateHistoricalActiveAddressCount,
		types.StatisticMethodHistoricalBytesPayment:            validateHistoricalBytesPayments,
		types.StatisticMethodHistoricalHoursStakingReward:      validateHistoricalHoursStakingReward,
		types.StatisticMethodHistoricalBytesStakingReward:      validateHistoricalBytesStakingRewards,
		types.StatisticMethodHistoricalStartSessionCount:       validateHistoricalStartSessionsCount,
		types.StatisticMethodHistoricalStartSubscriptionCount:  validateHistoricalStartSubscriptionsCount,
		types.StatisticMethodHistoricalSubscriptionDeposit:     validateHistoricalSubscriptionDeposits,
		types.StatisticMethodHistoricalPlanPayment:             validateHistoricalPlanPayments,
		types.StatisticMethodHistoricalPlanStakingReward:       validateHistoricalPlanStakingRewards,
		types.StatisticMethodHistoricalHoursPayment:            validateHistoricalHoursPayment,
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

func validateHistoricalSessionBytess(req *RequestGetStatistics) (err error) {
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

func validateHistoricalBytesPayments(req *RequestGetStatistics) (err error) {
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

func validateHistoricalBytesStakingRewards(req *RequestGetStatistics) (err error) {
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

func validateHistoricalPlanPayments(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalPlanStakingRewards(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalHoursPayment(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalActiveAddressCount(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalHoursStakingReward(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}
