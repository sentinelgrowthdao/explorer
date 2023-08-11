package subscription

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/accounts/:account_address/subscriptions", HandlerGetSubscriptions(db))
	router.GET("/nodes/:node_address/subscriptions", HandlerGetSubscriptions(db))
	router.GET("/plans/:id/subscriptions", HandlerGetSubscriptions(db))
	router.GET("/subscriptions", HandlerGetSubscriptions(db))
	router.GET("/subscriptions/:id", HandlerGetSubscription(db))
	router.GET("/subscriptions/:id/quotas", HandlerGetSubscriptionQuotas(db))
	router.GET("/subscriptions/:id/quotas/:account_address", HandlerGetSubscriptionQuota(db))
	router.GET("/subscriptions/:id/quotas/:account_address/events", HandlerGetSubscriptionQuotaEvents(db))
}
