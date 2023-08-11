package types

import (
	"time"

	commontypes "github.com/sentinel-official/explorer/types/common"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Deposit struct {
	Address   string            `json:"address,omitempty" bson:"address"`
	Coins     commontypes.Coins `json:"coins,omitempty" bson:"coins"`
	Height    int64             `json:"height,omitempty" bson:"height"`
	Timestamp time.Time         `json:"timestamp,omitempty" bson:"timestamp"`
}

func (d *Deposit) String() string {
	return explorerutils.MustMarshalIndent(d)
}
