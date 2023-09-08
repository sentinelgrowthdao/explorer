package plan

import (
	"strconv"

	"github.com/sentinel-official/explorer/types"
)

type EventAddPlan struct {
	ID uint64
}

func NewEventAddPlan(v *types.Event) (*EventAddPlan, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventAddPlan{
		ID: id,
	}, nil
}

func NewEventAddPlanFromEvents(v types.Events) (*EventAddPlan, error) {
	e, err := v.Get("sentinel.plan.v1.EventAddPlan")
	if err != nil {
		return nil, err
	}

	return NewEventAddPlan(e)
}
