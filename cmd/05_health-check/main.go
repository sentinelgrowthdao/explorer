package main

import (
	"context"
	"flag"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"

	"github.com/sentinel-official/explorer/database"
	nodetypes "github.com/sentinel-official/explorer/types/node"
	"github.com/sentinel-official/explorer/utils"
)

const appName = "05_health-check"

var (
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
	timeout    time.Duration
)

func init() {
	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.DurationVar(&timeout, "timeout", 15*time.Second, "")
	flag.Parse()
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "addr", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
		{
			Keys: bson.D{
				bson.E{Key: "status", Value: 1},
				bson.E{Key: "remote_url", Value: 1},
			},
		},
	}

	_, err := database.NodeIndexesCreateMany(ctx, db, indexes)
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

	if err := createIndexes(context.TODO(), db); err != nil {
		log.Panicln(err)
	}

	filter := bson.M{
		"remote_url": bson.M{
			"$exists": true,
		},
		"status": "active",
	}
	projection := bson.M{
		"addr":       1,
		"remote_url": 1,
	}
	opts := options.Find().
		SetProjection(projection)

	nodes, err := database.NodeFind(context.TODO(), db, filter, opts)
	if err != nil {
		log.Panicln(err)
	}

	group := &errgroup.Group{}
	group.SetLimit(64)

	for i := 0; i < len(nodes); i++ {
		var (
			nodeAddr  = nodes[i].Addr
			remoteURL = nodes[i].RemoteURL
		)

		group.Go(func() error {
			filter := bson.M{
				"addr": nodeAddr,
			}
			update := bson.M{}
			projection := bson.M{
				"_id": 1,
			}
			opts := options.FindOneAndUpdate().
				SetProjection(projection)

			info, err := nodetypes.FetchNewInfo(remoteURL, timeout)
			if err != nil {
				update = bson.M{
					"$set": bson.M{
						"health.status_fetch_error":     err.Error(),
						"health.status_fetch_timestamp": time.Now().UTC(),
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"handshake_dns":                 info.Handshake,
						"health.status_fetch_error":     "",
						"health.status_fetch_timestamp": time.Now().UTC(),
						"internet_speed": bson.M{
							"download": utils.MustStringFromInt64(info.Bandwidth.Download),
							"upload":   utils.MustStringFromInt64(info.Bandwidth.Upload),
						},
						"interval_set_sessions":    info.IntervalSetSessions,
						"interval_update_sessions": info.IntervalUpdateStatus,
						"interval_update_status":   info.IntervalUpdateStatus,
						"location":                 info.Location,
						"moniker":                  info.Moniker,
						"peers":                    info.Peers,
						"qos":                      info.QOS,
						"type":                     info.Type,
						"version":                  info.Version,
					},
				}
			}

			_, err = database.NodeFindOneAndUpdate(context.TODO(), db, filter, update, opts)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		log.Panicln(err)
	}
}
