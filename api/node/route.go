package node

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/providers/:provider_address/nodes", HandlerGetNodes(db))
	router.GET("/nodes", HandlerGetNodes(db))
	router.GET("/nodes/:node_address", HandlerGetNode(db))
	router.GET("/nodes/:node_address/events", HandlerGetNodeEvents(db))
}
