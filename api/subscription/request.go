package subscription

import (
	"github.com/gin-gonic/gin"
)

type RequestGetSubscriptions struct {
	URI struct {
		AccAddr  string `uri:"acc_addr"`
		ID       uint64 `uri:"id"`
		NodeAddr string `uri:"node_addr"`
	}
	Query struct {
		Status string `form:"status" binding:"omitempty,oneof=active inactive_pending inactive"`
		Skip   int64  `form:"skip" binding:"gte=0"`
		Limit  int64  `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetSubscriptions(c *gin.Context) (req *RequestGetSubscriptions, err error) {
	req = &RequestGetSubscriptions{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetSubscription struct {
	URI struct {
		ID uint64 `uri:"id"`
	}
}

func NewRequestGetSubscription(c *gin.Context) (req *RequestGetSubscription, err error) {
	req = &RequestGetSubscription{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetSubscriptionEvents struct {
	URI struct {
		ID uint64 `uri:"id"`
	}
	Query struct {
		Skip  int64 `form:"skip" binding:"gte=0"`
		Limit int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetSubscriptionEvents(c *gin.Context) (req *RequestGetSubscriptionEvents, err error) {
	req = &RequestGetSubscriptionEvents{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetAllocations struct {
	URI struct {
		ID uint64 `uri:"id"`
	}
	Query struct {
		Skip  int64 `form:"skip" binding:"gte=0"`
		Limit int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetAllocations(c *gin.Context) (req *RequestGetAllocations, err error) {
	req = &RequestGetAllocations{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetAllocation struct {
	URI struct {
		AccAddr string `uri:"acc_addr"`
		ID      uint64 `uri:"id"`
	}
}

func NewRequestGetAllocation(c *gin.Context) (req *RequestGetAllocation, err error) {
	req = &RequestGetAllocation{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetAllocationEvents struct {
	URI struct {
		AccAddr string `uri:"acc_addr"`
		ID      uint64 `uri:"id"`
	}
	Query struct {
		Skip  int64 `form:"skip" binding:"gte=0"`
		Limit int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetAllocationEvents(c *gin.Context) (req *RequestGetAllocationEvents, err error) {
	req = &RequestGetAllocationEvents{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}
