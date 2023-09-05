package models

import (
	"github.com/sentinel-official/explorer/utils"
)

type SubscriptionQuota struct {
	ID        uint64 `json:"id,omitempty" bson:"id"`
	Address   string `json:"address,omitempty" bson:"address"`
	Allocated int64  `bson:"allocated,omitempty" bson:"allocated"`
	Consumed  int64  `json:"consumed,omitempty" bson:"consumed"`
}

func (sq *SubscriptionQuota) String() string {
	return utils.MustMarshalIndentToString(sq)
}
