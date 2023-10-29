package main

import (
	"context"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	SessionEventStatistics struct {
		Timeframe        string
		SessionBandwidth map[uint64]*types.Bandwidth
		SessionDuration  map[uint64]int64
	}

	NodeEventStatistics struct {
		Timeframe  string
		ActiveNode map[string]bool
	}
)

func NewSessionEventStatistics(timeframe string) *SessionEventStatistics {
	return &SessionEventStatistics{
		Timeframe:        timeframe,
		SessionBandwidth: make(map[uint64]*types.Bandwidth),
		SessionDuration:  make(map[uint64]int64),
	}
}

func (s *SessionEventStatistics) Result(timestamp time.Time) []bson.M {
	var sessionBandwidth = &types.Bandwidth{}
	for _, v := range s.SessionBandwidth {
		sessionBandwidth = sessionBandwidth.Add(v)
	}

	var sessionDuration int64 = 0
	for _, v := range s.SessionDuration {
		sessionDuration = sessionDuration + v
	}

	return []bson.M{
		{
			"type":      types.StatisticTypeSessionBytes,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     sessionBandwidth,
		},
		{
			"type":      types.StatisticTypeSessionDuration,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     sessionDuration,
		},
	}
}

func NewNodeEventStatistics(timeframe string) *NodeEventStatistics {
	return &NodeEventStatistics{
		Timeframe:  timeframe,
		ActiveNode: make(map[string]bool),
	}
}

func (s *NodeEventStatistics) Result(timestamp time.Time) []bson.M {
	return []bson.M{
		{
			"type":      types.StatisticTypeActiveNode,
			"timeframe": s.Timeframe,
			"timestamp": timestamp,
			"value":     len(s.ActiveNode),
		},
	}
}

