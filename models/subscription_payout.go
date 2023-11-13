package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type SubscriptionPayout struct {
	ID            uint64      `json:"id,omitempty" bson:"id"`
	AccAddr       string      `json:"acc_addr,omitempty" bson:"acc_addr"`
	NodeAddr      string      `json:"node_addr,omitempty" bson:"node_addr"`
	Payment       *types.Coin `json:"payment,omitempty" bson:"payment"`
	StakingReward *types.Coin `json:"staking_reward,omitempty" bson:"staking_reward"`

	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (sq *SubscriptionPayout) String() string {
	return utils.MustMarshalIndentToString(sq)
}
