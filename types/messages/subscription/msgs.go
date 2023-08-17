package subscription

import (
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
)

type (
	MsgCancelRequest struct {
		From string `json:"from,omitempty" bson:"from"`
		ID   uint64 `json:"id,omitempty" bson:"id"`
	}
	MsgAllocateRequest struct {
		From    string `json:"from,omitempty" bson:"from"`
		ID      uint64 `json:"id,omitempty" bson:"id"`
		Address string `json:"address,omitempty" bson:"address"`
		Bytes   int64  `json:"bytes,omitempty" bson:"bytes"`
	}
)

func NewMsgMsgCancelRequestFromRaw(v *subscriptiontypes.MsgCancelRequest) *MsgCancelRequest {
	return &MsgCancelRequest{
		From: v.From,
		ID:   v.ID,
	}
}

func NewMsgAllocateRequestFromRaw(v *subscriptiontypes.MsgAllocateRequest) *MsgAllocateRequest {
	return &MsgAllocateRequest{
		From:    v.From,
		ID:      v.ID,
		Address: v.Address,
		Bytes:   v.Bytes.Int64(),
	}
}
