package session

import (
	"time"

	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	"github.com/tendermint/tendermint/libs/bytes"

	commontypes "github.com/sentinel-official/explorer/types/common"
)

type (
	MsgStartRequest struct {
		From string `json:"from,omitempty" bson:"from"`
		ID   uint64 `json:"id,omitempty" bson:"id"`
		Node string `json:"node,omitempty" bson:"node"`
	}
	MsgUpdateRequest struct {
		From      string                 `json:"from,omitempty" bson:"from"`
		ID        uint64                 `json:"id,omitempty" bson:"id"`
		Bandwidth *commontypes.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth"`
		Duration  time.Duration          `json:"duration,omitempty" bson:"duration"`
		Signature string                 `json:"signature,omitempty" bson:"signature"`
	}
	MsgEndRequest struct {
		From   string `json:"from,omitempty" bson:"from"`
		ID     uint64 `json:"id,omitempty" bson:"id"`
		Rating uint64 `json:"rating,omitempty" bson:"rating"`
	}
)

func NewMsgStartRequestFromRaw(v *sessiontypes.MsgStartRequest) *MsgStartRequest {
	return &MsgStartRequest{
		From: v.From,
		ID:   v.Id,
		Node: v.Node,
	}
}

func NewMsgUpdateRequestFromRaw(v *sessiontypes.MsgUpdateRequest) *MsgUpdateRequest {
	return &MsgUpdateRequest{
		From:      v.From,
		ID:        v.Proof.Id,
		Bandwidth: commontypes.NewBandwidthFromRaw(&v.Proof.Bandwidth),
		Duration:  v.Proof.Duration,
		Signature: bytes.HexBytes(v.Signature).String(),
	}
}

func NewMsgMsgEndRequestFromRaw(v *sessiontypes.MsgEndRequest) *MsgEndRequest {
	return &MsgEndRequest{
		From:   v.From,
		ID:     v.Id,
		Rating: v.Rating,
	}
}
