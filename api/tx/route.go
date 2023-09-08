package tx

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/blocks/:height/txs", HandlerGetTxs(db))

	router.GET("/txs", HandlerGetTxs(db))
	router.GET("/txs/:hash", HandlerGetTx(db))
}
