package types

import (
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type SubscriptionQuota struct {
	ID             uint64 `json:"id,omitempty" bson:"id"`
	Address        string `json:"address,omitempty" bson:"address"`
	ConsumedBytes  int64  `json:"consumed_bytes,omitempty" bson:"consumed_bytes"`
	AllocatedBytes int64  `bson:"allocated_bytes,omitempty" bson:"allocated_bytes"`
}

func (sq *SubscriptionQuota) String() string {
	return explorerutils.MustMarshalIndent(sq)
}
