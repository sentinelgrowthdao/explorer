package node

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/utils"
)

type RequestGetNodes struct {
	ProviderAddress string `uri:"provider_address"`

	Status string `form:"status" binding:"omitempty,oneof=STATUS_ACTIVE STATUS_INACTIVE"`
	Skip   int64  `form:"skip,default=0" binding:"gte=0"`
	Limit  int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	Sort   bson.D
	SortBy string `form:"sort_by"`
}

func NewRequestGetNodes(c *gin.Context) (req *RequestGetNodes, err error) {
	req = &RequestGetNodes{}
	if err = c.ShouldBindUri(&req); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req); err != nil {
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
	if req.Sort, err = utils.ParseQuerySortBy(allowed, req.SortBy); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetNode struct {
	NodeAddress string `uri:"node_address"`
}

func NewRequestGetNode(c *gin.Context) (req *RequestGetNode, err error) {
	req = &RequestGetNode{}
	if err = c.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetNodeEvents struct {
	NodeAddress string `uri:"node_address"`

	Status string `form:"status" binding:"omitempty,oneof=STATUS_ACTIVE STATUS_INACTIVE"`
	Skip   int64  `form:"skip,default=0" binding:"gte=0"`
	Limit  int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	Sort   bson.D
	SortBy string `form:"sort_by"`
}

func NewRequestGetNodeEvents(c *gin.Context) (req *RequestGetNodeEvents, err error) {
	req = &RequestGetNodeEvents{}
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
