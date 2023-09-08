package block

import (
	"github.com/gin-gonic/gin"
)

type RequestGetBlocks struct {
	Query struct {
		FromHeight int64 `form:"from_height"`
		ToHeight   int64 `form:"to_height,default=1000000000"`
		Skip       int64 `form:"skip" binding:"gte=0"`
		Limit      int64 `form:"limit,default=25" binding:"gte=0,lte=100"`
	}
}

func NewRequestGetBlocks(c *gin.Context) (req *RequestGetBlocks, err error) {
	req = &RequestGetBlocks{}
	if err = c.ShouldBindQuery(&req.Query); err != nil {
		return nil, err
	}

	return req, nil
}

type RequestGetBlock struct {
	URI struct {
		Height int64 `uri:"height"`
	}
}

func NewRequestGetBlock(c *gin.Context) (req *RequestGetBlock, err error) {
	req = &RequestGetBlock{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}
