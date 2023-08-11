package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	commontypes "github.com/sentinel-official/explorer/types/common"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "04-statistics"
)

var (
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func main() {
	db, err := database.PrepareDatabase(context.Background(), appName, dbAddress, dbUsername, dbPassword, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	now := time.Now()

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

	dayActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "day")
	if err != nil {
		log.Fatalln(err)
	}
	weekActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "week")
	if err != nil {
		log.Fatalln(err)
	}
	monthActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "month")
	if err != nil {
		log.Fatalln(err)
	}
	yearActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "year")
	if err != nil {
		log.Fatalln(err)
	}

	dayJoinNodes, err := HistoricalJoinNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "day")
	if err != nil {
		log.Fatalln(err)
	}
	weekJoinNodes, err := HistoricalJoinNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "week")
	if err != nil {
		log.Fatalln(err)
	}
	monthJoinNodes, err := HistoricalJoinNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "month")
	if err != nil {
		log.Fatalln(err)
	}
	yearJoinNodes, err := HistoricalJoinNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "year")
	if err != nil {
		log.Fatalln(err)
	}

	statisticsSession, _, _, err := HistoricalStatisticsSession(context.TODO(), db, fromTimestamp, toTimestamp)
	if err != nil {
		log.Println(err)
	}

	statisticsSessionEvents, err := HistoricalStatisticsSessionEvent(context.TODO(), db, fromTimestamp, toTimestamp)
	if err != nil {
		log.Println(err)
	}

	statisticsSubscription, err := HistoricalStatisticsSubscription(context.TODO(), db, fromTimestamp, toTimestamp)
	if err != nil {
		log.Println(err)
	}

	mArr := [][]bson.M{
		dayActiveNodes, weekActiveNodes, monthActiveNodes, yearActiveNodes,
		dayJoinNodes, weekJoinNodes, monthJoinNodes, yearJoinNodes,
		statisticsSession,
		statisticsSessionEvents,
		statisticsSubscription,
	}

	var items []bson.M
	for _, arr := range mArr {
		items = append(items, arr...)
	}

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "tag", Value: 1},
				bson.E{Key: "timeframe", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
			},
			Options: nil,
		},
	}

	if err := database.StatisticsIndexesCreateMany(context.TODO(), db, indexes); err != nil {
		log.Fatalln(err)
	}

	err = db.Client().UseSession(
		context.TODO(),
		func(sctx mongo.SessionContext) error {
			err = sctx.StartTransaction(
				options.Transaction().
					SetReadConcern(readconcern.Snapshot()).
					SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
			)
			if err != nil {
				return err
			}

			abort := true
			defer func() {
				if abort {
					_ = sctx.AbortTransaction(sctx)
				}
			}()

			filter := bson.M{}
			opts := options.Delete()

			if err = database.StatisticsDeleteMany(sctx, db, filter, opts); err != nil {
				return err
			}

			for _, item := range items {
				filter = bson.M{
					"tag":       item["tag"],
					"timeframe": item["timeframe"],
					"timestamp": item["timestamp"],
				}

				update := bson.M{
					"$set": bson.M{
						"value": item["value"],
					},
				}
				projection := bson.M{
					"_id": 1,
				}
				opts := options.FindOneAndUpdate().
					SetProjection(projection).
					SetUpsert(true)

				_, err = database.StatisticsFindOneAndUpdate(sctx, db, filter, update, opts)
				if err != nil {
					return err
				}
			}

			abort = false
			return sctx.CommitTransaction(sctx)
		},
	)

	log.Println("Duration", time.Since(now))
	if err != nil {
		log.Fatalln(err)
	}
}

