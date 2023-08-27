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

	buf, err := json.Marshal(v.Attributes["price"])
	if err != nil {
		return nil, err
	}

	var price sdk.Coin
	if err := json.Unmarshal(buf, &price); err != nil {
		return nil, err
	}

	buf, err = json.Marshal(v.Attributes["deposit"])
	if err != nil {
		return nil, err
	}

	var deposit sdk.Coin
	if err := json.Unmarshal(buf, &deposit); err != nil {
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
	Denom   string
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

	buf, err := json.Marshal(v.Attributes["price"])
	if err != nil {
		return nil, err
	}

	var payment sdk.Coin
	if err := json.Unmarshal(buf, &payment); err != nil {
		return nil, err
	}

	return &EventSubscribeToPlan{
		Owner:   v.Attributes["owner"],
		Denom:   v.Attributes["denom"],
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
