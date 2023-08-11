package deposit

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/utils"
)

type RequestGetDeposits struct {
	Skip   int64 `form:"skip,default=0" binding:"gte=0"`
	Limit  int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	Sort   bson.D
	SortBy string `form:"sort_by"`
}

func NewRequestGetDeposits(c *gin.Context) (req *RequestGetDeposits, err error) {
	req = &RequestGetDeposits{}
	if err = c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	allowed := []string{
		"coins.amount,coins.denom",
		"-coins.amount,coins.denom",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetDeposit struct {
	AccountAddress string `uri:"account_address"`
}

func NewRequestGetDeposit(c *gin.Context) (req *RequestGetDeposit, err error) {
	req = &RequestGetDeposit{}
	if err = c.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetDepositEvents struct {
	AccountAddress string `uri:"account_address"`

	Skip   int64 `form:"skip,default=0" binding:"gte=0"`
	Limit  int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	Sort   bson.D
	SortBy string `form:"sort_by"`
}

func NewRequestGetDepositEvents(c *gin.Context) (req *RequestGetDepositEvents, err error) {
	req = &RequestGetDepositEvents{}
	if err = c.ShouldBindUri(&req); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	allowed := []string{
		"height",
		"-height",
	}
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return nil, err
	}

	return req, nil
}
