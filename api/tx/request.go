package tx

type RequestGetTxs struct {
	Height int64 `uri:"height"`

	FromHeight int64 `form:"from_height"`
	ToHeight   int64 `form:"to_height"`
	Skip       int64 `form:"skip"`
	Limit      int64 `form:"limit"`
}

type RequestGetTx struct {
	Hash string `uri:"hash"`
}
