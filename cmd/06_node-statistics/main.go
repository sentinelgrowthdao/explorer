package main

import (
	"context"
	"flag"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "06_node-statistics"
)

var (
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	log.SetFlags(0)

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
				bson.E{Key: "timestamp", Value: -1},
			},
		},
	}

	_, err := database.EventIndexesCreateMany(ctx, db, indexes)
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
	_sort := bson.D{
		bson.E{Key: "height", Value: -1},
	}

	dBlocks, err := database.BlockFind(context.TODO(), db, filter, options.Find().SetProjection(projection).SetSort(_sort).SetLimit(1))
	if err != nil {
		log.Fatalln(err)
	}

	maxTimestamp := time.Now().UTC()
	if len(dBlocks) > 0 {
		maxTimestamp = dBlocks[0].Time
	}

	var (
		m     []bson.M
		group = errgroup.Group{}
	)

	group.Go(func() error {
		v, err := StatisticsFromEvents(context.TODO(), db)
		if err != nil {
			return err
		}

		m = append(m, v...)
		return nil
	})

	group.Go(func() error {
		v, err := StatisticsFromSessions(context.TODO(), db, time.Time{}, maxTimestamp)
		if err != nil {
			return err
		}

		m = append(m, v...)
		return nil
	})

	group.Go(func() error {
		v, err := StatisticsFromSubscriptions(context.TODO(), db, time.Time{}, maxTimestamp)
		if err != nil {
			return err
		}

		m = append(m, v...)
		return nil
	})

	group.Go(func() error {
		v, err := StatisticsFromSubscriptionPayouts(context.TODO(), db)
		if err != nil {
			return err
		}

		m = append(m, v...)
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Fatalln(err)
	}

	mm := make(map[string]map[time.Time]bson.M)
	for i := 0; i < len(m); i++ {
		addr := m[i]["addr"].(string)
		timestamp := m[i]["timestamp"].(time.Time)

		if _, ok := mm[addr]; !ok {
			mm[addr] = make(map[time.Time]bson.M)
		}
		if _, ok := mm[addr][timestamp]; !ok {
			mm[addr][timestamp] = make(bson.M)
		}

		for s, v := range m[i] {
			mm[addr][timestamp][s] = v
		}
	}

	var result bson.A
	for _, m := range mm {
		for _, item := range m {
			result = append(result, item)
		}
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
			if err := database.NodeStatisticDeleteMany(ctx, db, filter); err != nil {
				return err
			}

			if _, err := database.NodeStatisticInsertMany(ctx, db, result); err != nil {
				return err
			}

			abort = false
			return ctx.CommitTransaction(ctx)
		},
	)

	log.Println("Duration", time.Since(now))
	log.Println("")
	if err != nil {
		log.Fatalln(err)
	}
}
