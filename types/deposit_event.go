package types

import (
	"time"

	commontypes "github.com/sentinel-official/explorer/types/common"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type DepositEvent struct {
	Address   string            `json:"address,omitempty" bson:"address"`
	Coins     commontypes.Coins `json:"coins,omitempty" bson:"coins"`
	Type      string            `json:"type,omitempty" bson:"type"`
	Height    int64             `json:"height,omitempty" bson:"height"`
	Timestamp time.Time         `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string            `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (de *DepositEvent) String() string {
	return explorerutils.MustMarshalIndent(de)
}
