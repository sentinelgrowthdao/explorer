package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	blockapi "github.com/sentinel-official/explorer/api/block"
	depositapi "github.com/sentinel-official/explorer/api/deposit"
	nodeapi "github.com/sentinel-official/explorer/api/node"
	sessionapi "github.com/sentinel-official/explorer/api/session"
	statisticsapi "github.com/sentinel-official/explorer/api/statistics"
	subscriptionapi "github.com/sentinel-official/explorer/api/subscription"
	txapi "github.com/sentinel-official/explorer/api/tx"
	"github.com/sentinel-official/explorer/database"
	"github.com/sentinel-official/explorer/utils"
)

const (
	appName = "00_api-server"
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
				bson.E{Key: "status", Value: 1},
			},
		},
	}

	_, err := database.NodeIndexesCreateMany(ctx, db, indexes)
	if err != nil {
		return err
	}

	indexes = []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				bson.E{Key: "status", Value: 1},
				bson.E{Key: "acc_addr", Value: 1},
			},
		},
		{
			Keys: bson.D{
				bson.E{Key: "status", Value: 1},
				bson.E{Key: "node_addr", Value: 1},
			},
		},
	}

	_, err = database.SessionIndexesCreateMany(ctx, db, indexes)
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

	if err = db.Client().Ping(context.TODO(), nil); err != nil {
		log.Fatalln(err)
	}

	if err := createIndexes(context.TODO(), db); err != nil {
		log.Fatalln(err)
	}

	engine := gin.Default()
	engine.Use(cors.Default())

	router := engine.Group("/api/v1")

	blockapi.RegisterRoutes(router, db)
	depositapi.RegisterRoutes(router, db)
	nodeapi.RegisterRoutes(router, db)
	sessionapi.RegisterRoutes(router, db)
	statisticsapi.RegisterRoutes(router, db)
	subscriptionapi.RegisterRoutes(router, db)
	txapi.RegisterRoutes(router, db)

	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Fatalln(err)
	}
}
