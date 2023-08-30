package plan

import (
	"strconv"

	"github.com/sentinel-official/explorer/types"
)

type EventAdd struct {
	ID uint64
}

func NewEventAdd(v *types.Event) (*EventAdd, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventAdd{
		ID: id,
	}, nil
}

func NewEventAddFromEvents(v types.Events) (*EventAdd, error) {
	e, err := v.Get("sentinel.plan.v1.EventAdd")
	if err != nil {
		return nil, err
	}

	return NewEventAdd(e)
}
