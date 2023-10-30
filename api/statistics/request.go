package statistics

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type RequestGetStatistics struct {
	Sort bson.D

	Query struct {
		FromTimestamp time.Time `form:"from_timestamp"`
		Limit         int64     `form:"limit,default=30" binding:"gte=0,lte=100"`
		Method        string    `form:"method" binding:"required"`
		Skip          int64     `form:"skip,default=0" binding:"gte=0"`
		Sort          string    `form:"sort"`
		Status        string    `form:"status" binding:"omitempty,oneof=active inactive inactive_pending"`
		Timeframe     string    `form:"timeframe,default=day" binding:"oneof=day week month year"`
		ToTimestamp   time.Time `form:"to_timestamp,default=3001-01-01T00:00:00.0Z" binding:"gtfield=FromTimestamp"`
	}
}

func NewRequestGetStatistics(c *gin.Context) (req *RequestGetStatistics, err error) {
	req = &RequestGetStatistics{}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	validatorFunc, ok := validators[req.Query.Method]
	if !ok {
		return req, nil
	}
	if validatorFunc == nil {
		return req, nil
	}

	if err := validatorFunc(req); err != nil {
		return nil, err
	}

	return req, nil
}
