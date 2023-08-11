package main

import (
	"context"
	"fmt"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	commontypes "github.com/sentinel-official/explorer/types/common"
	"github.com/sentinel-official/explorer/utils"
)

func HistoricalStatisticsSessionEvents(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time) (result []bson.M, err error) {
	log.Println("HistoricalStatisticsSessionEvents", fromTimestamp, toTimestamp)

	pipeline := []bson.M{
		{
			"$sort": bson.M{
				"timestamp": 1,
				"id":        1,
			},
		},
		{
			"$project": bson.M{
				"id":        1,
				"timestamp": 1,
				"bandwidth": 1,
				"duration":  1,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"id": "$id",
					"timestamp": bson.M{
						"$dateToString": bson.M{
							"format": "%Y-%m-%d",
							"date":   "$timestamp",
						},
					},
				},
				"bandwidth": bson.M{
					"$last": "$bandwidth",
				},
				"duration": bson.M{
					"$last": "$duration",
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"timestamp": "$_id.timestamp",
				"upload":    "$bandwidth.upload",
				"download":  "$bandwidth.download",
				"duration":  "$duration",
			},
		},
	}

	dayBandwidthConsumption := make(map[time.Time]hubtypes.Bandwidth)
	weekBandwidthConsumption := make(map[time.Time]hubtypes.Bandwidth)
	monthBandwidthConsumption := make(map[time.Time]hubtypes.Bandwidth)
	yearBandwidthConsumption := make(map[time.Time]hubtypes.Bandwidth)

	dayDurationSpent := make(map[time.Time]int64)
	weekDurationSpent := make(map[time.Time]int64)
	monthDurationSpent := make(map[time.Time]int64)
	yearDurationSpent := make(map[time.Time]int64)

	items, err := database.SessionEventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		timestamp, err := time.Parse("2006-01-02", item["timestamp"].(string))
		if err != nil {
			return nil, err
		}

		dayTimestamp := utils.DayDate(timestamp)
		weekTimestamp := utils.ISOWeekDate(timestamp)
		monthTimestamp := utils.MonthDate(timestamp)
		yearTimestamp := utils.YearDate(timestamp)

		if _, ok := dayBandwidthConsumption[dayTimestamp]; !ok {
			dayBandwidthConsumption[dayTimestamp] = hubtypes.NewBandwidthFromInt64(0, 0)
		}
		if _, ok := weekBandwidthConsumption[weekTimestamp]; !ok {
			weekBandwidthConsumption[weekTimestamp] = hubtypes.NewBandwidthFromInt64(0, 0)
		}
		if _, ok := monthBandwidthConsumption[monthTimestamp]; !ok {
			monthBandwidthConsumption[monthTimestamp] = hubtypes.NewBandwidthFromInt64(0, 0)
		}
		if _, ok := yearBandwidthConsumption[yearTimestamp]; !ok {
			yearBandwidthConsumption[yearTimestamp] = hubtypes.NewBandwidthFromInt64(0, 0)
		}

		if _, ok := item["upload"]; !ok {
			item["upload"] = int64(0)
		}
		if _, ok := item["download"]; !ok {
			item["download"] = int64(0)
		}

		bandwidth := hubtypes.NewBandwidthFromInt64(item["upload"].(int64), item["download"].(int64))
		duration := item["duration"].(int64)

		dayBandwidthConsumption[dayTimestamp] = dayBandwidthConsumption[dayTimestamp].Add(bandwidth)
		weekBandwidthConsumption[weekTimestamp] = weekBandwidthConsumption[weekTimestamp].Add(bandwidth)
		monthBandwidthConsumption[monthTimestamp] = monthBandwidthConsumption[monthTimestamp].Add(bandwidth)
		yearBandwidthConsumption[yearTimestamp] = yearBandwidthConsumption[yearTimestamp].Add(bandwidth)

		dayDurationSpent[dayTimestamp] = dayDurationSpent[dayTimestamp] + duration
		weekDurationSpent[weekTimestamp] = weekDurationSpent[weekTimestamp] + duration
		monthDurationSpent[monthTimestamp] = monthDurationSpent[monthTimestamp] + duration
		yearDurationSpent[yearTimestamp] = yearDurationSpent[yearTimestamp] + duration
	}

	for timestamp, value := range dayBandwidthConsumption {
		result = append(result, bson.M{
			"tag":       types.TagBandwidthConsumption,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     commontypes.NewBandwidthFromRaw(&value),
		})
	}
	for timestamp, value := range weekBandwidthConsumption {
		result = append(result, bson.M{
			"tag":       types.TagBandwidthConsumption,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     commontypes.NewBandwidthFromRaw(&value),
		})
	}
	for timestamp, value := range monthBandwidthConsumption {
		result = append(result, bson.M{
			"tag":       types.TagBandwidthConsumption,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     commontypes.NewBandwidthFromRaw(&value),
		})
	}
	for timestamp, value := range yearBandwidthConsumption {
		result = append(result, bson.M{
			"tag":       types.TagBandwidthConsumption,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     commontypes.NewBandwidthFromRaw(&value),
		})
	}

	for timestamp, value := range dayDurationSpent {
		result = append(result, bson.M{
			"tag":       types.TagSessionDuration,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekDurationSpent {
		result = append(result, bson.M{
			"tag":       types.TagSessionDuration,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthDurationSpent {
		result = append(result, bson.M{
			"tag":       types.TagSessionDuration,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearDurationSpent {
		result = append(result, bson.M{
			"tag":       types.TagSessionDuration,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	return result, nil
}

/* func main() {
	// log.SetOutput(io.Discard)

	db, err := database.PrepareDatabase(context.TODO(), "test", "mongodb://127.0.0.1:27017", "", "", "sentinelhub-2")
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.TODO(), nil); err != nil {
		log.Fatalln(err)
	}

	fromTimestamp := time.Time{}
	if fromTimestamp.IsZero() {
		filter := bson.M{}
		project := bson.M{
			"_id":  0,
			"time": 1,
		}
		sort := bson.D{
			bson.E{Key: "height", Value: 1},
		}
		opts := options.Find().
			SetProjection(project).
			SetSort(sort).
			SetLimit(1)

		result, err := database.BlockFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			log.Fatalln(err)
		}
		if len(result) == 0 {
			log.Fatalln("nil result")
		}

		fromTimestamp = result[0].Time
	}

	toTimestamp := time.Time{}
	if toTimestamp.IsZero() {
		filter := bson.M{}
		project := bson.M{
			"_id":  0,
			"time": 1,
		}
		sort := bson.D{
			bson.E{Key: "height", Value: -1},
		}
		opts := options.Find().
			SetProjection(project).
			SetSort(sort).
			SetLimit(1)

		result, err := database.BlockFindAll(context.TODO(), db, filter, opts)
		if err != nil {
			log.Fatalln(err)
		}
		if len(result) == 0 {
			log.Fatalln("nil result")
		}

		toTimestamp = result[0].Time
	}

	if fromTimestamp.IsZero() {
		log.Fatalln("fromTimestamp is zero")
	}
	if toTimestamp.IsZero() {
		log.Fatalln("toTimestamp is zero")
	}

	result1, err := HistoricalStatisticsSessionEvents(context.TODO(), db, fromTimestamp, toTimestamp)
	if err != nil {
		log.Fatalln(err)
	}

	sort.Slice(result1, func(i, j int) bool {
		if result1[i]["tag"].(string) != result1[j]["tag"].(string) {
			return result1[i]["tag"].(string) < result1[j]["tag"].(string)
		}
		if result1[i]["timeframe"].(string) != result1[j]["timeframe"].(string) {
			return result1[i]["timeframe"].(string) < result1[j]["timeframe"].(string)
		}
		if !result1[i]["timestamp"].(time.Time).Equal(result1[j]["timestamp"].(time.Time)) {
			return result1[i]["timestamp"].(time.Time).Before(result1[j]["timestamp"].(time.Time))
		}

		return false
	})

	// result2, err := HistoricalStatisticsSubscription(context.TODO(), db, fromTimestamp, toTimestamp)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	//
	// sort.Slice(result2, func(i, j int) bool {
	// 	if result2[i]["tag"].(string) != result2[j]["tag"].(string) {
	// 		return result2[i]["tag"].(string) < result2[j]["tag"].(string)
	// 	}
	// 	if result2[i]["timeframe"].(string) != result2[j]["timeframe"].(string) {
	// 		return result2[i]["timeframe"].(string) < result2[j]["timeframe"].(string)
	// 	}
	// 	if !result2[i]["timestamp"].(time.Time).Equal(result2[j]["timestamp"].(time.Time)) {
	// 		return result2[i]["timestamp"].(time.Time).Before(result2[j]["timestamp"].(time.Time))
	// 	}
	//
	// 	return result2[i]["value"].(int64) < result2[j]["value"].(int64)
	// })
	//
	for _, m := range result1 {
		fmt.Println(m["tag"], m["timeframe"], m["timestamp"], m["value"])
	}
	// for _, m := range result2 {
	// 	fmt.Println(m["tag"], m["timeframe"], m["timestamp"], m["value"])
	// }
} */

func main() {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"tag":       types.TagBandwidthConsumption,
				"timeframe": "day",
			},
		},
		{
			"$unwind": "$_id",
		},
		{
			"$group": bson.M{
				"_id": "",
				"download": bson.M{
					"$sum": "$value.download",
				},
				"upload": bson.M{
					"$sum": "$value.upload",
				},
			},
		},
		{
			"$sort": bson.D{
				bson.E{Key: "_id", Value: 1},
			},
		},
		{
			"$project": bson.M{
				"_id": "$_id",
				"value": bson.M{
					"download": "$download",
					"upload":   "$upload",
				},
			},
		},
	}

	fmt.Println(string(utils.MustMarshal(pipeline)))
}
