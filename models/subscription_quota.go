package models

import (
	"math/big"

	"github.com/sentinel-official/explorer/utils"
)

type SubscriptionQuota struct {
	ID        uint64   `json:"id,omitempty" bson:"id"`
	Address   string   `json:"address,omitempty" bson:"address"`
	Allocated *big.Int `bson:"allocated,omitempty" bson:"allocated"`
	Consumed  *big.Int `json:"consumed,omitempty" bson:"consumed"`
}

func (sq *SubscriptionQuota) String() string {
	return utils.MustMarshalIndentToString(sq)
}
