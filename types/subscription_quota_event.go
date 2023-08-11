package types

import (
	"time"

	explorerutils "github.com/sentinel-official/explorer/utils"
)

type SubscriptionQuotaEvent struct {
	ID             uint64    `json:"id,omitempty" bson:"id"`
	Address        string    `json:"address,omitempty" bson:"address"`
	ConsumedBytes  int64     `json:"consumed_bytes,omitempty" bson:"consumed_bytes"`
	AllocatedBytes int64     `bson:"allocated_bytes,omitempty" bson:"allocated_bytes"`
	Height         int64     `json:"height,omitempty" bson:"height"`
	Timestamp      time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash         string    `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (sqe *SubscriptionQuotaEvent) String() string {
	return explorerutils.MustMarshalIndent(sqe)
}
