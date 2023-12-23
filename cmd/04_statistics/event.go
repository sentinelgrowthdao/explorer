package main

import (
	"context"
	"log"
	"sort"
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

func StatisticsFromSessionEvents(ctx context.Context, db *mongo.Database, excludeAddrs []string) (result []bson.M, err error) {
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
				"session_id": 1,
				"timestamp":  1,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"session_id": "$session_id",
					"timestamp": bson.M{
						"$dateFromParts": bson.M{
							"day": bson.M{
								"$dayOfMonth": "$timestamp",
							},
							"month": bson.M{
								"$month": "$timestamp",
							},
							"year": bson.M{
								"$year": "$timestamp",
							},
						},
					},
				},
				"bandwidth": bson.M{
					"$first": "$bandwidth",
				},
				"duration": bson.M{
					"$first": "$duration",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         database.SessionCollectionName,
				"localField":   "_id.session_id",
				"foreignField": "id",
				"as":           "session",
			},
		},
		{
			"$addFields": bson.M{
				"acc_addr": "$session.acc_addr",
			},
		},
		{
			"$unwind": "$acc_addr",
		},
		{
			"$match": bson.M{
				"acc_addr": bson.M{
					"$nin": excludeAddrs,
				},
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
		{
			"$sort": bson.M{
				"timestamp": 1,
			},
		},
	}

	cursor, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var (
		d = make(map[time.Time]*SessionEventStatistics)
		w = make(map[time.Time]*SessionEventStatistics)
		m = make(map[time.Time]*SessionEventStatistics)
		y = make(map[time.Time]*SessionEventStatistics)
	)

	for cursor.Next(ctx) {
		var item bson.M
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}

		bandwidth := types.BandwidthFromInterface(item["bandwidth"])
		duration := types.Int64FromInterface(item["duration"])
		sessionID := types.Uint64FromInterface(item["session_id"])
		timestamp := types.TimeFromInterface(item["timestamp"])

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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	var dKeys []time.Time
	for t := range d {
		dKeys = append(dKeys, t)
	}

	sort.Slice(dKeys, func(i, j int) bool {
		return dKeys[i].After(dKeys[j])
	})

	for _, t := range dKeys {
		for u := range d[t].SessionBandwidth {
			if v, ok := d[t.AddDate(0, 0, -1)]; ok {
				if v, ok := v.SessionBandwidth[u]; ok {
					d[t].SessionBandwidth[u] = d[t].SessionBandwidth[u].Sub(v)
				}
				if v, ok := v.SessionDuration[u]; ok {
					d[t].SessionDuration[u] = d[t].SessionDuration[u] - v
				}
			}
		}
	}

	var wKeys []time.Time
	for t := range w {
		wKeys = append(wKeys, t)
	}

	sort.Slice(wKeys, func(i, j int) bool {
		return wKeys[i].After(wKeys[j])
	})

	for _, t := range wKeys {
		for u := range w[t].SessionBandwidth {
			if v, ok := w[t.AddDate(0, 0, -7)]; ok {
				if v, ok := v.SessionBandwidth[u]; ok {
					w[t].SessionBandwidth[u] = w[t].SessionBandwidth[u].Sub(v)
				}
				if v, ok := v.SessionDuration[u]; ok {
					w[t].SessionDuration[u] = w[t].SessionDuration[u] - v
				}
			}
		}
	}

	var mKeys []time.Time
	for t := range m {
		mKeys = append(mKeys, t)
	}

	sort.Slice(mKeys, func(i, j int) bool {
		return mKeys[i].After(mKeys[j])
	})

	for _, t := range mKeys {
		for u := range m[t].SessionBandwidth {
			if v, ok := m[t.AddDate(0, -1, 0)]; ok {
				if v, ok := v.SessionBandwidth[u]; ok {
					m[t].SessionBandwidth[u] = m[t].SessionBandwidth[u].Sub(v)
				}
				if v, ok := v.SessionDuration[u]; ok {
					m[t].SessionDuration[u] = m[t].SessionDuration[u] - v
				}
			}
		}
	}

	var yKeys []time.Time
	for t := range y {
		yKeys = append(yKeys, t)
	}

	sort.Slice(yKeys, func(i, j int) bool {
		return yKeys[i].After(yKeys[j])
	})

	for _, t := range yKeys {
		for u := range y[t].SessionBandwidth {
			if v, ok := y[t.AddDate(-1, 0, 0)]; ok {
				if v, ok := v.SessionBandwidth[u]; ok {
					y[t].SessionBandwidth[u] = y[t].SessionBandwidth[u].Sub(v)
				}
				if v, ok := v.SessionDuration[u]; ok {
					y[t].SessionDuration[u] = y[t].SessionDuration[u] - v
				}
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
			"$group": bson.M{
				"_id": bson.M{
					"node_addr": "$node_addr",
					"timestamp": bson.M{
						"$dateFromParts": bson.M{
							"day": bson.M{
								"$dayOfMonth": "$timestamp",
							},
							"month": bson.M{
								"$month": "$timestamp",
							},
							"year": bson.M{
								"$year": "$timestamp",
							},
						},
					},
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

	cursor, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var (
		d = make(map[time.Time]*NodeEventStatistics)
		w = make(map[time.Time]*NodeEventStatistics)
		m = make(map[time.Time]*NodeEventStatistics)
		y = make(map[time.Time]*NodeEventStatistics)
	)

	for cursor.Next(ctx) {
		var item bson.M
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}

		nodeAddr := types.StringFromInterface(item["node_addr"])
		timestamp := types.TimeFromInterface(item["timestamp"])

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

	if err := cursor.Err(); err != nil {
		return nil, err
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
