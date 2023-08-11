package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	blockapi "github.com/sentinel-official/explorer/api/block"
	depositapi "github.com/sentinel-official/explorer/api/deposit"
	nodeapi "github.com/sentinel-official/explorer/api/node"
	sessionapi "github.com/sentinel-official/explorer/api/session"
	statisticsapi "github.com/sentinel-official/explorer/api/statistics"
	subscriptionapi "github.com/sentinel-official/explorer/api/subscription"
	txapi "github.com/sentinel-official/explorer/api/tx"
	"github.com/sentinel-official/explorer/database"
)

const (
	appName = "00-api_server"
)

var (
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
)

func init() {
	log.SetOutput(io.Discard)

	flag.StringVar(&dbAddress, "db-address", "mongodb://127.0.0.1:27017", "")
	flag.StringVar(&dbName, "db-name", "sentinelhub-2", "")
	flag.StringVar(&dbUsername, "db-username", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.Parse()
}

func main() {
	db, err := database.PrepareDatabase(context.TODO(), appName, dbAddress, dbUsername, dbPassword, dbName)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Client().Ping(context.TODO(), nil); err != nil {
		log.Fatalln(err)
	}

	engine := gin.Default()
	engine.Use(cors.Default())

	router := engine.Group("/api/v1")

	blockapi.RegisterRoutes(router, db)
	depositapi.RegisterRoutes(router, db)
	nodeapi.RegisterRoutes(router, db)
	sessionapi.RegisterRoutes(router, db)
	subscriptionapi.RegisterRoutes(router, db)
	statisticsapi.RegisterRoutes(router, db)
	txapi.RegisterRoutes(router, db)

	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Fatalln(err)
	}
}
