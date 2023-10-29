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
		EarningsForBytes  types.Coins
		SessionEndCount   int64
		SessionStartCount int64
	}
)

func NewSessionStatistics() *SessionStatistics {
	return &SessionStatistics{
		EarningsForBytes: types.NewCoins(nil),
	}
}

func (s *SessionStatistics) Result(addr string, timestamp time.Time) bson.M {
	return bson.M{
		"addr":                addr,
		"timestamp":           timestamp,
		"earnings_for_bytes":  s.EarningsForBytes,
		"session_end_count":   s.SessionEndCount,
		"session_start_count": s.SessionStartCount,
	}
}

func StatisticsFromSessions(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time) (result []bson.M, err error) {
	log.Println("StatisticsFromSessions", minTimestamp, maxTimestamp)

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"end_timestamp":   1,
		"node_addr":       1,
		"payment":         1,
		"start_timestamp": 1,
	}

	items, err := database.SessionFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	d := make(map[string]map[time.Time]*SessionStatistics)
	for i := 0; i < len(items); i++ {
		startTimestamp := items[i].StartTimestamp
		if items[i].StartTimestamp.IsZero() {
			startTimestamp = minTimestamp
		}

		endTimestamp := items[i].EndTimestamp
		if items[i].EndTimestamp.IsZero() {
			endTimestamp = maxTimestamp
		}

		if _, ok := d[items[i].NodeAddr]; !ok {
			d[items[i].NodeAddr] = make(map[time.Time]*SessionStatistics)
		}

		startTimestamp = utils.DayDate(startTimestamp)
		if _, ok := d[items[i].NodeAddr][startTimestamp]; !ok {
			d[items[i].NodeAddr][startTimestamp] = NewSessionStatistics()
		}

		endTimestamp = utils.DayDate(endTimestamp)
		if _, ok := d[items[i].NodeAddr][endTimestamp]; !ok {
			d[items[i].NodeAddr][endTimestamp] = NewSessionStatistics()
		}

		if !items[i].EndTimestamp.IsZero() {
			d[items[i].NodeAddr][endTimestamp].SessionEndCount += 1
		}
		if items[i].Payment != nil {
			d[items[i].NodeAddr][endTimestamp].EarningsForBytes = d[items[i].NodeAddr][endTimestamp].EarningsForBytes.Add(items[i].Payment)
		}
		if !items[i].StartTimestamp.IsZero() {
			d[items[i].NodeAddr][startTimestamp].SessionStartCount += 1
		}
	}

	for s, m := range d {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	return result, nil
}
