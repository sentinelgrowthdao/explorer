package main

import (
	"context"
	"flag"
	"log"
	"runtime"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "06_node-statistics"
)

var (
	batchSize    int
	dbAddress    string
	dbName       string
	dbUsername   string
	dbPassword   string
	excludeAddrs string
)

func init() {
	log.SetFlags(0)

	flag.IntVar(&batchSize, "batch-size", 25_000, "")
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.StringVar(&excludeAddrs, "exclude-addrs", "sent1c4nvz43tlw6d0c9nfu6r957y5d9pgjk5czl3n3", "")
	flag.Parse()
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "type", Value: 1},
				bson.E{Key: "timestamp", Value: -1},
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
				bson.E{Key: "id", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.SessionIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "addr", Value: 1},
				bson.E{Key: "timeframe", Value: 1},
				bson.E{Key: "timestamp", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	}

	_, err = database.NodeStatisticIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := utils.PrepareDatabase(context.TODO(), appName, dbUsername, dbPassword, dbAddress, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Client().Ping(context.TODO(), nil); err != nil {
		log.Fatalln(err)
	}

	now := time.Now()
	if err := createIndexes(context.TODO(), db); err != nil {
		log.Fatalln(err)
	}

	filter := bson.M{}
	projection := bson.M{
		"_id":    0,
		"height": 1,
		"time":   1,
	}
	opts := options.Find().
		SetProjection(projection).
		SetSort(bson.D{
			bson.E{Key: "height", Value: -1},
		}).
		SetLimit(1)

	dBlocks, err := database.BlockFind(context.TODO(), db, filter, opts)
	if err != nil {
		log.Fatalln(err)
	}

	maxTimestamp := time.Now().UTC()
	if len(dBlocks) > 0 {
		maxTimestamp = dBlocks[0].Time
	}

	excludeAddrs := strings.Split(excludeAddrs, ",")
	sort.Strings(excludeAddrs)

	var (
		models []mongo.WriteModel
		group  = errgroup.Group{}
	)

	addModels := func(m []bson.M) error {
		for i := 0; i < len(m); i++ {
			filter := bson.M{
				"addr":      m[i]["addr"],
				"timeframe": m[i]["timeframe"],
				"timestamp": m[i]["timestamp"],
			}
			update := bson.M{
				"$set": m[i],
			}
			model := mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(update).
				SetUpsert(true)

			models = append(models, model)
		}

		return nil
	}

	group.Go(func() error {
		defer func() { defer runtime.GC() }()

		m, err := StatisticsFromEvents(context.TODO(), db)
		if err != nil {
			return err
		}

		return addModels(m)
	})

	group.Go(func() error {
		defer func() { defer runtime.GC() }()

		m, err := StatisticsFromSessions(context.TODO(), db, time.Time{}, maxTimestamp, excludeAddrs)
		if err != nil {
			return err
		}

		return addModels(m)
	})

	group.Go(func() error {
		defer func() { defer runtime.GC() }()

		m, err := StatisticsFromSubscriptions(context.TODO(), db, time.Time{}, maxTimestamp, excludeAddrs)
		if err != nil {
			return err
		}

		return addModels(m)
	})

	group.Go(func() error {
		defer func() { defer runtime.GC() }()

		m, err := StatisticsFromSubscriptionPayouts(context.TODO(), db)
		if err != nil {
			return err
		}

		return addModels(m)
	})

	if err := group.Wait(); err != nil {
		log.Fatalln(err)
	}

	group = errgroup.Group{}
	group.SetLimit(runtime.NumCPU() / 2)

	log.Println("Models", len(models))
	for i := 0; ; {
		from, to := i, i+batchSize
		if to > len(models) {
			to = len(models)
		}

		group.Go(func() error {
			defer func() { defer runtime.GC() }()

			opts := options.BulkWrite().
				SetBypassDocumentValidation(false).
				SetOrdered(false)

			_, err = database.NodeStatisticBulkWrite(context.TODO(), db, models[from:to], opts)
			return err
		})

		i = i + batchSize
		if to == len(models) {
			break
		}
	}

	if err := group.Wait(); err != nil {
		log.Fatalln(err)
	}

	log.Println("Duration", time.Since(now))
	log.Println("")
	if err != nil {
		log.Fatalln(err)
	}
}
