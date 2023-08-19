package types

import (
	"time"

	explorerutils "github.com/sentinel-official/explorer/utils"
)

type AllocationEvent struct {
	ID            uint64 `json:"id,omitempty" bson:"id"`
	Address       string `json:"address,omitempty" bson:"address"`
	UtilisedBytes int64  `json:"utilised_bytes,omitempty" bson:"utilised_bytes"`
	GrantedBytes  int64  `json:"granted_bytes,omitempty" bson:"granted_bytes"`

	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (ae *AllocationEvent) String() string {
	return explorerutils.MustMarshalIndent(ae)
}
