package plan

import (
	"strconv"

	"github.com/sentinel-official/explorer/types"
)

type EventCreate struct {
	ID uint64
}

func NewEventCreate(v *types.Event) (*EventCreate, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventCreate{
		ID: id,
	}, nil
}

func NewEventCreateFromEvents(v types.Events) (int, *EventCreate, error) {
	i, e, err := v.Get("sentinel.plan.v2.EventCreate")
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventCreate(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
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
	i, e, err := v.Get("sentinel.plan.v2.EventCreateSubscription")
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventCreateSubscription(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
}
