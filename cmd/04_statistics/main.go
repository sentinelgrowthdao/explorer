package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "04_statistics"
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

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "type", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
				bson.E{Key: "status", Value: 1},
			},
			Options: options.Index().
				SetPartialFilterExpression(
					bson.M{
						"type":   types.EventTypeNodeUpdateStatus,
						"status": hubtypes.StatusActive.String(),
					},
				),
		},
	}

	_, err := database.EventIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	now := time.Now()

	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Panicln(err)
	}

	if err := db.Client().Ping(context.TODO(), nil); err != nil {
		log.Panicln(err)
	}

	if err := createIndexes(context.TODO(), db); err != nil {
		log.Panicln(err)
	}

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

	result, err := database.BlockFind(context.TODO(), db, filter, opts)
	if err != nil {
		log.Panicln(err)
	}
	if len(result) == 0 {
		log.Panicln("nil result")
	}

	fromTimestamp := result[0].Time
	if fromTimestamp.IsZero() {
		log.Panicln("zero fromTimestamp")
	}

	sort = bson.D{
		bson.E{Key: "height", Value: -1},
	}
	opts = options.Find().
		SetProjection(project).
		SetSort(sort).
		SetLimit(1)

	result, err = database.BlockFind(context.TODO(), db, filter, opts)
	if err != nil {
		log.Panicln(err)
	}
	if len(result) == 0 {
		log.Panicln("nil result")
	}

	toTimestamp := result[0].Time
	if toTimestamp.IsZero() {
		log.Panicln("zero toTimestamp")
	}

	dayActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "day")
	if err != nil {
		log.Panicln(err)
	}

	weekActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "week")
	if err != nil {
		log.Panicln(err)
	}

	monthActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "month")
	if err != nil {
		log.Panicln(err)
	}

	yearActiveNodes, err := HistoricalActiveNodesCount(context.TODO(), db, fromTimestamp, toTimestamp, "year")
	if err != nil {
		log.Panicln(err)
	}

	mArr := [][]bson.M{
		dayActiveNodes, weekActiveNodes, monthActiveNodes, yearActiveNodes,
	}

	var items []bson.M
	for _, arr := range mArr {
		items = append(items, arr...)
	}

	err = db.Client().UseSession(
		context.TODO(),
		func(ctx mongo.SessionContext) error {
			err := ctx.StartTransaction(
				options.Transaction().
					SetReadConcern(readconcern.Snapshot()).
					SetWriteConcern(writeconcern.Majority()),
			)
			if err != nil {
				return err
			}

			abort := true
			defer func() {
				if abort {
					_ = ctx.AbortTransaction(ctx)
				}
			}()

			filter := bson.M{}
			opts := options.Delete()

			if err := database.StatisticDeleteMany(ctx, db, filter, opts); err != nil {
				return err
			}

			for _, item := range items {
				filter := bson.M{
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

				if _, err := database.StatisticFindOneAndUpdate(ctx, db, filter, update, opts); err != nil {
					return err
				}
			}

			abort = false
			return ctx.CommitTransaction(ctx)
		},
	)

	log.Println("Duration", time.Since(now))
	if err != nil {
		log.Panicln(err)
	}
}

func HistoricalActiveNodesCount(ctx context.Context, db *mongo.Database, fromTimestamp, toTimestamp time.Time, timeframe string) ([]bson.M, error) {
	log.Println("HistoricalActiveNodesCount", fromTimestamp, toTimestamp, timeframe)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type": types.EventTypeNodeUpdateStatus,
				"timestamp": bson.M{
					"$gte": fromTimestamp,
					"$lt":  toTimestamp,
				},
				"status": hubtypes.StatusActive.String(),
			},
		},
		{
			"$project": bson.M{
				"_id":          0,
				"node_address": 1,
				"timestamp":    1,
			},
		},
		{
			"$group": bson.M{
				"_id": func() bson.M {
					b := bson.M{
						"address": "$node_address",
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
					"tag":       types.StatisticTypeActiveNode,
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

	result, err := database.EventAggregate(ctx, db, pipeline)
	if err != nil {
		return nil, err
	}

	return result, nil
}
