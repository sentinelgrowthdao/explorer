package types

import (
	"time"

	explorerutils "github.com/sentinel-official/explorer/utils"
)

type NodeEvent struct {
	Address   string    `json:"address,omitempty" bson:"address"`
	Status    string    `json:"status,omitempty" bson:"status"`
	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (ne *NodeEvent) String() string {
	return explorerutils.MustMarshalIndent(ne)
}
