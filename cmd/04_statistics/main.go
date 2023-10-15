package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/explorer/database"
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
			},
		},
		{
			Keys: bson.D{
				bson.E{Key: "type", Value: 1},
				bson.E{Key: "status", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
			},
		},
	}

	_, err := database.EventIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "register_timestamp", Value: 1},
			},
		},
	}

	_, err = database.NodeIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "start_timestamp", Value: 1},
			},
		},
	}

	_, err = database.SessionIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "start_timestamp", Value: 1},
			},
		},
	}

	_, err = database.SubscriptionIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Panicln(err)
	}

	if err := db.Client().Ping(context.TODO(), nil); err != nil {
		log.Panicln(err)
	}

	now := time.Now()

	if err := createIndexes(context.TODO(), db); err != nil {
		log.Panicln(err)
	}

	filter := bson.M{}
	projection := bson.M{
		"_id":    0,
		"height": 1,
		"time":   1,
	}
	_sort := bson.D{
		bson.E{Key: "height", Value: -1},
	}

	dBlocks, err := database.BlockFind(context.TODO(), db, filter, options.Find().SetProjection(projection).SetSort(_sort).SetLimit(1))
	if err != nil {
		log.Panicln(err)
	}

	var maxTimestamp time.Time
	if len(dBlocks) > 0 {
		maxTimestamp = dBlocks[0].Time
	}

	events, err := StatisticsFromEvents(context.TODO(), db)
	if err != nil {
		log.Panicln(err)
	}

	nodes, err := StatisticsFromNodes(context.TODO(), db)
	if err != nil {
		log.Panicln(err)
	}

	sessions, err := StatisticsFromSessions(context.TODO(), db, time.Time{}, maxTimestamp)
	if err != nil {
		log.Panicln(err)
	}

	subscriptions, err := StatisticsFromSubscriptions(context.TODO(), db, time.Time{}, maxTimestamp)
	if err != nil {
		log.Panicln(err)
	}

	result := append([]bson.M{}, events...)
	result = append(result, nodes...)
	result = append(result, sessions...)
	result = append(result, subscriptions...)

	sort.Slice(result, func(i, j int) bool {
		return result[i]["timestamp"].(time.Time).After(result[j]["timestamp"].(time.Time))
	})

	fmt.Println(utils.MustMarshalIndentToString(result))

	log.Println("Duration", time.Since(now))
	log.Println("")
	if err != nil {
		log.Panicln(err)
	}
}
