package swap

import (
	swaptypes "github.com/sentinel-official/hub/x/swap/types"
	"github.com/tendermint/tendermint/libs/bytes"
)

type (
	MsgSwapRequest struct {
		From     string `json:"from,omitempty" bson:"from"`
		TxHash   string `json:"tx_hash,omitempty" bson:"tx_hash"`
		Receiver string `json:"receiver,omitempty" bson:"receiver"`
		Amount   int64  `json:"amount,omitempty" bson:"amount"`
	}
)

func NewMsgSwapRequestFromRaw(v *swaptypes.MsgSwapRequest) *MsgSwapRequest {
	return &MsgSwapRequest{
		From:     v.From,
		TxHash:   bytes.HexBytes(v.TxHash).String(),
		Receiver: v.Receiver,
		Amount:   v.Amount.Int64(),
	}
}
