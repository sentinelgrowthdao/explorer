package deposit

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/types"
)

type EventAdd struct {
	Address  string
	Previous types.Coins
	Coins    types.Coins
	Current  types.Coins
}

func NewEventAdd(v *types.Event) (*EventAdd, error) {
	buf, err := json.Marshal(v.Attributes["previous"])
	if err != nil {
		return nil, err
	}

	var previous sdk.Coins
	if err := json.Unmarshal(buf, &previous); err != nil {
		return nil, err
	}

	buf, err = json.Marshal(v.Attributes["coins"])
	if err != nil {
		return nil, err
	}

	var coins sdk.Coins
	if err := json.Unmarshal(buf, &coins); err != nil {
		return nil, err
	}

	return &EventAdd{
		Address:  v.Attributes["address"],
		Previous: types.NewCoins(previous),
		Coins:    types.NewCoins(coins),
		Current:  types.NewCoins(previous.Add(coins...)),
	}, nil
}

func NewEventAddFromEvents(v types.Events) (*EventAdd, error) {
	e, err := v.Get("sentinel.deposit.v1.EventAdd")
	if err != nil {
		return nil, err
	}

	return NewEventAdd(e)
}

type EventSubtract struct {
	Address  string
	Previous types.Coins
	Coins    types.Coins
	Current  types.Coins
}

func NewEventSubtract(v *types.Event) (*EventSubtract, error) {
	buf, err := json.Marshal(v.Attributes["previous"])
	if err != nil {
		return nil, err
	}

	var previous sdk.Coins
	if err := json.Unmarshal(buf, &previous); err != nil {
		return nil, err
	}

	buf, err = json.Marshal(v.Attributes["coins"])
	if err != nil {
		return nil, err
	}

	var coins sdk.Coins
	if err := json.Unmarshal(buf, &coins); err != nil {
		return nil, err
	}

	return &EventSubtract{
		Address:  v.Attributes["address"],
		Previous: types.NewCoins(previous),
		Coins:    types.NewCoins(coins),
		Current:  types.NewCoins(previous.Sub(coins)),
	}, nil
}

func NewEventSubtractFromEvents(v types.Events) (*EventSubtract, error) {
	e, err := v.Get("sentinel.deposit.v1.EventSubtract")
	if err != nil {
		return nil, err
	}

	return NewEventSubtract(e)
}
