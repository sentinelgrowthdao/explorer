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

type RequestGetQuotas struct {
	URI struct {
		ID uint64 `uri:"id"`
	}
	Query struct {
		Skip  int64 `form:"skip" binding:"gte=0"`
		Limit int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetQuotas(c *gin.Context) (req *RequestGetQuotas, err error) {
	req = &RequestGetQuotas{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetQuota struct {
	URI struct {
		AccAddr string `uri:"acc_addr"`
		ID      uint64 `uri:"id"`
	}
}

func NewRequestGetQuota(c *gin.Context) (req *RequestGetQuota, err error) {
	req = &RequestGetQuota{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetQuotaEvents struct {
	URI struct {
		AccAddr string `uri:"acc_addr"`
		ID      uint64 `uri:"id"`
	}
	Query struct {
		Skip  int64 `form:"skip" binding:"gte=0"`
		Limit int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetQuotaEvents(c *gin.Context) (req *RequestGetQuotaEvents, err error) {
	req = &RequestGetQuotaEvents{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}
