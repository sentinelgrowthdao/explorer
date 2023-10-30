package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	EventStatistics struct {
		SessionBandwidth map[uint64]*types.Bandwidth
		SessionDuration  map[uint64]int64
	}
)

func NewStatistics() *EventStatistics {
	return &EventStatistics{
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

	return bson.M{
		"addr":              addr,
		"timestamp":         timestamp,
		"session_bandwidth": sessionBandwidth,
		"session_duration":  sessionDuration,
	}
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

	items, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	d := make(map[string]map[time.Time]*EventStatistics)
	for i := 0; i < len(items); i++ {
		bandwidth := types.BandwidthFromInterface(items[i]["bandwidth"])
		duration := types.Int64FromInterface(items[i]["duration"])
		nodeAddr := types.StringFromInterface(items[i]["node_addr"])
		sessionID := types.Uint64FromInterface(items[i]["session_id"])
		timestamp := types.TimeFromInterface(items[i]["timestamp"])

		if _, ok := d[nodeAddr]; !ok {
			d[nodeAddr] = make(map[time.Time]*EventStatistics)
		}

		timestamp = utils.DayDate(timestamp)
		if _, ok := d[nodeAddr][timestamp]; !ok {
			d[nodeAddr][timestamp] = NewStatistics()
		}

		d[nodeAddr][timestamp].SessionBandwidth[sessionID] = bandwidth.Copy()
		d[nodeAddr][timestamp].SessionDuration[sessionID] = duration
	}

	for i := 0; i < len(items); i++ {
		nodeAddr := types.StringFromInterface(items[i]["node_addr"])
		sessionID := types.Uint64FromInterface(items[i]["session_id"])
		timestamp := types.TimeFromInterface(items[i]["timestamp"])

		timestamp = utils.DayDate(timestamp)
		if v, ok := d[nodeAddr][timestamp.AddDate(0, 0, -1)]; ok {
			if v, ok := v.SessionBandwidth[sessionID]; ok {
				d[nodeAddr][timestamp].SessionBandwidth[sessionID] = d[nodeAddr][timestamp].SessionBandwidth[sessionID].Sub(v)
			}
			if v, ok := v.SessionDuration[sessionID]; ok {
				d[nodeAddr][timestamp].SessionDuration[sessionID] = d[nodeAddr][timestamp].SessionDuration[sessionID] - v
			}
		}
	}

	for s, m := range d {
		for t, statistics := range m {
			result = append(result, statistics.Result(s, t))
		}
	}

	return result, nil
}
