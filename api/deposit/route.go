package deposit

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, db *mongo.Database) {
	router.GET("/deposits", HandlerGetDeposits(db))
	router.GET("/deposits/:acc_addr", HandlerGetDeposit(db))
	router.GET("/deposits/:acc_addr/events", HandlerGetDepositEvents(db))
}
