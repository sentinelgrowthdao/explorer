package node

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/nodes", HandlerGetNodes(db))
	router.GET("/nodes/:node_addr", HandlerGetNode(db))
	router.GET("/nodes/:node_addr/events", HandlerGetNodeEvents(db))
}
