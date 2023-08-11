package block

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/blocks", HandlerGetBlocks(db))
	router.GET("/blocks/:height", HandlerGetBlock(db))
}
