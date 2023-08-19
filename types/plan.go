package types

import (
	"time"

	commontypes "github.com/sentinel-official/explorer/types/common"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Plan struct {
	ID              uint64            `json:"id,omitempty" bson:"id"`
	ProviderAddress string            `json:"provider_address,omitempty" bson:"provider_address"`
	Prices          commontypes.Coins `json:"prices,omitempty" bson:"prices"`
	Duration        int64             `json:"duration,omitempty" bson:"duration"`
	Gigabytes       int64             `json:"gigabytes,omitempty" bson:"gigabytes"`
	NodeAddresses   []string          `json:"node_addresses,omitempty" bson:"node_addresses"`

	CreateHeight    int64     `json:"create_height,omitempty" bson:"create_height"`
	CreateTimestamp time.Time `json:"create_timestamp,omitempty" bson:"create_timestamp"`
	CreateTxHash    string    `json:"create_tx_hash,omitempty" bson:"create_tx_hash"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`
}

func (p *Plan) String() string {
	return explorerutils.MustMarshalIndent(p)
}
