package deposit

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/deposits", HandlerGetDeposits(db))
	router.GET("/deposits/:account_address", HandlerGetDeposit(db))
	router.GET("/deposits/:account_address/events", HandlerGetDepositEvents(db))
}
