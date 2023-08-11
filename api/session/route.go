package session

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/accounts/:account_address/sessions", HandlerGetSessions(db))
	router.GET("/nodes/:node_address/sessions", HandlerGetSessions(db))
	router.GET("/sessions", HandlerGetSessions(db))
	router.GET("/sessions/:id", HandlerGetSession(db))
	router.GET("/sessions/:id/events", HandlerGetSessionEvents(db))
	router.GET("/subscriptions/:id/sessions", HandlerGetSessions(db))
}
