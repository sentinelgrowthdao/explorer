package types

const (
	EventTypeDepositAdd = 1 + iota
	EventTypeDepositSubtract
	EventTypeNodeUpdateDetails
	EventTypeNodeUpdateStatus
	EventTypePlanUpdateStatus
	EventTypePlanLinkNode
	EventTypePlanUnlinkNode
	EventTypeProviderUpdateDetails
	EventTypeSessionUpdateDetails
	EventTypeSessionUpdateStatus
	EventTypeSubscriptionUpdateDetails
	EventTypeSubscriptionUpdateStatus
	EventTypeSubscriptionAllocationUpdateDetails
)
