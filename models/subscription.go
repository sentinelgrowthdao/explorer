package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type Subscription struct {
	ID    uint64 `json:"id,omitempty" bson:"id"`
	Owner string `json:"owner,omitempty" bson:"owner"`
	Free  string `json:"free,omitempty" bson:"free"`

	Node    string      `json:"node,omitempty" bson:"node"`
	Price   *types.Coin `json:"price,omitempty" bson:"price"`
	Deposit *types.Coin `json:"deposit,omitempty" bson:"deposit"`
	Refund  *types.Coin `json:"refund,omitempty" bson:"refund"`

	Plan          uint64      `json:"plan,omitempty" bson:"plan"`
	Denom         string      `json:"denom,omitempty" bson:"denom"`
	Expiry        time.Time   `json:"expiry,omitempty" bson:"expiry"`
	Payment       *types.Coin `json:"payment,omitempty" bson:"payment"`
	StakingReward *types.Coin `json:"staking_reward,omitempty" bson:"staking_reward"`

	StartHeight    int64     `json:"start_height,omitempty" bson:"start_height"`
	StartTimestamp time.Time `json:"start_timestamp,omitempty" bson:"start_timestamp"`
	StartTxHash    string    `json:"start_tx_hash,omitempty" bson:"start_tx_hash"`
	EndHeight      int64     `json:"end_height,omitempty" bson:"end_height"`
	EndTimestamp   time.Time `json:"end_timestamp,omitempty" bson:"end_timestamp"`
	EndTxHash      string    `json:"end_tx_hash,omitempty" bson:"end_tx_hash"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`
}

func (s *Subscription) String() string {
	return utils.MustMarshalIndentToString(s)
}
