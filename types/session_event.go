package types

import (
	"time"

	commontypes "github.com/sentinel-official/explorer/types/common"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type SessionEvent struct {
	ID        uint64                 `json:"id,omitempty" bson:"id"`
	Bandwidth *commontypes.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth"`
	Duration  int64                  `json:"duration"`
	Signature string                 `json:"signature,omitempty" bson:"signature"`
	Height    int64                  `json:"height,omitempty" bson:"height"`
	Timestamp time.Time              `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string                 `json:"tx_hash,omitempty" bson:"tx_hash"`
}

func (se *SessionEvent) String() string {
	return explorerutils.MustMarshalIndent(se)
}
