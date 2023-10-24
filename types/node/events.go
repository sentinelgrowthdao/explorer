package node

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/types"
)

type EventUpdateDetails struct {
	Address        string
	GigabytePrices types.Coins
	HourlyPrices   types.Coins
	RemoteURL      string
}

func NewEventUpdateDetails(v *types.Event) (*EventUpdateDetails, error) {
	gigabytePrices, err := sdk.ParseCoinsNormalized(v.Attributes["gigabyte_prices"])
	if err != nil {
		return nil, err
	}

	hourlyPrices, err := sdk.ParseCoinsNormalized(v.Attributes["hourly_prices"])
	if err != nil {
		return nil, err
	}

	return &EventUpdateDetails{
		Address:        v.Attributes["address"],
		GigabytePrices: types.NewCoins(gigabytePrices),
		HourlyPrices:   types.NewCoins(hourlyPrices),
		RemoteURL:      v.Attributes["remote_url"],
	}, nil
}

type EventUpdateStatus struct {
	Address string
	Status  string
}

func NewEventUpdateStatus(v *types.Event) (*EventUpdateStatus, error) {
	return &EventUpdateStatus{
		Address: v.Attributes["address"],
		Status:  v.Attributes["status"],
	}, nil
}

type EventCreateSubscription struct {
	ID uint64
}

func NewEventCreateSubscription(v *types.Event) (*EventCreateSubscription, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventCreateSubscription{
		ID: id,
	}, nil
}

func NewEventCreateSubscriptionFromEvents(v types.Events) (int, *EventCreateSubscription, error) {
	i, e, err := v.Get("sentinel.node.v2.EventCreateSubscription")
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventCreateSubscription(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
}
