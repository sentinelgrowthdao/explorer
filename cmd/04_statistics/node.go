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
	NodeStatistics struct {
		Timeframe    string
		RegisterNode int64
	}
)

func NewNodeStatistics(timeframe string) *NodeStatistics {
	return &NodeStatistics{
		Timeframe: timeframe,
	}
}

func (ss *NodeStatistics) Result(timestamp time.Time) []bson.M {
	return []bson.M{
		{
			"type":      types.StatisticTypeRegisterNode,
			"timeframe": ss.Timeframe,
			"timestamp": timestamp,
			"value":     ss.RegisterNode,
		},
	}
}

func StatisticsFromNodes(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromNodes")

	filter := bson.M{}
	projection := bson.M{
		"_id":                0,
		"register_timestamp": 1,
	}
	sort := bson.D{
		bson.E{Key: "register_timestamp", Value: 1},
	}

	items, err := database.NodeFind(ctx, db, filter, options.Find().SetProjection(projection).SetSort(sort))
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*NodeStatistics)
		w = make(map[time.Time]*NodeStatistics)
		m = make(map[time.Time]*NodeStatistics)
		y = make(map[time.Time]*NodeStatistics)
	)

	for i := 0; i < len(items); i++ {
		dayRegisterTimestamp := utils.DayDate(items[i].RegisterTimestamp)
		if _, ok := d[dayRegisterTimestamp]; !ok {
			d[dayRegisterTimestamp] = NewNodeStatistics("day")
		}

		weekRegisterTimestamp := utils.ISOWeekDate(items[i].RegisterTimestamp)
		if _, ok := w[weekRegisterTimestamp]; !ok {
			w[weekRegisterTimestamp] = NewNodeStatistics("week")
		}

		monthRegisterTimestamp := utils.MonthDate(items[i].RegisterTimestamp)
		if _, ok := m[monthRegisterTimestamp]; !ok {
			m[monthRegisterTimestamp] = NewNodeStatistics("month")
		}

		yearRegisterTimestamp := utils.YearDate(items[i].RegisterTimestamp)
		if _, ok := y[yearRegisterTimestamp]; !ok {
			y[yearRegisterTimestamp] = NewNodeStatistics("year")
		}

		d[dayRegisterTimestamp].RegisterNode += 1
		w[weekRegisterTimestamp].RegisterNode += 1
		m[monthRegisterTimestamp].RegisterNode += 1
		y[yearRegisterTimestamp].RegisterNode += 1
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
