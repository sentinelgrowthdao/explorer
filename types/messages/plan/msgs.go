package plan

import (
	"time"

	plantypes "github.com/sentinel-official/hub/x/plan/types"

	commontypes "github.com/sentinel-official/explorer/types/common"
)

type (
	MsgCreateRequest struct {
		From      string            `json:"from,omitempty" bson:"from"`
		Duration  time.Duration     `json:"duration,omitempty" bson:"duration"`
		Gigabytes int64             `json:"gigabytes,omitempty" bson:"gigabytes"`
		Prices    commontypes.Coins `json:"prices,omitempty" bson:"prices"`
	}
	MsgUpdateStatusRequest struct {
		From   string `json:"from,omitempty" bson:"from"`
		ID     uint64 `json:"id,omitempty" bson:"id"`
		Status string `json:"status,omitempty" bson:"status"`
	}
	MsgLinkNodeRequest struct {
		From        string `json:"from,omitempty" bson:"from"`
		ID          uint64 `json:"id,omitempty" bson:"id"`
		NodeAddress string `json:"node_address,omitempty" bson:"node_address"`
	}
	MsgUnlinkNodeRequest struct {
		From        string `json:"from,omitempty" bson:"from"`
		ID          uint64 `json:"id,omitempty" bson:"id"`
		NodeAddress string `json:"node_address,omitempty" bson:"node_address"`
	}
	MsgSubscribeRequest struct {
		From  string `json:"from,omitempty" bson:"from"`
		ID    uint64 `json:"id,omitempty" bson:"id"`
		Denom string `json:"denom,omitempty" bson:"denom"`
	}
)

func NewMsgCreateRequestFromRaw(v *plantypes.MsgCreateRequest) *MsgCreateRequest {
	return &MsgCreateRequest{
		From:      v.From,
		Duration:  v.Duration,
		Gigabytes: v.Gigabytes,
		Prices:    commontypes.NewCoinsFromRaw(v.Prices),
	}
}

func NewMsgUpdateStatusRequestFromRaw(v *plantypes.MsgUpdateStatusRequest) *MsgUpdateStatusRequest {
	return &MsgUpdateStatusRequest{
		From:   v.From,
		ID:     v.ID,
		Status: v.Status.String(),
	}
}

func NewMsgLinkNodeRequestFromRaw(v *plantypes.MsgLinkNodeRequest) *MsgLinkNodeRequest {
	return &MsgLinkNodeRequest{
		From:        v.From,
		ID:          v.ID,
		NodeAddress: v.NodeAddress,
	}
}

func NewMsgUnlinkNodeRequestFromRaw(v *plantypes.MsgUnlinkNodeRequest) *MsgUnlinkNodeRequest {
	return &MsgUnlinkNodeRequest{
		From:        v.From,
		ID:          v.ID,
		NodeAddress: v.NodeAddress,
	}
}

func NewMsgSubscribeRequestFromRaw(v *plantypes.MsgSubscribeRequest) *MsgSubscribeRequest {
	return &MsgSubscribeRequest{
		From:  v.From,
		ID:    v.ID,
		Denom: v.Denom,
	}
}
