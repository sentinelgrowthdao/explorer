package subscription

import (
	"encoding/json"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/types"
)

type EventSubscribeToNode struct {
	Owner   string
	Node    string
	ID      uint64
	Price   *types.Coin
	Deposit *types.Coin
}

func NewEventSubscribeToNode(v *types.Event) (*EventSubscribeToNode, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	var price sdk.Coin
	if err := json.Unmarshal([]byte(v.Attributes["price"]), &price); err != nil {
		return nil, err
	}

	var deposit sdk.Coin
	if err := json.Unmarshal([]byte(v.Attributes["deposit"]), &deposit); err != nil {
		return nil, err
	}

	return &EventSubscribeToNode{
		Owner:   v.Attributes["owner"],
		Node:    v.Attributes["node"],
		ID:      id,
		Price:   types.NewCoin(&price),
		Deposit: types.NewCoin(&deposit),
	}, nil
}

func NewEventSubscribeToNodeFromEvents(v types.Events) (*EventSubscribeToNode, error) {
	e, err := v.Get("sentinel.subscription.v1.EventSubscribeToNode")
	if err != nil {
		return nil, err
	}

	return NewEventSubscribeToNode(e)
}

type EventSubscribeToPlan struct {
	Owner   string
	ID      uint64
	Plan    uint64
	Expiry  time.Time
	Payment *types.Coin
}

func NewEventSubscribeToPlan(v *types.Event) (*EventSubscribeToPlan, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	plan, err := strconv.ParseUint(v.Attributes["plan"], 10, 64)
	if err != nil {
		return nil, err
	}

	expiry, err := time.Parse(time.RFC3339Nano, v.Attributes["expiry"])
	if err != nil {
		return nil, err
	}

	var payment sdk.Coin
	if err := json.Unmarshal([]byte(v.Attributes["price"]), &payment); err != nil {
		return nil, err
	}

	return &EventSubscribeToPlan{
		Owner:   v.Attributes["owner"],
		ID:      id,
		Plan:    plan,
		Expiry:  expiry,
		Payment: types.NewCoin(&payment),
	}, nil
}

func NewEventSubscribeToPlanFromEvents(v types.Events) (*EventSubscribeToPlan, error) {
	e, err := v.Get("sentinel.subscription.v1.EventSubscribeToPlan")
	if err != nil {
		return nil, err
	}

	return NewEventSubscribeToPlan(e)
}

type EventAddQuota struct {
	ID        uint64
	Address   string
	Consumed  int64
	Allocated int64
	Free      int64
}

func NewEventAddQuota(v *types.Event) (*EventAddQuota, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	consumed, err := strconv.ParseInt(v.Attributes["consumed"], 10, 64)
	if err != nil {
		return nil, err
	}

	allocated, err := strconv.ParseInt(v.Attributes["allocated"], 10, 64)
	if err != nil {
		return nil, err
	}

	free, err := strconv.ParseInt(v.Attributes["free"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventAddQuota{
		ID:        id,
		Address:   v.Attributes["address"],
		Consumed:  consumed,
		Allocated: allocated,
		Free:      free,
	}, nil
}

func NewEventAddQuotaFromEvents(v types.Events) (*EventAddQuota, error) {
	e, err := v.Get("sentinel.subscription.v1.EventAddQuota")
	if err != nil {
		return nil, err
	}

	return NewEventAddQuota(e)
}

type EventUpdateQuota struct {
	ID        uint64
	Address   string
	Consumed  int64
	Allocated int64
	Free      int64
}

func NewEventUpdateQuota(v *types.Event) (*EventUpdateQuota, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	consumed, err := strconv.ParseInt(v.Attributes["consumed"], 10, 64)
	if err != nil {
		return nil, err
	}

	allocated, err := strconv.ParseInt(v.Attributes["allocated"], 10, 64)
	if err != nil {
		return nil, err
	}

	free, err := strconv.ParseInt(v.Attributes["free"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventUpdateQuota{
		ID:        id,
		Address:   v.Attributes["address"],
		Consumed:  consumed,
		Allocated: allocated,
		Free:      free,
	}, nil
}

func NewEventUpdateQuotaFromEvents(v types.Events) (*EventUpdateQuota, error) {
	e, err := v.Get("sentinel.subscription.v1.EventUpdateQuota")
	if err != nil {
		return nil, err
	}

	return NewEventUpdateQuota(e)
}

type EventCancelSubscription struct {
	ID     uint64
	Status string
}

func NewEventCancelSubscription(v *types.Event) (*EventCancelSubscription, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventCancelSubscription{
		ID:     id,
		Status: v.Attributes["status"],
	}, nil
}

func NewEventCancelSubscriptionFromEvents(v types.Events) (*EventCancelSubscription, error) {
	e, err := v.Get("sentinel.subscription.v1.EventCancelSubscription")
	if err != nil {
		return nil, err
	}

	return NewEventCancelSubscription(e)
}

type EventRefund struct {
	ID     uint64
	Refund *types.Coin
}

func NewEventRefund(v *types.Event) (*EventRefund, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	var refund sdk.Coin
	if err := json.Unmarshal([]byte(v.Attributes["refund"]), &refund); err != nil {
		return nil, err
	}

	return &EventRefund{
		ID:     id,
		Refund: types.NewCoin(&refund),
	}, nil
}

func NewEventRefundFromEvents(v types.Events) (*EventRefund, error) {
	e, err := v.Get("sentinel.subscription.v1.EventRefund")
	if err != nil {
		return nil, err
	}

	return NewEventRefund(e)
}
