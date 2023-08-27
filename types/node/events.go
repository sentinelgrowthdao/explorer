package node

import (
	"github.com/sentinel-official/explorer/types"
)

type EventSetNodeStatus struct {
	Address string
	Status  string
}

func NewEventSetNodeStatus(v *types.Event) (*EventSetNodeStatus, error) {
	return &EventSetNodeStatus{
		Address: v.Attributes["address"],
		Status:  v.Attributes["status"],
	}, nil
}
