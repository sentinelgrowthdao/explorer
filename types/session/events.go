package session

import (
	"strconv"

	"github.com/sentinel-official/explorer/types"
)

type EventStart struct {
	ID uint64
}

func NewEventStart(v *types.Event) (*EventStart, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventStart{
		ID: id,
	}, nil
}

func NewEventStartFromEvents(v types.Events, skip int) (int, *EventStart, error) {
	i, e, err := v.Get("sentinel.session.v2.EventStart", skip)
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventStart(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
}

type EventUpdateStatus struct {
	ID     uint64
	Status string
}

func NewEventUpdateStatus(v *types.Event) (*EventUpdateStatus, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventUpdateStatus{
		ID:     id,
		Status: v.Attributes["status"],
	}, nil
}
