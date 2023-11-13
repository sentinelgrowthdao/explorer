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
		EarningsForHours types.Coins
	}
)

func NewSubscriptionPayoutStatistics() *SubscriptionPayoutStatistics {
	return &SubscriptionPayoutStatistics{
		EarningsForHours: types.NewCoins(nil),
	}
}

func (s *SubscriptionPayoutStatistics) Result(addr string, timestamp time.Time) bson.M {
	return bson.M{
		"addr":               addr,
		"timestamp":          timestamp,
		"earnings_for_bytes": s.EarningsForHours,
	}
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

	d := make(map[string]map[time.Time]*SubscriptionPayoutStatistics)
	for i := 0; i < len(items); i++ {
		if _, ok := d[items[i].NodeAddr]; !ok {
			d[items[i].NodeAddr] = make(map[time.Time]*SubscriptionPayoutStatistics)
		}

		timestamp := utils.DayDate(items[i].Timestamp)
		if _, ok := d[items[i].NodeAddr][timestamp]; !ok {
			d[items[i].NodeAddr][timestamp] = NewSubscriptionPayoutStatistics()
		}

		d[items[i].NodeAddr][timestamp].EarningsForHours = d[items[i].NodeAddr][timestamp].EarningsForHours.Add(items[i].Payment)
	}

	for s, m := range d {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	return result, nil
}
