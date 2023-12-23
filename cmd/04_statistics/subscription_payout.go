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

func (s *SubscriptionPayoutStatistics) Result(timestamp time.Time) []bson.M {
	return []bson.M{
		{
			"type":      types.StatisticTypeHoursPayment,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.HoursPayment,
		},
		{
			"type":      types.StatisticTypeHoursStakingReward,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     s.HoursStakingReward,
		},
	}
}

func StatisticsFromSubscriptionPayouts(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromSubscriptionPayouts")

	filter := bson.M{}
	projection := bson.M{
		"_id":            0,
		"payment":        1,
		"staking_reward": 1,
		"timestamp":      1,
	}

	items, err := database.SubscriptionPayoutFind(ctx, db, filter, options.Find().SetProjection(projection))
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

		d[dayTimestamp].HoursStakingReward = d[dayTimestamp].HoursStakingReward.Add(items[i].StakingReward)
		w[weekTimestamp].HoursStakingReward = w[weekTimestamp].HoursStakingReward.Add(items[i].StakingReward)
		m[monthTimestamp].HoursStakingReward = m[monthTimestamp].HoursStakingReward.Add(items[i].StakingReward)
		y[yearTimestamp].HoursStakingReward = y[yearTimestamp].HoursStakingReward.Add(items[i].StakingReward)
	}

	for t := range d {
		result = append(result, d[t].Result(t)...)
	}
	for t := range w {
		result = append(result, w[t].Result(t)...)
	}
	for t := range m {
		result = append(result, m[t].Result(t)...)
	}
	for t := range y {
		result = append(result, y[t].Result(t)...)
	}

	return result, nil
}
