package block

type RequestGetBlocks struct {
	FromHeight int64 `form:"from_height"`
	ToHeight   int64 `form:"to_height"`
	Skip       int64 `form:"skip"`
	Limit      int64 `form:"limit"`
}

type RequestGetBlock struct {
	Height int64 `uri:"height"`
}
