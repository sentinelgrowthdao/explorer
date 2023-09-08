package tx

import (
	"github.com/gin-gonic/gin"
)

type RequestGetTxs struct {
	URI struct {
		Height int64 `uri:"height"`
	}
	Query struct {
		FromHeight int64 `form:"from_height"`
		ToHeight   int64 `form:"to_height,default=1000000000"`
		Skip       int64 `form:"skip" binding:"gte=0"`
		Limit      int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetTxs(c *gin.Context) (req *RequestGetTxs, err error) {
	req = &RequestGetTxs{}
	if err := c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetTx struct {
	URI struct {
		Hash string `uri:"hash"`
	}
}

func NewRequestGetTx(c *gin.Context) (req *RequestGetTx, err error) {
	req = &RequestGetTx{}
	if err := c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}
