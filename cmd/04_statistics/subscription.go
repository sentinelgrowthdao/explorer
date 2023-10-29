package main

import (
	"context"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	SubscriptionStatistics struct {
		Timeframe           string
		ActiveSubscription  int64
		BytesSubscription   int64
		EndSubscription     int64
		HoursSubscription   int64
		PlanPayment         types.Coins
		PlanStakingReward   types.Coins
		PlanSubscription    int64
		StartSubscription   int64
		SubscriptionBytes   string
		SubscriptionDeposit types.Coins
		SubscriptionHours   int64
		SubscriptionRefund  types.Coins
	}
)

func NewSubscriptionStatistics(timeframe string) *SubscriptionStatistics {
	return &SubscriptionStatistics{
		Timeframe:           timeframe,
		PlanPayment:         types.NewCoins(nil),
		PlanStakingReward:   types.NewCoins(nil),
		SubscriptionDeposit: types.NewCoins(nil),
		SubscriptionRefund:  types.NewCoins(nil),
	}
}

func (s *SubscriptionStatistics) Result(timestamp time.Time) []bson.M {
	return []bson.M{
		{
			"type":      types.StatisticTypeActiveSubscription,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.ActiveSubscription,
		},
		{
			"type":      types.StatisticTypeBytesSubscription,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.BytesSubscription,
		},
		{
			"type":      types.StatisticTypeEndSubscription,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.EndSubscription,
		},
		{
			"type":      types.StatisticTypeHoursSubscription,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.HoursSubscription,
		},
		{
			"type":      types.StatisticTypePlanPayment,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.PlanPayment,
		},
		{
			"type":      types.StatisticTypePlanStakingReward,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.PlanStakingReward,
		},
		{
			"type":      types.StatisticTypePlanSubscription,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.PlanSubscription,
		},
		{
			"type":      types.StatisticTypeStartSubscription,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.StartSubscription,
		},
		{
			"type":      types.StatisticTypeSubscriptionBytes,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.SubscriptionBytes,
		},
		{
			"type":      types.StatisticTypeSubscriptionDeposit,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.SubscriptionDeposit,
		},
		{
			"type":      types.StatisticTypeSubscriptionHours,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.SubscriptionHours,
		},
		{
			"type":      types.StatisticTypeSubscriptionRefund,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.SubscriptionRefund,
		},
	}
}

