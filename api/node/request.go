package node

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/utils"
)

type RequestGetNodes struct {
	Sort bson.D

	Query struct {
		Status string `form:"status" binding:"omitempty,oneof=active inactive"`
		Sort   string `form:"sort"`
		Skip   int64  `form:"skip" binding:"gte=0"`
		Limit  int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetNodes(c *gin.Context) (req *RequestGetNodes, err error) {
	req = &RequestGetNodes{}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	allowed := []string{
		"bandwidth.download,bandwidth.upload,-peers",
		"-bandwidth.download,-bandwidth.upload,peers",
		"peers",
		"-peers",
		"join_height",
		"-join_height",
		"location.country,location.city",
		"-location.country,-location.city",
		"moniker",
		"-moniker",
		"price.amount,price.denom",
		"-price.amount,price.denom",
		"version",
		"-version",
	}
	if req.Sort, err = utils.ParseQuerySort(allowed, req.Query.Sort); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetNode struct {
	URI struct {
		NodeAddr string `uri:"node_addr"`
	}
}

func NewRequestGetNode(c *gin.Context) (req *RequestGetNode, err error) {
	req = &RequestGetNode{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetNodeEvents struct {
	Sort bson.D

	URI struct {
		NodeAddr string `uri:"node_addr"`
	}
	Query struct {
		Sort  string `form:"sort"`
		Skip  int64  `form:"skip" binding:"gte=0"`
		Limit int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetNodeEvents(c *gin.Context) (req *RequestGetNodeEvents, err error) {
	req = &RequestGetNodeEvents{}
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
