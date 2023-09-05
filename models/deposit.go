package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Deposit struct {
	Address string      `json:"address,omitempty" bson:"address"`
	Coins   types.Coins `json:"coins,omitempty" bson:"coins"`

	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (d *Deposit) String() string {
	return explorerutils.MustMarshalIndentToString(d)
}
