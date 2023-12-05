package deposit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/types"
)

type EventAdd struct {
	Address string
	Coins   types.Coins
}

func NewEventAdd(v *types.Event) (*EventAdd, error) {
	coins, err := sdk.ParseCoinsNormalized(v.Attributes["coins"])
	if err != nil {
		return nil, err
	}

	return &EventAdd{
		Address: v.Attributes["address"],
		Coins:   types.NewCoins(coins),
	}, nil
}

func NewEventAddFromEvents(v types.Events, skip int) (int, *EventAdd, error) {
	i, e, err := v.Get("sentinel.deposit.v1.EventAdd", skip)
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventAdd(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
}

type EventSubtract struct {
	Address string
	Coins   types.Coins
}

func NewEventSubtract(v *types.Event) (*EventSubtract, error) {
	coins, err := sdk.ParseCoinsNormalized(v.Attributes["coins"])
	if err != nil {
		return nil, err
	}

	return &EventSubtract{
		Address: v.Attributes["address"],
		Coins:   types.NewCoins(coins),
	}, nil
}
