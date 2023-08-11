package session

type RequestGetSessions struct {
	AccountAddress string `uri:"account_address"`
	NodeAddress    string `uri:"node_address"`
	ID             uint64 `uri:"id"`
	Status         string `form:"status"`
	Skip           int64  `form:"skip"`
	Limit          int64  `form:"limit"`
}

type RequestGetSession struct {
	ID uint64 `uri:"id"`
}

type RequestGetSessionEvents struct {
	ID    uint64 `uri:"id"`
	Skip  int64  `form:"skip"`
	Limit int64  `form:"limit"`
}
