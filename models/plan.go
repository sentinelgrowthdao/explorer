package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type Plan struct {
	ID              uint64      `json:"id,omitempty" bson:"id"`
	ProviderAddress string      `json:"provider_address,omitempty" bson:"provider_address"`
	Price           types.Coins `json:"price,omitempty" bson:"price"`
	Validity        int64       `json:"validity,omitempty" bson:"validity"`
	Bytes           int64       `json:"bytes,omitempty" bson:"bytes"`
	NodeAddresses   []string    `json:"node_addresses,omitempty" bson:"node_addresses"`

	AddHeight    int64     `json:"add_height,omitempty" bson:"add_height"`
	AddTimestamp time.Time `json:"add_timestamp,omitempty" bson:"add_timestamp"`
	AddTxHash    string    `json:"add_tx_hash,omitempty" bson:"add_tx_hash"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`
}

func (p *Plan) String() string {
	return utils.MustMarshalIndentToString(p)
}
