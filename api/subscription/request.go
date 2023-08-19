package subscription

type RequestGetSubscriptions struct {
	AccountAddress string `uri:"account_address"`
	ID             uint64 `uri:"id"`
	NodeAddress    string `uri:"node_address"`
	Status         string `form:"status"`
	Skip           int64  `form:"skip"`
	Limit          int64  `form:"limit"`
}

type RequestGetSubscription struct {
	ID uint64 `uri:"id"`
}

type RequestGetAllocations struct {
	ID    uint64 `uri:"id"`
	Skip  int64  `form:"skip"`
	Limit int64  `form:"limit"`
}

type RequestGetAllocation struct {
	ID             uint64 `uri:"id"`
	AccountAddress string `uri:"account_address"`
	Skip           int64  `form:"skip"`
	Limit          int64  `form:"limit"`
}

type RequestGetAllocationEvents struct {
	ID             uint64 `uri:"id"`
	AccountAddress string `uri:"account_address"`
	Skip           int64  `form:"skip"`
	Limit          int64  `form:"limit"`
}