func StatisticsFromSubscriptions(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time) (result []bson.M, err error) {
	log.Println("StatisticsFromSubscriptions", minTimestamp, maxTimestamp)

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"end_timestamp":   1,
		"deposit":         1,
		"gigabytes":       1,
		"hours":           1,
		"payment":         1,
		"plan_id":         1,
		"refund":          1,
		"staking_reward":  1,
		"start_timestamp": 1,
	}

	items, err := database.SubscriptionFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*SubscriptionStatistics)
		w = make(map[time.Time]*SubscriptionStatistics)
		m = make(map[time.Time]*SubscriptionStatistics)
		y = make(map[time.Time]*SubscriptionStatistics)
	)

	for i := 0; i < len(items); i++ {
		startTimestamp := items[i].StartTimestamp
		if items[i].StartTimestamp.IsZero() {
			startTimestamp = minTimestamp
		}

		endTimestamp := items[i].EndTimestamp
		if items[i].EndTimestamp.IsZero() {
			endTimestamp = maxTimestamp
		}

		dayStartTimestamp, dayEndTimestamp := utils.DayDate(startTimestamp), utils.DayDate(endTimestamp)
		for t := dayStartTimestamp; !t.After(dayEndTimestamp); t = t.AddDate(0, 0, 1) {
			if _, ok := d[t]; !ok {
				d[t] = NewSubscriptionStatistics("day")
			}

			d[t].ActiveSubscription += 1
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(startTimestamp), utils.ISOWeekDate(endTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			if _, ok := w[t]; !ok {
				w[t] = NewSubscriptionStatistics("week")
			}

			w[t].ActiveSubscription += 1
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(startTimestamp), utils.MonthDate(endTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			if _, ok := m[t]; !ok {
				m[t] = NewSubscriptionStatistics("month")
			}

			m[t].ActiveSubscription += 1
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(startTimestamp), utils.YearDate(endTimestamp)
		for t := yearStartTimestamp; !t.After(yearEndTimestamp); t = t.AddDate(1, 0, 0) {
			if _, ok := y[t]; !ok {
				y[t] = NewSubscriptionStatistics("year")
			}

			y[t].ActiveSubscription += 1
		}

		if !items[i].EndTimestamp.IsZero() {
			d[dayEndTimestamp].EndSubscription += 1
			w[weekEndTimestamp].EndSubscription += 1
			m[monthEndTimestamp].EndSubscription += 1
			y[yearEndTimestamp].EndSubscription += 1
		}
		if items[i].Deposit != nil {
			d[dayStartTimestamp].SubscriptionDeposit = d[dayStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			w[weekStartTimestamp].SubscriptionDeposit = w[weekStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			m[monthStartTimestamp].SubscriptionDeposit = m[monthStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			y[yearStartTimestamp].SubscriptionDeposit = y[yearStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
		}
		if items[i].Gigabytes != 0 {
			d[dayStartTimestamp].BytesSubscription += 1
			w[weekStartTimestamp].BytesSubscription += 1
			m[monthStartTimestamp].BytesSubscription += 1
			y[yearStartTimestamp].BytesSubscription += 1

			bytes := hubtypes.Gigabyte.MulRaw(items[i].Gigabytes)
			d[dayStartTimestamp].SubscriptionBytes = utils.MustIntFromString(d[dayStartTimestamp].SubscriptionBytes).Add(bytes).String()
			w[weekStartTimestamp].SubscriptionBytes = utils.MustIntFromString(w[weekStartTimestamp].SubscriptionBytes).Add(bytes).String()
			m[monthStartTimestamp].SubscriptionBytes = utils.MustIntFromString(m[monthStartTimestamp].SubscriptionBytes).Add(bytes).String()
			y[yearStartTimestamp].SubscriptionBytes = utils.MustIntFromString(y[yearStartTimestamp].SubscriptionBytes).Add(bytes).String()
		}
		if items[i].Hours != 0 {
			d[dayStartTimestamp].HoursSubscription += 1
			w[weekStartTimestamp].HoursSubscription += 1
			m[monthStartTimestamp].HoursSubscription += 1
			y[yearStartTimestamp].HoursSubscription += 1

			d[dayStartTimestamp].SubscriptionHours += items[i].Hours
			w[weekStartTimestamp].SubscriptionHours += items[i].Hours
			m[monthStartTimestamp].SubscriptionHours += items[i].Hours
			y[yearStartTimestamp].SubscriptionHours += items[i].Hours
		}
		if items[i].Payment != nil {
			d[dayStartTimestamp].PlanPayment = d[dayStartTimestamp].PlanPayment.Add(items[i].Payment)
			w[weekStartTimestamp].PlanPayment = w[weekStartTimestamp].PlanPayment.Add(items[i].Payment)
			m[monthStartTimestamp].PlanPayment = m[monthStartTimestamp].PlanPayment.Add(items[i].Payment)
			y[yearStartTimestamp].PlanPayment = y[yearStartTimestamp].PlanPayment.Add(items[i].Payment)
		}
		if items[i].PlanID != 0 {
			d[dayStartTimestamp].PlanSubscription += 1
			w[weekStartTimestamp].PlanSubscription += 1
			m[monthStartTimestamp].PlanSubscription += 1
			y[yearStartTimestamp].PlanSubscription += 1
		}
		if items[i].Refund != nil {
			d[dayEndTimestamp].SubscriptionRefund = d[dayEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
			w[weekEndTimestamp].SubscriptionRefund = w[weekEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
			m[monthEndTimestamp].SubscriptionRefund = m[monthEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
			y[yearEndTimestamp].SubscriptionRefund = y[yearEndTimestamp].SubscriptionRefund.Add(items[i].Refund)
		}
		if items[i].StakingReward != nil {
			d[dayStartTimestamp].PlanStakingReward = d[dayStartTimestamp].PlanStakingReward.Add(items[i].StakingReward)
			w[weekStartTimestamp].PlanStakingReward = w[weekStartTimestamp].PlanStakingReward.Add(items[i].StakingReward)
			m[monthStartTimestamp].PlanStakingReward = m[monthStartTimestamp].PlanStakingReward.Add(items[i].StakingReward)
			y[yearStartTimestamp].PlanStakingReward = y[yearStartTimestamp].PlanStakingReward.Add(items[i].StakingReward)
		}
		if !items[i].StartTimestamp.IsZero() {
			d[dayStartTimestamp].StartSubscription += 1
			w[weekStartTimestamp].StartSubscription += 1
			m[monthStartTimestamp].StartSubscription += 1
			y[yearStartTimestamp].StartSubscription += 1
		}
	}

	for t, statistics := range d {
		result = append(result, statistics.Result(t)...)
	}
	for t, statistics := range w {
		result = append(result, statistics.Result(t)...)
	}
	for t, statistics := range m {
		result = append(result, statistics.Result(t)...)
	}
	for t, statistics := range y {
		result = append(result, statistics.Result(t)...)
	}

	return result, nil
}
