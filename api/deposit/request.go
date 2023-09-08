package deposit

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/utils"
)

type RequestGetDeposits struct {
	Sort bson.D

	Query struct {
		Sort  string `form:"sort"`
		Skip  int64  `form:"skip,default=0" binding:"gte=0"`
		Limit int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetDeposits(c *gin.Context) (req *RequestGetDeposits, err error) {
	req = &RequestGetDeposits{}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	allowed := []string{
		"coins.amount,coins.denom",
		"-coins.amount,coins.denom",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetDeposit struct {
	URI struct {
		AccAddr string `uri:"acc_addr"`
	}
}

func NewRequestGetDeposit(c *gin.Context) (req *RequestGetDeposit, err error) {
	req = &RequestGetDeposit{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetDepositEvents struct {
	Sort bson.D

	URI struct {
		AccAddr string `uri:"acc_addr"`
	}
	Query struct {
		Sort  string `form:"sort"`
		Skip  int64  `form:"skip,default=0" binding:"gte=0"`
		Limit int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetDepositEvents(c *gin.Context) (req *RequestGetDepositEvents, err error) {
	req = &RequestGetDepositEvents{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	allowed := []string{
		"height",
		"-height",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return nil, err
	}

	return req, nil
}
