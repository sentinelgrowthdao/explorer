package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	SubscriptionPayoutStatistics struct {
		Timeframe          string
		HoursPayment       types.Coins
		HoursStakingReward types.Coins
	}
)

func NewSubscriptionPayoutStatistics(timeframe string) *SubscriptionPayoutStatistics {
	return &SubscriptionPayoutStatistics{
		Timeframe:          timeframe,
		HoursPayment:       types.NewCoins(nil),
		HoursStakingReward: types.NewCoins(nil),
	}
}

func (sps *SubscriptionPayoutStatistics) Result(timestamp time.Time) bson.A {
	return bson.A{
		bson.M{
			"type":      types.StatisticTypeHoursPayment,
			"timeframe": sps.Timeframe,
			"timestamp": timestamp,
			"value":     sps.HoursPayment,
		},
		bson.M{
			"type":      types.StatisticTypeHoursStakingReward,
			"timeframe": sps.Timeframe,
			"timestamp": timestamp,
			"value":     sps.HoursStakingReward,
		},
	}
}

func StatisticsFromSubscriptionPayouts(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time) (result bson.A, err error) {
	log.Println("StatisticsFromSubscriptionPayouts", minTimestamp, maxTimestamp)

	filter := bson.M{}
	sort := bson.D{
		bson.E{Key: "timestamp", Value: 1},
	}

	items, err := database.SubscriptionPayoutFind(ctx, db, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*SubscriptionPayoutStatistics)
		w = make(map[time.Time]*SubscriptionPayoutStatistics)
		m = make(map[time.Time]*SubscriptionPayoutStatistics)
		y = make(map[time.Time]*SubscriptionPayoutStatistics)
	)

	for i := 0; i < len(items); i++ {
		dayTimestamp := utils.DayDate(items[i].Timestamp)
		if _, ok := d[dayTimestamp]; !ok {
			d[dayTimestamp] = NewSubscriptionPayoutStatistics("day")
		}

		weekTimestamp := utils.ISOWeekDate(items[i].Timestamp)
		if _, ok := w[weekTimestamp]; !ok {
			w[weekTimestamp] = NewSubscriptionPayoutStatistics("week")
		}

		monthTimestamp := utils.MonthDate(items[i].Timestamp)
		if _, ok := m[monthTimestamp]; !ok {
			m[monthTimestamp] = NewSubscriptionPayoutStatistics("month")
		}

		yearTimestamp := utils.YearDate(items[i].Timestamp)
		if _, ok := y[yearTimestamp]; !ok {
			y[yearTimestamp] = NewSubscriptionPayoutStatistics("year")
		}

		d[dayTimestamp].HoursPayment = d[dayTimestamp].HoursPayment.Add(items[i].Payment)
		w[weekTimestamp].HoursPayment = w[weekTimestamp].HoursPayment.Add(items[i].Payment)
		m[monthTimestamp].HoursPayment = m[monthTimestamp].HoursPayment.Add(items[i].Payment)
		y[yearTimestamp].HoursPayment = y[yearTimestamp].HoursPayment.Add(items[i].Payment)

		d[dayTimestamp].HoursStakingReward = d[dayTimestamp].HoursStakingReward.Add(items[i].Payment)
		w[weekTimestamp].HoursStakingReward = w[weekTimestamp].HoursStakingReward.Add(items[i].Payment)
		m[monthTimestamp].HoursStakingReward = m[monthTimestamp].HoursStakingReward.Add(items[i].Payment)
		y[yearTimestamp].HoursStakingReward = y[yearTimestamp].HoursStakingReward.Add(items[i].Payment)
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
