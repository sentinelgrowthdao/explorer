package main

import (
	"context"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	EventStatistics struct {
		Timeframe        string
		SessionBandwidth map[uint64]*types.Bandwidth
		SessionDuration  map[uint64]int64
	}
)

func NewStatistics(timeframe string) *EventStatistics {
	return &EventStatistics{
		Timeframe:        timeframe,
		SessionBandwidth: make(map[uint64]*types.Bandwidth),
		SessionDuration:  make(map[uint64]int64),
	}
}

func (s *EventStatistics) Result(addr string, timestamp time.Time) bson.M {
	var sessionBandwidth = &types.Bandwidth{}
	for _, v := range s.SessionBandwidth {
		sessionBandwidth = sessionBandwidth.Add(v)
	}

	var sessionDuration int64 = 0
	for _, v := range s.SessionDuration {
		sessionDuration = sessionDuration + v
	}

	res := bson.M{
		"addr":      addr,
		"timeframe": s.Timeframe,
		"timestamp": timestamp,
	}

	if !sessionBandwidth.IsZero() {
		res["session_bandwidth"] = sessionBandwidth
	}
	if sessionDuration != 0 {
		res["session_duration"] = sessionDuration
	}

	return res
}

func StatisticsFromEvents(ctx context.Context, db *mongo.Database) (result []bson.M, err error) {
	log.Println("StatisticsFromEvents")

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
			"$lookup": bson.M{
				"from":         database.SessionCollectionName,
				"localField":   "_id.session_id",
				"foreignField": "id",
				"as":           "session",
			},
		},
		{
			"$addFields": bson.M{
				"node_addr": "$session.node_addr",
			},
		},
		{
			"$unwind": "$node_addr",
		},
		{
			"$project": bson.M{
				"_id":        0,
				"bandwidth":  "$bandwidth",
				"duration":   "$duration",
				"node_addr":  "$node_addr",
				"session_id": "$_id.session_id",
				"timestamp":  "$_id.timestamp",
			},
		},
	}

	cursor, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var (
		d = make(map[string]map[time.Time]*EventStatistics)
		w = make(map[string]map[time.Time]*EventStatistics)
		m = make(map[string]map[time.Time]*EventStatistics)
		y = make(map[string]map[time.Time]*EventStatistics)
	)

	for cursor.Next(ctx) {
		var item bson.M
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}

		bandwidth := types.BandwidthFromInterface(item["bandwidth"])
		duration := types.Int64FromInterface(item["duration"])
		nodeAddr := types.StringFromInterface(item["node_addr"])
		sessionID := types.Uint64FromInterface(item["session_id"])
		timestamp := types.TimeFromInterface(item["timestamp"])

		if _, ok := d[nodeAddr]; !ok {
			d[nodeAddr] = make(map[time.Time]*EventStatistics)
		}
		if _, ok := w[nodeAddr]; !ok {
			w[nodeAddr] = make(map[time.Time]*EventStatistics)
		}
		if _, ok := m[nodeAddr]; !ok {
			m[nodeAddr] = make(map[time.Time]*EventStatistics)
		}
		if _, ok := y[nodeAddr]; !ok {
			y[nodeAddr] = make(map[time.Time]*EventStatistics)
		}

		dayTimestamp := utils.DayDate(timestamp)
		if _, ok := d[nodeAddr][dayTimestamp]; !ok {
			d[nodeAddr][dayTimestamp] = NewStatistics("day")
		}
		weekTimestamp := utils.ISOWeekDate(timestamp)
		if _, ok := w[nodeAddr][weekTimestamp]; !ok {
			w[nodeAddr][weekTimestamp] = NewStatistics("week")
		}
		monthTimestamp := utils.MonthDate(timestamp)
		if _, ok := m[nodeAddr][monthTimestamp]; !ok {
			m[nodeAddr][monthTimestamp] = NewStatistics("month")
		}
		yearTimestamp := utils.YearDate(timestamp)
		if _, ok := y[nodeAddr][yearTimestamp]; !ok {
			y[nodeAddr][yearTimestamp] = NewStatistics("year")
		}

		d[nodeAddr][dayTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		d[nodeAddr][dayTimestamp].SessionDuration[sessionID] = duration

		w[nodeAddr][weekTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		w[nodeAddr][weekTimestamp].SessionDuration[sessionID] = duration

		m[nodeAddr][monthTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		m[nodeAddr][monthTimestamp].SessionDuration[sessionID] = duration

		y[nodeAddr][yearTimestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		y[nodeAddr][yearTimestamp].SessionDuration[sessionID] = duration
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	for s := range d {
		var tKeys []time.Time
		for t := range d[s] {
			tKeys = append(tKeys, t)
		}

		sort.Slice(tKeys, func(i, j int) bool {
			return tKeys[i].After(tKeys[j])
		})

		for _, t := range tKeys {
			for u := range d[s][t].SessionBandwidth {
				if v, ok := d[s][t.AddDate(0, 0, -1)]; ok {
					if v, ok := v.SessionBandwidth[u]; ok {
						d[s][t].SessionBandwidth[u] = d[s][t].SessionBandwidth[u].Sub(v)
					}
					if v, ok := v.SessionDuration[u]; ok {
						d[s][t].SessionDuration[u] = d[s][t].SessionDuration[u] - v
					}
				}
			}
		}
	}

	for s := range w {
		var tKeys []time.Time
		for t := range w[s] {
			tKeys = append(tKeys, t)
		}

		sort.Slice(tKeys, func(i, j int) bool {
			return tKeys[i].After(tKeys[j])
		})

		for _, t := range tKeys {
			for u := range w[s][t].SessionBandwidth {
				if v, ok := w[s][t.AddDate(0, 0, -7)]; ok {
					if v, ok := v.SessionBandwidth[u]; ok {
						w[s][t].SessionBandwidth[u] = w[s][t].SessionBandwidth[u].Sub(v)
					}
					if v, ok := v.SessionDuration[u]; ok {
						w[s][t].SessionDuration[u] = w[s][t].SessionDuration[u] - v
					}
				}
			}
		}
	}

	for s := range m {
		var tKeys []time.Time
		for t := range m[s] {
			tKeys = append(tKeys, t)
		}

		sort.Slice(tKeys, func(i, j int) bool {
			return tKeys[i].After(tKeys[j])
		})

		for _, t := range tKeys {
			for u := range m[s][t].SessionBandwidth {
				if v, ok := m[s][t.AddDate(0, -1, 0)]; ok {
					if v, ok := v.SessionBandwidth[u]; ok {
						m[s][t].SessionBandwidth[u] = m[s][t].SessionBandwidth[u].Sub(v)
					}
					if v, ok := v.SessionDuration[u]; ok {
						m[s][t].SessionDuration[u] = m[s][t].SessionDuration[u] - v
					}
				}
			}
		}
	}

	for s := range y {
		var tKeys []time.Time
		for t := range y[s] {
			tKeys = append(tKeys, t)
		}

		sort.Slice(tKeys, func(i, j int) bool {
			return tKeys[i].After(tKeys[j])
		})

		for _, t := range tKeys {
			for u := range y[s][t].SessionBandwidth {
				if v, ok := y[s][t.AddDate(-1, 0, 0)]; ok {
					if v, ok := v.SessionBandwidth[u]; ok {
						y[s][t].SessionBandwidth[u] = y[s][t].SessionBandwidth[u].Sub(v)
					}
					if v, ok := v.SessionDuration[u]; ok {
						y[s][t].SessionDuration[u] = y[s][t].SessionDuration[u] - v
					}
				}
			}
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
