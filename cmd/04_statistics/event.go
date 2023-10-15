package main

import (
	"context"
	"fmt"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	EventStatistics struct {
		Timeframe        string
		ActiveNode       map[string]int64
		SessionBandwidth map[uint64]*types.Bandwidth
		SessionDuration  map[uint64]int64
	}
)

func NewEventStatistics(timeframe string) *EventStatistics {
	return &EventStatistics{
		Timeframe:        timeframe,
		ActiveNode:       make(map[string]int64),
		SessionBandwidth: make(map[uint64]*types.Bandwidth),
		SessionDuration:  make(map[uint64]int64),
	}
}

func (es *EventStatistics) Result(timestamp time.Time) []bson.M {
	var activeNode int64 = 0
	for _, v := range es.ActiveNode {
		activeNode = activeNode + v
	}

	var sessionBandwidth = &types.Bandwidth{}
	for _, v := range es.SessionBandwidth {
		sessionBandwidth = sessionBandwidth.Add(v)
	}

	var sessionDuration int64 = 0
	for _, v := range es.SessionDuration {
		sessionDuration = sessionDuration + v
	}

	return []bson.M{
		{
			"type":      types.StatisticTypeActiveNode,
			"timeframe": es.Timeframe,
			"timestamp": timestamp,
			"value":     activeNode,
		},
		{
			"type":      types.StatisticTypeSessionBandwidth,
			"timeframe": es.Timeframe,
			"timestamp": timestamp,
			"value":     sessionBandwidth,
		},
		{
			"type":      types.StatisticTypeSessionDuration,
			"timeframe": es.Timeframe,
			"timestamp": timestamp,
			"value":     sessionDuration,
		},
	}
}

func StatisticsFromEvents(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromEvents")

	filter := bson.M{
		"type": bson.M{
			"$in": []string{
				types.EventTypeNodeUpdateStatus,
				types.EventTypeSessionUpdateDetails,
			},
		},
		"$or": []bson.M{
			{
				"status": bson.M{"$exists": false},
			},
			{
				"status": hubtypes.StatusActive.String(),
			},
		},
	}
	projection := bson.M{
		"_id":          0,
		"bandwidth":    1,
		"duration":     1,
		"node_address": 1,
		"session_id":   1,
		"timestamp":    1,
		"type":         1,
	}
	sort := bson.D{
		bson.E{Key: "timestamp", Value: 1},
	}
	opts := options.Find().
		SetProjection(projection).
		SetSort(sort)

	items, err := database.EventFind(ctx, db, filter, opts)
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*EventStatistics)
		w = make(map[time.Time]*EventStatistics)
		m = make(map[time.Time]*EventStatistics)
		y = make(map[time.Time]*EventStatistics)
	)

	for i := 0; i < len(items); i++ {
		dayTimestamp := utils.DayDate(items[i].Timestamp)
		if _, ok := d[dayTimestamp]; !ok {
			d[dayTimestamp] = NewEventStatistics("day")
		}

		weekTimestamp := utils.ISOWeekDate(items[i].Timestamp)
		if _, ok := w[weekTimestamp]; !ok {
			w[weekTimestamp] = NewEventStatistics("week")
		}

		monthTimestamp := utils.MonthDate(items[i].Timestamp)
		if _, ok := m[monthTimestamp]; !ok {
			m[monthTimestamp] = NewEventStatistics("month")
		}

		yearTimestamp := utils.YearDate(items[i].Timestamp)
		if _, ok := y[yearTimestamp]; !ok {
			y[yearTimestamp] = NewEventStatistics("year")
		}

		switch items[i].Type {
		case types.EventTypeNodeUpdateStatus:
			d[dayTimestamp].ActiveNode[items[i].NodeAddress] = 1
			w[weekTimestamp].ActiveNode[items[i].NodeAddress] = 1
			m[monthTimestamp].ActiveNode[items[i].NodeAddress] = 1
			y[yearTimestamp].ActiveNode[items[i].NodeAddress] = 1
		case types.EventTypeSessionUpdateDetails:
			d[dayTimestamp].SessionBandwidth[items[i].SessionID] = items[i].Bandwidth
			d[dayTimestamp].SessionDuration[items[i].SessionID] = items[i].Duration
			if v, ok := d[dayTimestamp.AddDate(0, 0, -1)]; ok {
				if v, ok := v.SessionBandwidth[items[i].SessionID]; ok {
					d[dayTimestamp].SessionBandwidth[items[i].SessionID] = d[dayTimestamp].SessionBandwidth[items[i].SessionID].Sub(v)
				}
				if v, ok := v.SessionDuration[items[i].SessionID]; ok {
					d[dayTimestamp].SessionDuration[items[i].SessionID] = d[dayTimestamp].SessionDuration[items[i].SessionID] - v
				}
			}

			w[weekTimestamp].SessionBandwidth[items[i].SessionID] = items[i].Bandwidth
			w[weekTimestamp].SessionDuration[items[i].SessionID] = items[i].Duration
			if v, ok := w[weekTimestamp.AddDate(0, 0, -7)]; ok {
				if v, ok := v.SessionBandwidth[items[i].SessionID]; ok {
					w[weekTimestamp].SessionBandwidth[items[i].SessionID] = w[weekTimestamp].SessionBandwidth[items[i].SessionID].Sub(v)
				}
				if v, ok := v.SessionDuration[items[i].SessionID]; ok {
					w[weekTimestamp].SessionDuration[items[i].SessionID] = w[weekTimestamp].SessionDuration[items[i].SessionID] - v
				}
			}

			m[monthTimestamp].SessionBandwidth[items[i].SessionID] = items[i].Bandwidth
			m[monthTimestamp].SessionDuration[items[i].SessionID] = items[i].Duration
			if v, ok := m[monthTimestamp.AddDate(0, -1, 0)]; ok {
				if v, ok := v.SessionBandwidth[items[i].SessionID]; ok {
					m[monthTimestamp].SessionBandwidth[items[i].SessionID] = m[monthTimestamp].SessionBandwidth[items[i].SessionID].Sub(v)
				}
				if v, ok := v.SessionDuration[items[i].SessionID]; ok {
					m[monthTimestamp].SessionDuration[items[i].SessionID] = m[monthTimestamp].SessionDuration[items[i].SessionID] - v
				}
			}

			y[yearTimestamp].SessionBandwidth[items[i].SessionID] = items[i].Bandwidth
			y[yearTimestamp].SessionDuration[items[i].SessionID] = items[i].Duration
			if v, ok := y[yearTimestamp.AddDate(-1, 0, 0)]; ok {
				if v, ok := v.SessionBandwidth[items[i].SessionID]; ok {
					y[yearTimestamp].SessionBandwidth[items[i].SessionID] = y[yearTimestamp].SessionBandwidth[items[i].SessionID].Sub(v)
				}
				if v, ok := v.SessionDuration[items[i].SessionID]; ok {
					y[yearTimestamp].SessionDuration[items[i].SessionID] = y[yearTimestamp].SessionDuration[items[i].SessionID] - v
				}
			}
		default:
			return nil, fmt.Errorf("invalid type %s", items[i].Type)
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
