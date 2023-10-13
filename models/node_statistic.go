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

	SessionEndCount        int64 `json:"session_end_count" bson:"session_end_count"`
	SessionStartCount      int64 `json:"session_start_count" bson:"session_start_count"`
	SubscriptionEndCount   int64 `json:"subscription_end_count" bson:"subscription_end_count"`
	SubscriptionStartCount int64 `json:"subscription_start_count" bson:"subscription_start_count"`

	SessionBandwidth *types.Bandwidth `json:"session_bandwidth" bson:"session_bandwidth"`
	SessionDuration  int64            `json:"session_duration" bson:"session_duration"`

	SubscriptionBytes            string      `json:"subscription_bytes" bson:"subscription_bytes"`
	SubscriptionHours            int64       `json:"subscription_hours" bson:"subscription_hours"`
	SubscriptionEarningsForBytes types.Coins `json:"subscription_earnings_for_bytes" bson:"subscription_earnings_for_bytes"`
	SubscriptionEarningsForHours types.Coins `json:"subscription_earnings_for_hours" bson:"subscription_earnings_for_hours"`
}

func (ns *NodeStatistic) String() string {
	return utils.MustMarshalIndentToString(ns)
}