func StatisticsFromSessionEvents(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromSessionEvents")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type": types.EventTypeSessionUpdateDetails,
			},
		},
		{
			"$sort": bson.M{
				"timestamp": -1,
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"bandwidth":  1,
				"duration":   1,
				"timestamp":  1,
				"session_id": 1,
			},
		},
		{
			"$addFields": bson.M{
				"timestamp": bson.M{
					"$dateFromParts": bson.M{
						"day":   bson.M{"$dayOfMonth": "$timestamp"},
						"month": bson.M{"$month": "$timestamp"},
						"year":  bson.M{"$year": "$timestamp"},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"session_id": "$session_id",
					"timestamp":  "$timestamp",
				},
				"bandwidth": bson.M{"$first": "$bandwidth"},
				"duration":  bson.M{"$first": "$duration"},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"bandwidth":  "$bandwidth",
				"duration":   "$duration",
				"session_id": "$_id.session_id",
				"timestamp":  "$_id.timestamp",
			},
		},
	}

	items, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*SessionEventStatistics)
		w = make(map[time.Time]*SessionEventStatistics)
		m = make(map[time.Time]*SessionEventStatistics)
		y = make(map[time.Time]*SessionEventStatistics)
	)

	for i := 0; i < len(items); i++ {
		bandwidth := types.BandwidthFromInterface(items[i]["bandwidth"])
		duration := types.Int64FromInterface(items[i]["duration"])
		sessionID := types.Uint64FromInterface(items[i]["session_id"])
		timestamp := types.TimeFromInterface(items[i]["timestamp"])

		dayTimestamp := utils.DayDate(timestamp)
		if _, ok := d[dayTimestamp]; !ok {
			d[dayTimestamp] = NewSessionEventStatistics("day")
		}

		weekTimestamp := utils.ISOWeekDate(timestamp)
		if _, ok := w[weekTimestamp]; !ok {
			w[weekTimestamp] = NewSessionEventStatistics("week")
		}

		monthTimestamp := utils.MonthDate(timestamp)
		if _, ok := m[monthTimestamp]; !ok {
			m[monthTimestamp] = NewSessionEventStatistics("month")
		}

		yearTimestamp := utils.YearDate(timestamp)
		if _, ok := y[yearTimestamp]; !ok {
			y[yearTimestamp] = NewSessionEventStatistics("year")
		}

		d[dayTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		d[dayTimestamp].SessionDuration[sessionID] = duration

		w[weekTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		w[weekTimestamp].SessionDuration[sessionID] = duration

		m[monthTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		m[monthTimestamp].SessionDuration[sessionID] = duration

		y[yearTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		y[yearTimestamp].SessionDuration[sessionID] = duration
	}

	for i := 0; i < len(items); i++ {
		sessionID := types.Uint64FromInterface(items[i]["session_id"])
		timestamp := types.TimeFromInterface(items[i]["timestamp"])

		dayTimestamp := utils.DayDate(timestamp)
		if v, ok := d[dayTimestamp.AddDate(0, 0, -1)]; ok {
			if v, ok := v.SessionBandwidth[sessionID]; ok {
				d[dayTimestamp].SessionBandwidth[sessionID] = d[dayTimestamp].SessionBandwidth[sessionID].Sub(v)
			}
			if v, ok := v.SessionDuration[sessionID]; ok {
				d[dayTimestamp].SessionDuration[sessionID] = d[dayTimestamp].SessionDuration[sessionID] - v
			}
		}

		weekTimestamp := utils.ISOWeekDate(timestamp)
		if v, ok := d[weekTimestamp.AddDate(0, 0, -7)]; ok {
			if v, ok := v.SessionBandwidth[sessionID]; ok {
				d[weekTimestamp].SessionBandwidth[sessionID] = d[weekTimestamp].SessionBandwidth[sessionID].Sub(v)
			}
			if v, ok := v.SessionDuration[sessionID]; ok {
				d[weekTimestamp].SessionDuration[sessionID] = d[weekTimestamp].SessionDuration[sessionID] - v
			}
		}

		monthTimestamp := utils.MonthDate(timestamp)
		if v, ok := d[monthTimestamp.AddDate(0, -1, 0)]; ok {
			if v, ok := v.SessionBandwidth[sessionID]; ok {
				d[monthTimestamp].SessionBandwidth[sessionID] = d[monthTimestamp].SessionBandwidth[sessionID].Sub(v)
			}
			if v, ok := v.SessionDuration[sessionID]; ok {
				d[monthTimestamp].SessionDuration[sessionID] = d[monthTimestamp].SessionDuration[sessionID] - v
			}
		}

		yearTimestamp := utils.YearDate(timestamp)
		if v, ok := d[yearTimestamp.AddDate(-1, 0, 0)]; ok {
			if v, ok := v.SessionBandwidth[sessionID]; ok {
				d[yearTimestamp].SessionBandwidth[sessionID] = d[yearTimestamp].SessionBandwidth[sessionID].Sub(v)
			}
			if v, ok := v.SessionDuration[sessionID]; ok {
				d[yearTimestamp].SessionDuration[sessionID] = d[yearTimestamp].SessionDuration[sessionID] - v
			}
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

func StatisticsFromNodeEvents(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromNodeEvents")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":   types.EventTypeNodeUpdateStatus,
				"status": hubtypes.StatusActive.String(),
			},
		},
		{
			"$addFields": bson.M{
				"timestamp": bson.M{
					"$dateFromParts": bson.M{
						"day":   bson.M{"$dayOfMonth": "$timestamp"},
						"month": bson.M{"$month": "$timestamp"},
						"year":  bson.M{"$year": "$timestamp"},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"node_addr": "$node_addr",
					"timestamp": "$timestamp",
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"node_addr": "$_id.node_addr",
				"timestamp": "$_id.timestamp",
			},
		},
	}

	items, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	var (
		d = make(map[time.Time]*NodeEventStatistics)
		w = make(map[time.Time]*NodeEventStatistics)
		m = make(map[time.Time]*NodeEventStatistics)
		y = make(map[time.Time]*NodeEventStatistics)
	)

	for i := 0; i < len(items); i++ {
		nodeAddr := types.StringFromInterface(items[i]["node_addr"])
		timestamp := types.TimeFromInterface(items[i]["timestamp"])

		dayTimestamp := utils.DayDate(timestamp)
		if _, ok := d[dayTimestamp]; !ok {
			d[dayTimestamp] = NewNodeEventStatistics("day")
		}

		weekTimestamp := utils.ISOWeekDate(timestamp)
		if _, ok := w[weekTimestamp]; !ok {
			w[weekTimestamp] = NewNodeEventStatistics("week")
		}

		monthTimestamp := utils.MonthDate(timestamp)
		if _, ok := m[monthTimestamp]; !ok {
			m[monthTimestamp] = NewNodeEventStatistics("month")
		}

		yearTimestamp := utils.YearDate(timestamp)
		if _, ok := y[yearTimestamp]; !ok {
			y[yearTimestamp] = NewNodeEventStatistics("year")
		}

		d[dayTimestamp].ActiveNode[nodeAddr] = true
		w[weekTimestamp].ActiveNode[nodeAddr] = true
		m[monthTimestamp].ActiveNode[nodeAddr] = true
		y[yearTimestamp].ActiveNode[nodeAddr] = true
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
