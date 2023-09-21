package models

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/utils"
)

type SubscriptionQuota struct {
	ID        uint64  `json:"id,omitempty" bson:"id"`
	Address   string  `json:"address,omitempty" bson:"address"`
	Allocated sdk.Int `bson:"allocated,omitempty" bson:"allocated"`
	Consumed  sdk.Int `json:"consumed,omitempty" bson:"consumed"`
}

func (sq *SubscriptionQuota) String() string {
	return utils.MustMarshalIndentToString(sq)
}
