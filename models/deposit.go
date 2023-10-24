package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Deposit struct {
	Addr  string      `json:"addr,omitempty" bson:"addr"`
	Coins types.Coins `json:"coins,omitempty" bson:"coins"`

	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func NewDeposit() *Deposit {
	return &Deposit{
		Coins: types.NewCoins(nil),
	}
}

func (d *Deposit) String() string {
	return explorerutils.MustMarshalIndentToString(d)
}
