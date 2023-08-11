package subscription

import (
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"

	commontypes "github.com/sentinel-official/explorer/types/common"
)

type (
	MsgSubscribeToNodeRequest struct {
		From    string            `json:"from,omitempty" bson:"from"`
		Address string            `json:"address,omitempty" bson:"address"`
		Deposit *commontypes.Coin `json:"deposit,omitempty" bson:"deposit"`
	}
	MsgSubscribeToPlanRequest struct {
		From  string `json:"from,omitempty" bson:"from"`
		ID    uint64 `json:"id,omitempty" bson:"id"`
		Denom string `json:"denom,omitempty" bson:"denom"`
	}
	MsgCancelRequest struct {
		From string `json:"from,omitempty" bson:"from"`
		ID   uint64 `json:"id,omitempty" bson:"id"`
	}
	MsgAddQuotaRequest struct {
		From    string `json:"from,omitempty" bson:"from"`
		ID      uint64 `json:"id,omitempty" bson:"id"`
		Address string `json:"address,omitempty" bson:"address"`
		Bytes   int64  `json:"bytes,omitempty" bson:"bytes"`
	}
	MsgUpdateQuotaRequest struct {
		From    string `json:"from,omitempty" bson:"from"`
		ID      uint64 `json:"id,omitempty" bson:"id"`
		Address string `json:"address,omitempty" bson:"address"`
		Bytes   int64  `json:"bytes,omitempty" bson:"bytes"`
	}
)

func NewMsgSubscribeToNodeRequestFromRaw(v *subscriptiontypes.MsgSubscribeToNodeRequest) *MsgSubscribeToNodeRequest {
	return &MsgSubscribeToNodeRequest{
		From:    v.From,
		Address: v.Address,
		Deposit: commontypes.NewCoinFromRaw(&v.Deposit),
	}
}

func NewMsgSubscribeToPlanRequestFromRaw(v *subscriptiontypes.MsgSubscribeToPlanRequest) *MsgSubscribeToPlanRequest {
	return &MsgSubscribeToPlanRequest{
		From:  v.From,
		ID:    v.Id,
		Denom: v.Denom,
	}
}

func NewMsgMsgCancelRequestFromRaw(v *subscriptiontypes.MsgCancelRequest) *MsgCancelRequest {
	return &MsgCancelRequest{
		From: v.From,
		ID:   v.Id,
	}
}

func NewMsgAddQuotaRequestFromRaw(v *subscriptiontypes.MsgAddQuotaRequest) *MsgAddQuotaRequest {
	return &MsgAddQuotaRequest{
		From:    v.From,
		ID:      v.Id,
		Address: v.Address,
		Bytes:   v.Bytes.Int64(),
	}
}

func NewMsgUpdateQuotaRequestFromRaw(v *subscriptiontypes.MsgUpdateQuotaRequest) *MsgUpdateQuotaRequest {
	return &MsgUpdateQuotaRequest{
		From:    v.From,
		ID:      v.Id,
		Address: v.Address,
		Bytes:   v.Bytes.Int64(),
	}
}
