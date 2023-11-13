package models

import (
	"github.com/sentinel-official/explorer/utils"
)

type SubscriptionAllocation struct {
	ID            uint64 `json:"id,omitempty" bson:"id"`
	AccAddr       string `json:"acc_addr,omitempty" bson:"acc_addr"`
	GrantedBytes  string `bson:"granted_bytes,omitempty" bson:"granted_bytes"`
	UtilisedBytes string `json:"utilised_bytes,omitempty" bson:"utilised_bytes"`
}

func (sq *SubscriptionAllocation) String() string {
	return utils.MustMarshalIndentToString(sq)
}
