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
	SessionStatistics struct {
		Timeframe          string
		ActiveSession      int64
		EndSession         int64
		BytesPayment       types.Coins
		BytesStakingReward types.Coins
		StartSession       int64
	}
)

func NewSessionStatistics(timeframe string) *SessionStatistics {
	return &SessionStatistics{
		Timeframe:          timeframe,
		BytesPayment:       types.NewCoins(nil),
		BytesStakingReward: types.NewCoins(nil),
	}
}

func (ss *SessionStatistics) Result(timestamp time.Time) bson.A {
	return bson.A{
		bson.M{
			"type":      types.StatisticTypeActiveSession,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.ActiveSession,
		},
		bson.M{
			"type":      types.StatisticTypeEndSession,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.EndSession,
		},
		bson.M{
			"type":      types.StatisticTypeBytesPayment,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.BytesPayment,
		},
		bson.M{
			"type":      types.StatisticTypeBytesStakingReward,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.BytesStakingReward,
		},
		bson.M{
			"type":      types.StatisticTypeStartSession,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.StartSession,
		},
	}
}

func StatisticsFromSessions(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time) (result bson.A, err error) {
	log.Println("StatisticsFromSessions", minTimestamp, maxTimestamp)

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"end_timestamp":   1,
		"payment":         1,
		"staking_reward":  1,
		"start_timestamp": 1,
	}
	sort := bson.D{
		bson.E{Key: "start_timestamp", Value: 1},
	}

	items, err := database.SessionFind(ctx, db, filter, options.Find().SetProjection(projection).SetSort(sort))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*SessionStatistics)
		w = make(map[time.Time]*SessionStatistics)
		m = make(map[time.Time]*SessionStatistics)
		y = make(map[time.Time]*SessionStatistics)
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
				d[t] = NewSessionStatistics("day")
			}

			d[t].ActiveSession += 1
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(startTimestamp), utils.ISOWeekDate(endTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			if _, ok := w[t]; !ok {
				w[t] = NewSessionStatistics("week")
			}

			w[t].ActiveSession += 1
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(startTimestamp), utils.MonthDate(endTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			if _, ok := m[t]; !ok {
				m[t] = NewSessionStatistics("month")
			}

			m[t].ActiveSession += 1
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(startTimestamp), utils.YearDate(endTimestamp)
		for t := yearStartTimestamp; !t.After(yearEndTimestamp); t = t.AddDate(1, 0, 0) {
			if _, ok := y[t]; !ok {
				y[t] = NewSessionStatistics("year")
			}

			y[t].ActiveSession += 1
		}

		if !items[i].EndTimestamp.IsZero() {
			d[dayEndTimestamp].EndSession += 1
			w[weekEndTimestamp].EndSession += 1
			m[monthEndTimestamp].EndSession += 1
			y[yearEndTimestamp].EndSession += 1
		}
		if items[i].Payment != nil {
			d[dayEndTimestamp].BytesPayment = d[dayEndTimestamp].BytesPayment.Add(items[i].Payment)
			w[weekEndTimestamp].BytesPayment = w[weekEndTimestamp].BytesPayment.Add(items[i].Payment)
			m[monthEndTimestamp].BytesPayment = m[monthEndTimestamp].BytesPayment.Add(items[i].Payment)
			y[yearEndTimestamp].BytesPayment = y[yearEndTimestamp].BytesPayment.Add(items[i].Payment)
		}
		if items[i].StakingReward != nil {
			d[dayEndTimestamp].BytesStakingReward = d[dayEndTimestamp].BytesStakingReward.Add(items[i].StakingReward)
			w[weekEndTimestamp].BytesStakingReward = w[weekEndTimestamp].BytesStakingReward.Add(items[i].StakingReward)
			m[monthEndTimestamp].BytesStakingReward = m[monthEndTimestamp].BytesStakingReward.Add(items[i].StakingReward)
			y[yearEndTimestamp].BytesStakingReward = y[yearEndTimestamp].BytesStakingReward.Add(items[i].StakingReward)
		}
		if !items[i].StartTimestamp.IsZero() {
			d[dayStartTimestamp].StartSession += 1
			w[weekStartTimestamp].StartSession += 1
			m[monthStartTimestamp].StartSession += 1
			y[yearStartTimestamp].StartSession += 1
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
