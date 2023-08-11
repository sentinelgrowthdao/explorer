package types

import (
	"time"

	explorerutils "github.com/sentinel-official/explorer/utils"
)

type NodeReachEvent struct {
	Address      string    `json:"address,omitempty" bson:"address"`
	ErrorMessage string    `json:"error_message,omitempty" bson:"error_message"`
	Timestamp    time.Time `json:"timestamp,omitempty" bson:"timestamp"`
}

func (nve *NodeReachEvent) String() string {
	return explorerutils.MustMarshalIndent(nve)
}
