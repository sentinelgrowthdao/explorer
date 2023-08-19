package types

import (
	"time"

	commontypes "github.com/sentinel-official/explorer/types/common"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Subscription struct {
	ID         uint64    `json:"id,omitempty" bson:"id"`
	Address    string    `json:"address,omitempty" bson:"address"`
	InactiveAt time.Time `json:"inactive_at,omitempty" bson:"inactive_at"`

	NodeAddress string            `json:"node_address,omitempty" bson:"node_address"`
	Gigabytes   int64             `json:"gigabytes,omitempty" bson:"gigabytes"`
	Hours       int64             `json:"hours,omitempty" bson:"hours"`
	Deposit     *commontypes.Coin `json:"deposit,omitempty" bson:"deposit"`

	PlanID        uint64            `json:"plan_id,omitempty" bson:"plan_id"`
	Denom         string            `json:"denom,omitempty" bson:"denom"`
	StakingReward *commontypes.Coin `json:"staking_reward,omitempty" bson:"staking_reward"`
	Payment       *commontypes.Coin `json:"payment,omitempty" bson:"payment"`

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
	return explorerutils.MustMarshalIndent(s)
}
