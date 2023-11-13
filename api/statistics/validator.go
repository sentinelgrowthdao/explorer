package statistics

import (
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

var (
	validators = map[string]func(req *RequestGetStatistics) error{
		types.StatisticMethodHistoricalActiveNodeCount:         validateHistoricalActiveNodeCount,
		types.StatisticMethodHistoricalActiveSessionCount:      validateHistoricalActiveSessionCount,
		types.StatisticMethodHistoricalActiveSubscriptionCount: validateHistoricalActiveSubscriptionCount,
		types.StatisticMethodHistoricalBytesPayment:            validateHistoricalBytesPayment,
		types.StatisticMethodHistoricalBytesStakingReward:      validateHistoricalBytesStakingReward,
		types.StatisticMethodHistoricalEndSessionCount:         validateHistoricalEndSessionCount,
		types.StatisticMethodHistoricalEndSubscriptionCount:    validateHistoricalEndSubscriptionCount,
		types.StatisticMethodHistoricalHoursPayment:            validateHistoricalHoursPayment,
		types.StatisticMethodHistoricalHoursStakingReward:      validateHistoricalHoursStakingReward,
		types.StatisticMethodHistoricalPlanPayment:             validateHistoricalPlanPayment,
		types.StatisticMethodHistoricalPlanStakingReward:       validateHistoricalPlanStakingReward,
		types.StatisticMethodHistoricalRegisterNodeCount:       validateHistoricalRegisterNodeCount,
		types.StatisticMethodHistoricalSessionAddressCount:     validateHistoricalSessionAddressCount,
		types.StatisticMethodHistoricalSessionBytes:            validateHistoricalSessionBytes,
		types.StatisticMethodHistoricalSessionDuration:         validateHistoricalSessionDuration,
		types.StatisticMethodHistoricalSessionNodeCount:        validateHistoricalSessionNodeCount,
		types.StatisticMethodHistoricalStartSessionCount:       validateHistoricalStartSessionCount,
		types.StatisticMethodHistoricalStartSubscriptionCount:  validateHistoricalStartSubscriptionCount,
		types.StatisticMethodHistoricalSubscriptionDeposit:     validateHistoricalSubscriptionDeposit,
	}
)

func validateHistoricalActiveNodeCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalActiveSessionCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalActiveSubscriptionCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalSessionBytes(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalEndSessionCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalEndSubscriptionCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalRegisterNodeCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalSessionDuration(req *RequestGetStatistics) (err error) {
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

func validateHistoricalBytesPayment(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalBytesStakingReward(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalStartSessionCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalStartSubscriptionCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalSubscriptionDeposit(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalPlanPayment(req *RequestGetStatistics) (err error) {
	allowed := []string{
		"-timestamp",
		"timestamp",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return err
	}

	return nil
}

func validateHistoricalPlanStakingReward(req *RequestGetStatistics) (err error) {
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

func validateHistoricalSessionAddressCount(req *RequestGetStatistics) (err error) {
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

func validateHistoricalSessionNodeCount(req *RequestGetStatistics) (err error) {
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
