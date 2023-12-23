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
		Timeframe    string
		HoursEarning types.Coins
	}
)

func NewSubscriptionPayoutStatistics(timeframe string) *SubscriptionPayoutStatistics {
	return &SubscriptionPayoutStatistics{
		Timeframe:    timeframe,
		HoursEarning: types.NewCoins(nil),
	}
}

func (s *SubscriptionPayoutStatistics) Result(addr string, timestamp time.Time) bson.M {
	res := bson.M{
		"addr":      addr,
		"timeframe": s.Timeframe,
		"timestamp": timestamp,
	}

	if s.HoursEarning.Len() != 0 {
		res["hours_earning"] = s.HoursEarning
	}

	return res
}

func StatisticsFromSubscriptionPayouts(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromSubscriptionPayouts")

	filter := bson.M{}
	projection := bson.M{
		"_id":       0,
		"node_addr": 1,
		"payment":   1,
		"timestamp": 1,
	}

	items, err := database.SubscriptionPayoutFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[string]map[time.Time]*SubscriptionPayoutStatistics)
		w = make(map[string]map[time.Time]*SubscriptionPayoutStatistics)
		m = make(map[string]map[time.Time]*SubscriptionPayoutStatistics)
		y = make(map[string]map[time.Time]*SubscriptionPayoutStatistics)
	)

	for i := 0; i < len(items); i++ {
		if _, ok := d[items[i].NodeAddr]; !ok {
			d[items[i].NodeAddr] = make(map[time.Time]*SubscriptionPayoutStatistics)
		}
		if _, ok := w[items[i].NodeAddr]; !ok {
			w[items[i].NodeAddr] = make(map[time.Time]*SubscriptionPayoutStatistics)
		}
		if _, ok := m[items[i].NodeAddr]; !ok {
			m[items[i].NodeAddr] = make(map[time.Time]*SubscriptionPayoutStatistics)
		}
		if _, ok := y[items[i].NodeAddr]; !ok {
			y[items[i].NodeAddr] = make(map[time.Time]*SubscriptionPayoutStatistics)
		}

		dayTimestamp := utils.DayDate(items[i].Timestamp)
		if _, ok := d[items[i].NodeAddr][dayTimestamp]; !ok {
			d[items[i].NodeAddr][dayTimestamp] = NewSubscriptionPayoutStatistics("day")
		}

		weekTimestamp := utils.ISOWeekDate(items[i].Timestamp)
		if _, ok := w[items[i].NodeAddr][weekTimestamp]; !ok {
			w[items[i].NodeAddr][weekTimestamp] = NewSubscriptionPayoutStatistics("week")
		}

		monthTimestamp := utils.MonthDate(items[i].Timestamp)
		if _, ok := m[items[i].NodeAddr][monthTimestamp]; !ok {
			m[items[i].NodeAddr][monthTimestamp] = NewSubscriptionPayoutStatistics("month")
		}

		yearTimestamp := utils.YearDate(items[i].Timestamp)
		if _, ok := y[items[i].NodeAddr][yearTimestamp]; !ok {
			y[items[i].NodeAddr][yearTimestamp] = NewSubscriptionPayoutStatistics("year")
		}

		d[items[i].NodeAddr][dayTimestamp].HoursEarning = d[items[i].NodeAddr][dayTimestamp].HoursEarning.Add(items[i].Payment)
		w[items[i].NodeAddr][weekTimestamp].HoursEarning = w[items[i].NodeAddr][weekTimestamp].HoursEarning.Add(items[i].Payment)
		m[items[i].NodeAddr][monthTimestamp].HoursEarning = m[items[i].NodeAddr][monthTimestamp].HoursEarning.Add(items[i].Payment)
		y[items[i].NodeAddr][yearTimestamp].HoursEarning = y[items[i].NodeAddr][yearTimestamp].HoursEarning.Add(items[i].Payment)
	}

	for s := range d {
		for t := range d[s] {
			result = append(result, d[s][t].Result(s, t))
		}
	}

	for s := range w {
		for t := range w[s] {
			result = append(result, w[s][t].Result(s, t))
		}
	}

	for s := range m {
		for t := range m[s] {
			result = append(result, m[s][t].Result(s, t))
		}
	}

	for s := range y {
		for t := range y[s] {
			result = append(result, y[s][t].Result(s, t))
		}
	}

	return result, nil
}
