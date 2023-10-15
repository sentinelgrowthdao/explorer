package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type NodeStatistic struct {
	Address   string    `json:"address" bson:"address"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`

	GasFeeSpent types.Coins `json:"gas_fee_spent" bson:"gas_fee_spent"`

	SessionBandwidth  *types.Bandwidth `json:"session_bandwidth" bson:"session_bandwidth"`
	SessionDuration   int64            `json:"session_duration" bson:"session_duration"`
	SessionEndCount   int64            `json:"session_end_count" bson:"session_end_count"`
	SessionStartCount int64            `json:"session_start_count" bson:"session_start_count"`

	SubscriptionBytes      string `json:"subscription_bytes" bson:"subscription_bytes"`
	SubscriptionEndCount   int64  `json:"subscription_end_count" bson:"subscription_end_count"`
	SubscriptionHours      int64  `json:"subscription_hours" bson:"subscription_hours"`
	SubscriptionStartCount int64  `json:"subscription_start_count" bson:"subscription_start_count"`

	EarningsForBytes types.Coins `json:"earnings_for_bytes" bson:"earnings_for_bytes"`
	EarningsForHours types.Coins `json:"earnings_for_hours" bson:"earnings_for_hours"`
}

func NewNodeStatistic() *NodeStatistic {
	return &NodeStatistic{
		GasFeeSpent:      types.NewCoins(nil),
		SessionBandwidth: types.NewBandwidth(nil),
		EarningsForBytes: types.NewCoins(nil),
		EarningsForHours: types.NewCoins(nil),
	}
}

func (ns *NodeStatistic) String() string {
	return utils.MustMarshalIndentToString(ns)
}