func HistoricalActiveNodesCount(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time, timeframe string) ([]bson.M, error) {
	log.Println("HistoricalActiveNodesCount", fromTimestamp, toTimestamp, timeframe)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"status": "STATUS_ACTIVE",
				"timestamp": bson.M{
					"$gte": fromTimestamp,
					"$lt":  toTimestamp,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"address":   1,
				"timestamp": 1,
			},
		},
		{
			"$group": bson.M{
				"_id": func() bson.M {
					b := bson.M{
						"address": "$address",
					}

					if timeframe == "year" {
						b["year"] = bson.M{"$year": "$timestamp"}
						return b
					}
					if timeframe == "month" {
						b["year"], b["month"] = bson.M{"$year": "$timestamp"}, bson.M{"$month": "$timestamp"}
						return b
					}
					if timeframe == "week" {
						b["year"], b["week"] = bson.M{"$isoWeekYear": "$timestamp"}, bson.M{"$isoWeek": "$timestamp"}
						return b
					}
					if timeframe == "day" {
						b["year"], b["month"], b["day"] = bson.M{"$year": "$timestamp"}, bson.M{"$month": "$timestamp"}, bson.M{"$dayOfMonth": "$timestamp"}
						return b
					}

					panic(fmt.Errorf("invalid timeframe %s", timeframe))
				}(),
				"value": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": func() bson.M {
					b := bson.M{}

					if timeframe == "year" {
						b["year"] = "$_id.year"
						return b
					}
					if timeframe == "month" {
						b["year"], b["month"] = "$_id.year", "$_id.month"
						return b
					}
					if timeframe == "week" {
						b["year"], b["week"] = "$_id.year", "$_id.week"
						return b
					}
					if timeframe == "day" {
						b["year"], b["month"], b["day"] = "$_id.year", "$_id.month", "$_id.day"
						return b
					}

					panic(fmt.Errorf("invalid timeframe %s", timeframe))
				}(),
				"value": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$project": func() bson.M {
				b := bson.M{
					"_id":       0,
					"tag":       types.TagActiveNode,
					"timeframe": timeframe,
					"value":     "$value",
				}

				if timeframe == "year" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"year": "$_id.year"}}
					return b
				}
				if timeframe == "month" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"year": "$_id.year", "month": "$_id.month"}}
					return b
				}
				if timeframe == "week" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"isoWeekYear": "$_id.year", "isoWeek": "$_id.week"}}
					return b
				}
				if timeframe == "day" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"year": "$_id.year", "month": "$_id.month", "day": "$_id.day"}}
					return b
				}

				panic(fmt.Errorf("invalid timeframe %s", timeframe))
			}(),
		},
	}

	result, err := database.NodeEventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func HistoricalJoinNodesCount(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time, timeframe string) ([]bson.M, error) {
	log.Println("HistoricalJoinNodesCount", fromTimestamp, toTimestamp, timeframe)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"join_timestamp": bson.M{
					"$gte": fromTimestamp,
					"$lt":  toTimestamp,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":            0,
				"join_timestamp": 1,
			},
		},
		{
			"$group": bson.M{
				"_id": func() bson.M {
					b := bson.M{}

					if timeframe == "year" {
						b["year"] = bson.M{"$year": "$join_timestamp"}
						return b
					}
					if timeframe == "month" {
						b["year"], b["month"] = bson.M{"$year": "$join_timestamp"}, bson.M{"$month": "$join_timestamp"}
						return b
					}
					if timeframe == "week" {
						b["year"], b["week"] = bson.M{"$isoWeekYear": "$join_timestamp"}, bson.M{"$isoWeek": "$join_timestamp"}
						return b
					}
					if timeframe == "day" {
						b["year"], b["month"], b["day"] = bson.M{"$year": "$join_timestamp"}, bson.M{"$month": "$join_timestamp"}, bson.M{"$dayOfMonth": "$join_timestamp"}
						return b
					}

					panic(fmt.Errorf("invalid timeframe %s", timeframe))
				}(),
				"value": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$project": func() bson.M {
				b := bson.M{
					"_id":       0,
					"tag":       types.TagJoinNode,
					"timeframe": timeframe,
					"value":     "$value",
				}

				if timeframe == "year" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"year": "$_id.year"}}
					return b
				}
				if timeframe == "month" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"year": "$_id.year", "month": "$_id.month"}}
					return b
				}
				if timeframe == "week" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"isoWeekYear": "$_id.year", "isoWeek": "$_id.week"}}
					return b
				}
				if timeframe == "day" {
					b["timestamp"] = bson.M{"$dateFromParts": bson.M{"year": "$_id.year", "month": "$_id.month", "day": "$_id.day"}}
					return b
				}

				panic(fmt.Errorf("invalid timeframe %s", timeframe))
			}(),
		},
	}

	result, err := database.NodeAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func HistoricalStatisticsSession(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time) (result, addrResult, nodeAddrResult []bson.M, err error) {
	log.Println("HistoricalStatisticsSession", fromTimestamp, toTimestamp)

	dayStartSessions := make(map[time.Time]int64)
	weekStartSessions := make(map[time.Time]int64)
	monthStartSessions := make(map[time.Time]int64)
	yearStartSessions := make(map[time.Time]int64)

	dayEndSessions := make(map[time.Time]int64)
	weekEndSessions := make(map[time.Time]int64)
	monthEndSessions := make(map[time.Time]int64)
	yearEndSessions := make(map[time.Time]int64)

	dayActiveSessions := make(map[time.Time]int64)
	weekActiveSessions := make(map[time.Time]int64)
	monthActiveSessions := make(map[time.Time]int64)
	yearActiveSessions := make(map[time.Time]int64)

	dayPayments := make(map[time.Time]sdk.Coins)
	weekPayments := make(map[time.Time]sdk.Coins)
	monthPayments := make(map[time.Time]sdk.Coins)
	yearPayments := make(map[time.Time]sdk.Coins)

	dayStakingRewards := make(map[time.Time]sdk.Coins)
	weekStakingRewards := make(map[time.Time]sdk.Coins)
	monthStakingRewards := make(map[time.Time]sdk.Coins)
	yearStakingRewards := make(map[time.Time]sdk.Coins)

	nodeAddrDayStartSessions := make(map[string]map[time.Time]int64)
	nodeAddrDayEndSessions := make(map[string]map[time.Time]int64)
	nodeAddrDayActiveSessions := make(map[string]map[time.Time]int64)
	nodeAddrDayPayments := make(map[string]map[time.Time]sdk.Coins)
	nodeAddrDayStakingRewards := make(map[string]map[time.Time]sdk.Coins)

	addrDayStartSessions := make(map[string]map[time.Time]int64)
	addrDayEndSessions := make(map[string]map[time.Time]int64)
	addrDayActiveSessions := make(map[string]map[time.Time]int64)
	addrDayPayments := make(map[string]map[time.Time]sdk.Coins)
	addrDayStakingRewards := make(map[string]map[time.Time]sdk.Coins)

	year, month, day := fromTimestamp.Date()
	fromTimestamp = time.Date(year, month, day, 0, 0, 0, 0, fromTimestamp.Location())

	year, month, day = toTimestamp.Date()
	toTimestamp = time.Date(year, month, day, 0, 0, 0, 0, toTimestamp.Location())

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"address":         1,
		"node_address":    1,
		"start_timestamp": 1,
		"end_timestamp":   1,
		"staking_reward":  1,
		"payment":         1,
	}
	opts := options.Find().
		SetProjection(projection)

	items, err := database.SessionFindAll(ctx, db, filter, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, item := range items {
		startTimestamp := item.StartTimestamp
		if startTimestamp.IsZero() {
			startTimestamp = fromTimestamp
		}

		endTimestamp := item.EndTimestamp
		if endTimestamp.IsZero() {
			endTimestamp = toTimestamp
		}

		if _, ok := addrDayStartSessions[item.Address]; !ok {
			addrDayStartSessions[item.Address] = make(map[time.Time]int64)
		}
		if _, ok := addrDayActiveSessions[item.Address]; !ok {
			addrDayActiveSessions[item.Address] = make(map[time.Time]int64)
		}
		if _, ok := addrDayEndSessions[item.Address]; !ok {
			addrDayEndSessions[item.Address] = make(map[time.Time]int64)
		}
		if _, ok := addrDayPayments[item.Address]; !ok {
			addrDayPayments[item.Address] = make(map[time.Time]sdk.Coins)
		}
		if _, ok := addrDayStakingRewards[item.Address]; !ok {
			addrDayStakingRewards[item.Address] = make(map[time.Time]sdk.Coins)
		}

		if _, ok := nodeAddrDayStartSessions[item.NodeAddress]; !ok {
			nodeAddrDayStartSessions[item.NodeAddress] = make(map[time.Time]int64)
		}
		if _, ok := nodeAddrDayActiveSessions[item.NodeAddress]; !ok {
			nodeAddrDayActiveSessions[item.NodeAddress] = make(map[time.Time]int64)
		}
		if _, ok := nodeAddrDayEndSessions[item.NodeAddress]; !ok {
			nodeAddrDayEndSessions[item.NodeAddress] = make(map[time.Time]int64)
		}
		if _, ok := nodeAddrDayPayments[item.NodeAddress]; !ok {
			nodeAddrDayPayments[item.NodeAddress] = make(map[time.Time]sdk.Coins)
		}
		if _, ok := nodeAddrDayStakingRewards[item.NodeAddress]; !ok {
			nodeAddrDayStakingRewards[item.NodeAddress] = make(map[time.Time]sdk.Coins)
		}

		dayStartTimestamp := utils.DayDate(startTimestamp)
		dayEndTimestamp := utils.DayDate(endTimestamp)
		for t := dayStartTimestamp; !t.After(dayEndTimestamp); t = t.AddDate(0, 0, 1) {
			dayActiveSessions[t] += 1

			addrDayActiveSessions[item.Address][t] += 1
			nodeAddrDayActiveSessions[item.NodeAddress][t] += 1
		}

		weekStartTimestamp := utils.ISOWeekDate(startTimestamp)
		weekEndTimestamp := utils.ISOWeekDate(endTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			weekActiveSessions[t] += 1
		}

		monthStartTimestamp := utils.MonthDate(startTimestamp)
		monthEndTimestamp := utils.MonthDate(endTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			monthActiveSessions[t] += 1
		}

		yearStartTimestamp := utils.YearDate(startTimestamp)
		yearEndTimestamp := utils.YearDate(endTimestamp)
		for t := yearStartTimestamp; !t.After(yearEndTimestamp); t = t.AddDate(1, 0, 0) {
			yearActiveSessions[t] += 1
		}

		if !item.StartTimestamp.IsZero() {
			dayStartSessions[dayStartTimestamp] += 1
			weekStartSessions[weekStartTimestamp] += 1
			monthStartSessions[monthStartTimestamp] += 1
			yearStartSessions[yearStartTimestamp] += 1

			addrDayStartSessions[item.Address][dayStartTimestamp] += 1
			nodeAddrDayStartSessions[item.NodeAddress][dayStartTimestamp] += 1
		}
		if !item.EndTimestamp.IsZero() {
			dayEndSessions[dayEndTimestamp] += 1
			weekEndSessions[weekEndTimestamp] += 1
			monthEndSessions[monthEndTimestamp] += 1
			yearEndSessions[yearEndTimestamp] += 1

			addrDayEndSessions[item.Address][dayEndTimestamp] += 1
			nodeAddrDayEndSessions[item.NodeAddress][dayEndTimestamp] += 1
		}
		if item.Payment != nil {
			paymentRaw := item.Payment.Raw()
			dayPayments[dayEndTimestamp] = dayPayments[dayEndTimestamp].Add(paymentRaw)
			weekPayments[weekEndTimestamp] = weekPayments[weekEndTimestamp].Add(paymentRaw)
			monthPayments[monthEndTimestamp] = monthPayments[monthEndTimestamp].Add(paymentRaw)
			yearPayments[yearEndTimestamp] = yearPayments[yearEndTimestamp].Add(paymentRaw)

			addrDayPayments[item.Address][dayEndTimestamp] = addrDayPayments[item.Address][dayEndTimestamp].Add(paymentRaw)
			nodeAddrDayPayments[item.NodeAddress][dayEndTimestamp] = nodeAddrDayPayments[item.NodeAddress][dayEndTimestamp].Add(paymentRaw)
		}
		if item.StakingReward != nil {
			stakingRewardRaw := item.StakingReward.Raw()
			dayStakingRewards[dayEndTimestamp] = dayStakingRewards[dayEndTimestamp].Add(stakingRewardRaw)
			weekStakingRewards[weekEndTimestamp] = weekStakingRewards[weekEndTimestamp].Add(stakingRewardRaw)
			monthStakingRewards[monthEndTimestamp] = monthStakingRewards[monthEndTimestamp].Add(stakingRewardRaw)
			yearStakingRewards[yearEndTimestamp] = yearStakingRewards[yearEndTimestamp].Add(stakingRewardRaw)

			addrDayStakingRewards[item.Address][dayEndTimestamp] = addrDayStakingRewards[item.Address][dayEndTimestamp].Add(stakingRewardRaw)
			nodeAddrDayStakingRewards[item.NodeAddress][dayEndTimestamp] = nodeAddrDayStakingRewards[item.NodeAddress][dayEndTimestamp].Add(stakingRewardRaw)
		}
	}

	for timestamp, value := range dayStartSessions {
		result = append(result, bson.M{
			"tag":       types.TagStartSession,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekStartSessions {
		result = append(result, bson.M{
			"tag":       types.TagStartSession,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthStartSessions {
		result = append(result, bson.M{
			"tag":       types.TagStartSession,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearStartSessions {
		result = append(result, bson.M{
			"tag":       types.TagStartSession,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	for timestamp, value := range dayActiveSessions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSession,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekActiveSessions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSession,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthActiveSessions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSession,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearActiveSessions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSession,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	for timestamp, value := range dayEndSessions {
		result = append(result, bson.M{
			"tag":       types.TagEndSession,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekEndSessions {
		result = append(result, bson.M{
			"tag":       types.TagEndSession,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthEndSessions {
		result = append(result, bson.M{
			"tag":       types.TagEndSession,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearEndSessions {
		result = append(result, bson.M{
			"tag":       types.TagEndSession,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	for timestamp, coins := range dayPayments {
		result = append(result, bson.M{
			"tag":       types.TagSessionPayment,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range weekPayments {
		result = append(result, bson.M{
			"tag":       types.TagSessionPayment,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range monthPayments {
		result = append(result, bson.M{
			"tag":       types.TagSessionPayment,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range yearPayments {
		result = append(result, bson.M{
			"tag":       types.TagSessionPayment,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}

	for timestamp, coins := range dayStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSessionStakingReward,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range weekStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSessionStakingReward,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range monthStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSessionStakingReward,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range yearStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSessionStakingReward,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}

	for addr, dayStartSessions := range addrDayStartSessions {
		for timestamp, value := range dayStartSessions {
			addrResult = append(addrResult, bson.M{
				"tag":       types.TagStartSession,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   addr,
				"value":     value,
			})
		}
	}
	for addr, dayActiveSessions := range addrDayActiveSessions {
		for timestamp, value := range dayActiveSessions {
			addrResult = append(addrResult, bson.M{
				"tag":       types.TagActiveSession,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   addr,
				"value":     value,
			})
		}
	}
	for addr, dayEndSessions := range addrDayEndSessions {
		for timestamp, value := range dayEndSessions {
			addrResult = append(addrResult, bson.M{
				"tag":       types.TagEndSession,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   addr,
				"value":     value,
			})
		}
	}
	for addr, dayPayments := range addrDayPayments {
		for timestamp, coins := range dayPayments {
			addrResult = append(addrResult, bson.M{
				"tag":       types.TagSessionPayment,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   addr,
				"value":     commontypes.NewCoinsFromRaw(coins),
			})
		}
	}
	for addr, dayStakingRewards := range addrDayStakingRewards {
		for timestamp, coins := range dayStakingRewards {
			addrResult = append(addrResult, bson.M{
				"tag":       types.TagSessionStakingReward,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   addr,
				"value":     commontypes.NewCoinsFromRaw(coins),
			})
		}
	}

	for nodeAddr, dayStartSessions := range nodeAddrDayStartSessions {
		for timestamp, value := range dayStartSessions {
			nodeAddrResult = append(nodeAddrResult, bson.M{
				"tag":       types.TagStartSession,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   nodeAddr,
				"value":     value,
			})
		}
	}
	for nodeAddr, dayActiveSessions := range nodeAddrDayActiveSessions {
		for timestamp, value := range dayActiveSessions {
			nodeAddrResult = append(nodeAddrResult, bson.M{
				"tag":       types.TagActiveSession,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   nodeAddr,
				"value":     value,
			})
		}
	}
	for nodeAddr, dayEndSessions := range nodeAddrDayEndSessions {
		for timestamp, value := range dayEndSessions {
			nodeAddrResult = append(nodeAddrResult, bson.M{
				"tag":       types.TagEndSession,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   nodeAddr,
				"value":     value,
			})
		}
	}
	for nodeAddr, dayPayments := range nodeAddrDayPayments {
		for timestamp, coins := range dayPayments {
			nodeAddrResult = append(nodeAddrResult, bson.M{
				"tag":       types.TagSessionPayment,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   nodeAddr,
				"value":     commontypes.NewCoinsFromRaw(coins),
			})
		}
	}
	for nodeAddr, dayStakingRewards := range nodeAddrDayStakingRewards {
		for timestamp, coins := range dayStakingRewards {
			nodeAddrResult = append(nodeAddrResult, bson.M{
				"tag":       types.TagSessionStakingReward,
				"timeframe": "day",
				"timestamp": timestamp,
				"address":   nodeAddr,
				"value":     commontypes.NewCoinsFromRaw(coins),
			})
		}
	}

	return result, addrResult, nodeAddrResult, nil
}

func HistoricalStatisticsSessionEvent(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time) (result []bson.M, err error) {
	log.Println("HistoricalStatisticsSessionEvent", fromTimestamp, toTimestamp)

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

func HistoricalStatisticsSubscription(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time) (result []bson.M, err error) {
	log.Println("HistoricalStatisticsSubscription", fromTimestamp, toTimestamp)

	dayStartSubscriptions := make(map[time.Time]int64)
	weekStartSubscriptions := make(map[time.Time]int64)
	monthStartSubscriptions := make(map[time.Time]int64)
	yearStartSubscriptions := make(map[time.Time]int64)

	dayEndSubscriptions := make(map[time.Time]int64)
	weekEndSubscriptions := make(map[time.Time]int64)
	monthEndSubscriptions := make(map[time.Time]int64)
	yearEndSubscriptions := make(map[time.Time]int64)

	dayActiveSubscriptions := make(map[time.Time]int64)
	weekActiveSubscriptions := make(map[time.Time]int64)
	monthActiveSubscriptions := make(map[time.Time]int64)
	yearActiveSubscriptions := make(map[time.Time]int64)

	dayDeposits := make(map[time.Time]sdk.Coins)
	weekDeposits := make(map[time.Time]sdk.Coins)
	monthDeposits := make(map[time.Time]sdk.Coins)
	yearDeposits := make(map[time.Time]sdk.Coins)

	dayPayments := make(map[time.Time]sdk.Coins)
	weekPayments := make(map[time.Time]sdk.Coins)
	monthPayments := make(map[time.Time]sdk.Coins)
	yearPayments := make(map[time.Time]sdk.Coins)

	dayStakingRewards := make(map[time.Time]sdk.Coins)
	weekStakingRewards := make(map[time.Time]sdk.Coins)
	monthStakingRewards := make(map[time.Time]sdk.Coins)
	yearStakingRewards := make(map[time.Time]sdk.Coins)

	year, month, day := fromTimestamp.Date()
	fromTimestamp = time.Date(year, month, day, 0, 0, 0, 0, fromTimestamp.Location())

	year, month, day = toTimestamp.Date()
	toTimestamp = time.Date(year, month, day, 0, 0, 0, 0, toTimestamp.Location())

	filter := bson.M{}
	projection := bson.M{
		"_id":             0,
		"start_timestamp": 1,
		"end_timestamp":   1,
		"staking_reward":  1,
		"payment":         1,
		"deposit":         1,
	}
	opts := options.Find().
		SetProjection(projection)

	items, err := database.SubscriptionFindAll(ctx, db, filter, opts)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		startTimestamp := item.StartTimestamp
		if startTimestamp.IsZero() {
			startTimestamp = fromTimestamp
		}

		endTimestamp := item.EndTimestamp
		if endTimestamp.IsZero() {
			endTimestamp = toTimestamp
		}

		dayStartTimestamp := utils.DayDate(startTimestamp)
		dayEndTimestamp := utils.DayDate(endTimestamp)
		for t := dayStartTimestamp; !t.After(dayEndTimestamp); t = t.AddDate(0, 0, 1) {
			dayActiveSubscriptions[t] += 1
		}

		weekStartTimestamp := utils.ISOWeekDate(startTimestamp)
		weekEndTimestamp := utils.ISOWeekDate(endTimestamp)
		for t := weekStartTimestamp; !t.After(weekEndTimestamp); t = t.AddDate(0, 0, 7) {
			weekActiveSubscriptions[t] += 1
		}

		monthStartTimestamp := utils.MonthDate(startTimestamp)
		monthEndTimestamp := utils.MonthDate(endTimestamp)
		for t := monthStartTimestamp; !t.After(monthEndTimestamp); t = t.AddDate(0, 1, 0) {
			monthActiveSubscriptions[t] += 1
		}

		yearStartTimestamp := utils.YearDate(startTimestamp)
		yearEndTimestamp := utils.YearDate(endTimestamp)
		for t := yearStartTimestamp; !t.After(yearEndTimestamp); t = t.AddDate(1, 0, 0) {
			yearActiveSubscriptions[t] += 1
		}

		if !item.StartTimestamp.IsZero() {
			dayStartSubscriptions[dayStartTimestamp] += 1
			weekStartSubscriptions[weekStartTimestamp] += 1
			monthStartSubscriptions[monthStartTimestamp] += 1
			yearStartSubscriptions[yearStartTimestamp] += 1
		}
		if !item.EndTimestamp.IsZero() {
			dayEndSubscriptions[dayEndTimestamp] += 1
			weekEndSubscriptions[weekEndTimestamp] += 1
			monthEndSubscriptions[monthEndTimestamp] += 1
			yearEndSubscriptions[yearEndTimestamp] += 1
		}
		if item.Deposit != nil {
			depositRaw := item.Deposit.Raw()
			dayDeposits[dayStartTimestamp] = dayDeposits[dayStartTimestamp].Add(depositRaw)
			weekDeposits[weekStartTimestamp] = weekDeposits[weekStartTimestamp].Add(depositRaw)
			monthDeposits[monthStartTimestamp] = monthDeposits[monthStartTimestamp].Add(depositRaw)
			yearDeposits[yearStartTimestamp] = yearDeposits[yearStartTimestamp].Add(depositRaw)
		}
		if item.Payment != nil {
			paymentRaw := item.Payment.Raw()
			dayPayments[dayStartTimestamp] = dayPayments[dayStartTimestamp].Add(paymentRaw)
			weekPayments[weekStartTimestamp] = weekPayments[weekStartTimestamp].Add(paymentRaw)
			monthPayments[monthStartTimestamp] = monthPayments[monthStartTimestamp].Add(paymentRaw)
			yearPayments[yearStartTimestamp] = yearPayments[yearStartTimestamp].Add(paymentRaw)
		}
		if item.StakingReward != nil {
			stakingRewardRaw := item.StakingReward.Raw()
			dayStakingRewards[dayEndTimestamp] = dayStakingRewards[dayEndTimestamp].Add(stakingRewardRaw)
			weekStakingRewards[weekEndTimestamp] = weekStakingRewards[weekEndTimestamp].Add(stakingRewardRaw)
			monthStakingRewards[monthEndTimestamp] = monthStakingRewards[monthEndTimestamp].Add(stakingRewardRaw)
			yearStakingRewards[yearEndTimestamp] = yearStakingRewards[yearEndTimestamp].Add(stakingRewardRaw)
		}
	}

	for timestamp, value := range dayStartSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagStartSubscription,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekStartSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagStartSubscription,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthStartSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagStartSubscription,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearStartSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagStartSubscription,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	for timestamp, value := range dayActiveSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSubscription,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekActiveSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSubscription,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthActiveSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSubscription,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearActiveSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagActiveSubscription,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	for timestamp, value := range dayEndSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagEndSubscription,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range weekEndSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagEndSubscription,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range monthEndSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagEndSubscription,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     value,
		})
	}
	for timestamp, value := range yearEndSubscriptions {
		result = append(result, bson.M{
			"tag":       types.TagEndSubscription,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     value,
		})
	}

	for timestamp, coins := range dayDeposits {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionDeposit,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range weekDeposits {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionDeposit,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range monthDeposits {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionDeposit,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range yearDeposits {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionDeposit,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}

	for timestamp, coins := range dayPayments {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionPayment,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range weekPayments {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionPayment,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range monthPayments {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionPayment,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range yearPayments {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionPayment,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}

	for timestamp, coins := range dayStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionStakingReward,
			"timeframe": "day",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range weekStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionStakingReward,
			"timeframe": "week",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range monthStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionStakingReward,
			"timeframe": "month",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}
	for timestamp, coins := range yearStakingRewards {
		result = append(result, bson.M{
			"tag":       types.TagSubscriptionStakingReward,
			"timeframe": "year",
			"timestamp": timestamp,
			"value":     commontypes.NewCoinsFromRaw(coins),
		})
	}

	return result, nil
}
