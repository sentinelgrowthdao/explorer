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
	"github.com/sentinel-official/explorer/utils"
)

type (
	SubscriptionStatistics struct {
		SubscriptionBytes      string
		SubscriptionEndCount   int64
		SubscriptionHours      int64
		SubscriptionStartCount int64
	}
)

func NewSubscriptionStatistics() *SubscriptionStatistics {
	return &SubscriptionStatistics{}
}

func (s *SubscriptionStatistics) Result(addr string, timestamp time.Time) bson.M {
	return bson.M{
		"addr":                     addr,
		"timestamp":                timestamp,
		"subscription_bytes":       s.SubscriptionBytes,
		"subscription_end_count":   s.SubscriptionEndCount,
		"subscription_hours":       s.SubscriptionHours,
		"subscription_start_count": s.SubscriptionStartCount,
	}
}

func StatisticsFromSubscriptions(ctx context.Context, db *mongo.Database, minTimestamp, maxTimestamp time.Time) (result []bson.M, err error) {
	log.Println("StatisticsFromSubscriptions", minTimestamp, maxTimestamp)

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"end_timestamp":   1,
		"gigabytes":       1,
		"hours":           1,
		"node_addr":       1,
		"start_timestamp": 1,
	}

	items, err := database.SubscriptionFind(ctx, db, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	d := make(map[string]map[time.Time]*SubscriptionStatistics)
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
			d[items[i].NodeAddr] = make(map[time.Time]*SubscriptionStatistics)
		}

		startTimestamp = utils.DayDate(startTimestamp)
		if _, ok := d[items[i].NodeAddr][startTimestamp]; !ok {
			d[items[i].NodeAddr][startTimestamp] = NewSubscriptionStatistics()
		}

		endTimestamp = utils.DayDate(endTimestamp)
		if _, ok := d[items[i].NodeAddr][endTimestamp]; !ok {
			d[items[i].NodeAddr][endTimestamp] = NewSubscriptionStatistics()
		}

		if items[i].Gigabytes != 0 {
			bytes := hubtypes.Gigabyte.MulRaw(items[i].Gigabytes)
			d[items[i].NodeAddr][startTimestamp].SubscriptionBytes = utils.MustIntFromString(d[items[i].NodeAddr][startTimestamp].SubscriptionBytes).Add(bytes).String()
		}
		if !items[i].EndTimestamp.IsZero() {
			d[items[i].NodeAddr][endTimestamp].SubscriptionEndCount += 1
		}
		if items[i].Hours != 0 {
			d[items[i].NodeAddr][startTimestamp].SubscriptionHours += items[i].Hours
		}
		if !items[i].StartTimestamp.IsZero() {
			d[items[i].NodeAddr][startTimestamp].SubscriptionStartCount += 1
		}
	}

	for s, m := range d {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	return result, nil
}
