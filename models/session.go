package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type Session struct {
	ID           uint64           `json:"id,omitempty" bson:"id"`
	Subscription uint64           `json:"subscription,omitempty" bson:"subscription"`
	Address      string           `json:"address,omitempty" bson:"address"`
	Node         string           `json:"node,omitempty" bson:"node"`
	Duration     int64            `json:"duration,omitempty" bson:"duration"`
	Bandwidth    *types.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth"`

	StartHeight    int64     `json:"start_height,omitempty" bson:"start_height"`
	StartTimestamp time.Time `json:"start_timestamp,omitempty" bson:"start_timestamp"`
	StartTxHash    string    `json:"start_tx_hash,omitempty" bson:"start_tx_hash"`
	EndHeight      uint64    `json:"end_height,omitempty" bson:"end_height"`
	EndTimestamp   time.Time `json:"end_timestamp,omitempty" bson:"end_timestamp"`
	EndTxHash      string    `json:"end_tx_hash,omitempty" bson:"end_tx_hash"`

	Payment *types.Coin `json:"payment,omitempty" bson:"payment"`
	Rating  uint64      `json:"rating,omitempty" bson:"rating"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`
}

func (s *Session) String() string {
	return utils.MustMarshalIndentToString(s)
}
