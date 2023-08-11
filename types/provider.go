package types

import (
	"time"

	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Provider struct {
	Address     string `json:"address,omitempty" bson:"address"`
	Name        string `json:"name,omitempty" bson:"name"`
	Identity    string `json:"identity,omitempty" bson:"identity"`
	Website     string `json:"website,omitempty" bson:"website"`
	Description string `json:"description,omitempty" bson:"description"`

	JoinHeight    int64     `json:"join_height,omitempty" bson:"join_height"`
	JoinTimestamp time.Time `json:"join_timestamp,omitempty" bson:"join_timestamp"`
	JoinTxHash    string    `json:"join_tx_hash,omitempty" bson:"join_tx_hash"`
}

func (p *Provider) String() string {
	return explorerutils.MustMarshalIndent(p)
}
