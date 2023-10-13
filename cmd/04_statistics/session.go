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
	SessionStatistics struct {
		Timeframe      string
		ActiveSession  int64
		EndSession     int64
		SessionPayment types.Coins
		StartSession   int64
	}
)

func NewSessionStatistics(timeframe string) *SessionStatistics {
	return &SessionStatistics{
		Timeframe: timeframe,
	}
}

func (ss *SessionStatistics) Result(timestamp time.Time) []bson.M {
	return []bson.M{
		{
			"type":      types.StatisticTypeActiveSession,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.ActiveSession,
		},
		{
			"type":      types.StatisticTypeEndSession,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.EndSession,
		},
		{
			"type":      types.StatisticTypeSessionPayment,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.SessionPayment,
		},
		{
			"type":      types.StatisticTypeStartSession,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.StartSession,
		},
	}
}

func StatisticsFromSessions(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	filter := bson.M{
		"start_timestamp": bson.M{
			"$exists": true,
		},
	}
	projection := bson.M{}

	items, err := database.SessionFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	var (
		d map[time.Time]*SessionStatistics
		w map[time.Time]*SessionStatistics
		m map[time.Time]*SessionStatistics
		y map[time.Time]*SessionStatistics
	)

	for i := 0; i < len(items); i++ {
		dayStartTimestamp, dayEndTimestamp := utils.DayDate(items[i].StartTimestamp), utils.DayDate(items[i].EndTimestamp)
		for t := dayStartTimestamp; !t.After(dayEndTimestamp); t = t.AddDate(0, 0, 1) {
			if _, ok := d[t]; !ok {
				d[t] = NewSessionStatistics("day")
			}

			d[t].ActiveSession += 1
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(items[i].StartTimestamp), utils.ISOWeekDate(items[i].EndTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			if _, ok := w[t]; !ok {
				w[t] = NewSessionStatistics("week")
			}

			w[t].ActiveSession += 1
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(items[i].StartTimestamp), utils.MonthDate(items[i].EndTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			if _, ok := m[t]; !ok {
				m[t] = NewSessionStatistics("month")
			}

			m[t].ActiveSession += 1
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(items[i].StartTimestamp), utils.YearDate(items[i].EndTimestamp)
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
			d[dayEndTimestamp].SessionPayment = d[dayStartTimestamp].SessionPayment.Add(items[i].Payment)
			w[weekEndTimestamp].SessionPayment = w[weekStartTimestamp].SessionPayment.Add(items[i].Payment)
			m[monthEndTimestamp].SessionPayment = m[monthStartTimestamp].SessionPayment.Add(items[i].Payment)
			y[yearEndTimestamp].SessionPayment = y[yearStartTimestamp].SessionPayment.Add(items[i].Payment)
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
