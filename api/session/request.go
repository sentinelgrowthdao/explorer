package session

import (
	"github.com/gin-gonic/gin"
)

type RequestGetSessions struct {
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

func NewRequestGetSessions(c *gin.Context) (req *RequestGetSessions, err error) {
	req = &RequestGetSessions{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetSession struct {
	URI struct {
		ID uint64 `uri:"id"`
	}
}

func NewRequestGetSession(c *gin.Context) (req *RequestGetSession, err error) {
	req = &RequestGetSession{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetSessionEvents struct {
	URI struct {
		ID uint64 `uri:"id"`
	}
	Query struct {
		Skip  int64 `form:"skip" binding:"gte=0"`
		Limit int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetSessionEvents(c *gin.Context) (req *RequestGetSessionEvents, err error) {
	req = &RequestGetSessionEvents{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}
