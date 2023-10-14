package main

import (
	"context"
	"time"

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
		EndSubscription     int64
		StartSubscription   int64
		SubscriptionDeposit types.Coins
		SubscriptionPayment types.Coins
	}
)

func NewSubscriptionStatistics(timeframe string) *SubscriptionStatistics {
	return &SubscriptionStatistics{
		Timeframe:           timeframe,
		SubscriptionDeposit: types.NewCoins(nil),
		SubscriptionPayment: types.NewCoins(nil),
	}
}

func (ss *SubscriptionStatistics) Result(timestamp time.Time) []bson.M {
	return []bson.M{
		{
			"type":      types.StatisticTypeActiveSubscription,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.ActiveSubscription,
		},
		{
			"type":      types.StatisticTypeEndSubscription,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.EndSubscription,
		},
		{
			"type":      types.StatisticTypeStartSubscription,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.StartSubscription,
		},
		{
			"type":      types.StatisticTypeSubscriptionDeposit,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.SubscriptionDeposit,
		},
		{
			"type":      types.StatisticTypeSubscriptionPayment,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.SubscriptionPayment,
		},
	}
}

func StatisticsFromSubscriptions(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	filter := bson.M{
		"start_timestamp": bson.M{
			"$exists": true,
		},
	}
	projection := bson.M{}
	sort := bson.D{
		bson.E{Key: "start_timestamp", Value: 1},
	}

	items, err := database.SubscriptionFind(ctx, db, filter, options.Find().SetProjection(projection).SetSort(sort))
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
		dayStartTimestamp, dayEndTimestamp := utils.DayDate(items[i].StartTimestamp), utils.DayDate(items[i].EndTimestamp)
		for t := dayStartTimestamp; !t.After(dayEndTimestamp); t = t.AddDate(0, 0, 1) {
			if _, ok := d[t]; !ok {
				d[t] = NewSubscriptionStatistics("day")
			}

			d[t].ActiveSubscription += 1
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(items[i].StartTimestamp), utils.ISOWeekDate(items[i].EndTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			if _, ok := w[t]; !ok {
				w[t] = NewSubscriptionStatistics("week")
			}

			w[t].ActiveSubscription += 1
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(items[i].StartTimestamp), utils.MonthDate(items[i].EndTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			if _, ok := m[t]; !ok {
				m[t] = NewSubscriptionStatistics("month")
			}

			m[t].ActiveSubscription += 1
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(items[i].StartTimestamp), utils.YearDate(items[i].EndTimestamp)
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
		if !items[i].StartTimestamp.IsZero() {
			d[dayStartTimestamp].StartSubscription += 1
			w[weekStartTimestamp].StartSubscription += 1
			m[monthStartTimestamp].StartSubscription += 1
			y[yearStartTimestamp].StartSubscription += 1
		}
		if items[i].Deposit != nil {
			d[dayEndTimestamp].SubscriptionDeposit = d[dayStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			w[weekEndTimestamp].SubscriptionDeposit = w[weekStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			m[monthEndTimestamp].SubscriptionDeposit = m[monthStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
			y[yearEndTimestamp].SubscriptionDeposit = y[yearStartTimestamp].SubscriptionDeposit.Add(items[i].Deposit)
		}
		if items[i].Payment != nil {
			d[dayEndTimestamp].SubscriptionPayment = d[dayStartTimestamp].SubscriptionPayment.Add(items[i].Payment)
			w[weekEndTimestamp].SubscriptionPayment = w[weekStartTimestamp].SubscriptionPayment.Add(items[i].Payment)
			m[monthEndTimestamp].SubscriptionPayment = m[monthStartTimestamp].SubscriptionPayment.Add(items[i].Payment)
			y[yearEndTimestamp].SubscriptionPayment = y[yearStartTimestamp].SubscriptionPayment.Add(items[i].Payment)
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
