package node

import (
	nodetypes "github.com/sentinel-official/hub/x/node/types"

	commontypes "github.com/sentinel-official/explorer/types/common"
)

type (
	MsgRegisterRequest struct {
		From      string            `json:"from,omitempty" bson:"from"`
		Provider  string            `json:"provider,omitempty" bson:"provider"`
		Price     commontypes.Coins `json:"price,omitempty" bson:"price"`
		RemoteURL string            `json:"remote_url,omitempty" bson:"remote_url"`
	}
	MsgUpdateRequest struct {
		From      string            `json:"from,omitempty" bson:"from"`
		Provider  string            `json:"provider,omitempty" bson:"provider"`
		Price     commontypes.Coins `json:"price,omitempty" bson:"price"`
		RemoteURL string            `json:"remote_url,omitempty" bson:"remote_url"`
	}
	MsgSetStatusRequest struct {
		From   string `json:"from,omitempty" bson:"from"`
		Status string `json:"status,omitempty" bson:"status"`
	}
)

func NewMsgRegisterRequestFromRaw(v *nodetypes.MsgRegisterRequest) *MsgRegisterRequest {
	return &MsgRegisterRequest{
		From:      v.From,
		Provider:  v.Provider,
		Price:     commontypes.NewCoinsFromRaw(v.Price),
		RemoteURL: v.RemoteURL,
	}
}

func NewMsgMsgUpdateRequestFromRaw(v *nodetypes.MsgUpdateRequest) *MsgUpdateRequest {
	return &MsgUpdateRequest{
		From:      v.From,
		Provider:  v.Provider,
		Price:     commontypes.NewCoinsFromRaw(v.Price),
		RemoteURL: v.RemoteURL,
	}
}

func NewMsgSetStatusRequestFromRaw(v *nodetypes.MsgSetStatusRequest) *MsgSetStatusRequest {
	return &MsgSetStatusRequest{
		From:   v.From,
		Status: v.Status.String(),
	}
}
