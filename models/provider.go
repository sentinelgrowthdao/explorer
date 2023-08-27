package models

import (
	"time"

	"github.com/sentinel-official/explorer/utils"
)

type Provider struct {
	Address     string `json:"address,omitempty" bson:"address"`
	Name        string `json:"name,omitempty" bson:"name"`
	Identity    string `json:"identity,omitempty" bson:"identity"`
	Website     string `json:"website,omitempty" bson:"website"`
	Description string `json:"description,omitempty" bson:"description"`

	RegisterHeight    int64     `json:"register_height,omitempty" bson:"register_height"`
	RegisterTimestamp time.Time `json:"register_timestamp,omitempty" bson:"register_timestamp"`
	RegisterTxHash    string    `json:"register_tx_hash,omitempty" bson:"register_tx_hash"`
}

func (p *Provider) String() string {
	return utils.MustMarshalIndentToString(p)
}
