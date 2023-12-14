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
		Timeframe      string
		ActiveSession  int64
		BytesEarning   types.Coins
		EndSession     int64
		SessionAddress map[string]bool
		StartSession   int64
	}
)

func NewSessionStatistics(timeframe string) *SessionStatistics {
	return &SessionStatistics{
		Timeframe:      timeframe,
		BytesEarning:   types.NewCoins(nil),
		SessionAddress: make(map[string]bool),
	}
}

func (s *SessionStatistics) Result(addr string, timestamp time.Time) bson.M {
	res := bson.M{
		"addr":      addr,
		"timeframe": s.Timeframe,
		"timestamp": timestamp,
	}

	if s.ActiveSession != 0 {
		res["active_session"] = s.ActiveSession
	}
	if s.BytesEarning.Len() != 0 {
		res["bytes_earning"] = s.BytesEarning
	}
	if s.EndSession != 0 {
		res["end_session"] = s.EndSession
	}
	if len(s.SessionAddress) != 0 {
		res["session_address"] = len(s.SessionAddress)
	}
	if s.StartSession != 0 {
		res["start_session"] = s.StartSession
	}

	return res
}

func StatisticsFromSessions(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time, excludeAddrs []string) (result []bson.M, err error) {
	log.Println("StatisticsFromSessions", minTimestamp, maxTimestamp)

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"acc_addr":        1,
		"end_timestamp":   1,
		"node_addr":       1,
		"payment":         1,
		"start_timestamp": 1,
	}

	items, err := database.SessionFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[string]map[time.Time]*SessionStatistics)
		w = make(map[string]map[time.Time]*SessionStatistics)
		m = make(map[string]map[time.Time]*SessionStatistics)
		y = make(map[string]map[time.Time]*SessionStatistics)
	)

	for i := 0; i < len(items); i++ {
		if _, ok := d[items[i].NodeAddr]; !ok {
			d[items[i].NodeAddr] = make(map[time.Time]*SessionStatistics)
		}
		if _, ok := w[items[i].NodeAddr]; !ok {
			w[items[i].NodeAddr] = make(map[time.Time]*SessionStatistics)
		}
		if _, ok := m[items[i].NodeAddr]; !ok {
			m[items[i].NodeAddr] = make(map[time.Time]*SessionStatistics)
		}
		if _, ok := y[items[i].NodeAddr]; !ok {
			y[items[i].NodeAddr] = make(map[time.Time]*SessionStatistics)
		}

		exclude := utils.ContainsString(excludeAddrs, items[i].AccAddr)
		if exclude {
			continue
		}

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
			if _, ok := d[items[i].NodeAddr][t]; !ok {
				d[items[i].NodeAddr][t] = NewSessionStatistics("day")
			}

			d[items[i].NodeAddr][t].ActiveSession += 1
			d[items[i].NodeAddr][t].SessionAddress[items[i].AccAddr] = true
		}

		weekStartTimestamp, weekEndTimestamp := utils.ISOWeekDate(startTimestamp), utils.ISOWeekDate(endTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			if _, ok := w[items[i].NodeAddr][t]; !ok {
				w[items[i].NodeAddr][t] = NewSessionStatistics("week")
			}

			w[items[i].NodeAddr][t].ActiveSession += 1
			w[items[i].NodeAddr][t].SessionAddress[items[i].AccAddr] = true
		}

		monthStartTimestamp, monthEndTimestamp := utils.MonthDate(startTimestamp), utils.MonthDate(endTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			if _, ok := m[items[i].NodeAddr][t]; !ok {
				m[items[i].NodeAddr][t] = NewSessionStatistics("month")
			}

			m[items[i].NodeAddr][t].ActiveSession += 1
			m[items[i].NodeAddr][t].SessionAddress[items[i].AccAddr] = true
		}

		yearStartTimestamp, yearEndTimestamp := utils.YearDate(startTimestamp), utils.YearDate(endTimestamp)
		for t := yearStartTimestamp; !t.After(yearEndTimestamp); t = t.AddDate(1, 0, 0) {
			if _, ok := y[items[i].NodeAddr][t]; !ok {
				y[items[i].NodeAddr][t] = NewSessionStatistics("year")
			}

			y[items[i].NodeAddr][t].ActiveSession += 1
			y[items[i].NodeAddr][t].SessionAddress[items[i].AccAddr] = true
		}

		if !items[i].EndTimestamp.IsZero() {
			d[items[i].NodeAddr][dayEndTimestamp].EndSession += 1
			w[items[i].NodeAddr][weekEndTimestamp].EndSession += 1
			m[items[i].NodeAddr][monthEndTimestamp].EndSession += 1
			y[items[i].NodeAddr][yearEndTimestamp].EndSession += 1
		}
		if !items[i].StartTimestamp.IsZero() {
			d[items[i].NodeAddr][dayStartTimestamp].StartSession += 1
			w[items[i].NodeAddr][weekStartTimestamp].StartSession += 1
			m[items[i].NodeAddr][monthStartTimestamp].StartSession += 1
			y[items[i].NodeAddr][yearStartTimestamp].StartSession += 1
		}
	}

	for i := 0; i < len(items); i++ {
		endTimestamp := items[i].EndTimestamp
		if items[i].EndTimestamp.IsZero() {
			endTimestamp = maxTimestamp
		}

		dayEndTimestamp := utils.DayDate(endTimestamp)
		if _, ok := d[items[i].NodeAddr][dayEndTimestamp]; !ok {
			d[items[i].NodeAddr][dayEndTimestamp] = NewSessionStatistics("day")
		}

		weekEndTimestamp := utils.ISOWeekDate(endTimestamp)
		if _, ok := w[items[i].NodeAddr][weekEndTimestamp]; !ok {
			w[items[i].NodeAddr][weekEndTimestamp] = NewSessionStatistics("week")
		}

		monthEndTimestamp := utils.MonthDate(endTimestamp)
		if _, ok := m[items[i].NodeAddr][monthEndTimestamp]; !ok {
			m[items[i].NodeAddr][monthEndTimestamp] = NewSessionStatistics("month")
		}

		yearEndTimestamp := utils.YearDate(endTimestamp)
		if _, ok := y[items[i].NodeAddr][yearEndTimestamp]; !ok {
			y[items[i].NodeAddr][yearEndTimestamp] = NewSessionStatistics("year")
		}

		if items[i].Payment != nil {
			d[items[i].NodeAddr][dayEndTimestamp].BytesEarning = d[items[i].NodeAddr][dayEndTimestamp].BytesEarning.Add(items[i].Payment)
			w[items[i].NodeAddr][weekEndTimestamp].BytesEarning = w[items[i].NodeAddr][weekEndTimestamp].BytesEarning.Add(items[i].Payment)
			m[items[i].NodeAddr][monthEndTimestamp].BytesEarning = m[items[i].NodeAddr][monthEndTimestamp].BytesEarning.Add(items[i].Payment)
			y[items[i].NodeAddr][yearEndTimestamp].BytesEarning = y[items[i].NodeAddr][yearEndTimestamp].BytesEarning.Add(items[i].Payment)
		}
	}

	for s, m := range d {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	for s, m := range w {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	for s, m := range m {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	for s, m := range y {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	return result, nil
}
