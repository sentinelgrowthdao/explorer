package plan

import (
	"time"

	plantypes "github.com/sentinel-official/hub/x/plan/types"

	commontypes "github.com/sentinel-official/explorer/types/common"
)

type (
	MsgAddRequest struct {
		From     string            `json:"from,omitempty" bson:"from"`
		Price    commontypes.Coins `json:"price,omitempty" bson:"price"`
		Validity time.Duration     `json:"validity,omitempty" bson:"validity"`
		Bytes    int64             `json:"bytes,omitempty" bson:"bytes"`
	}
	MsgSetStatusRequest struct {
		From   string `json:"from,omitempty" bson:"from"`
		ID     uint64 `json:"id,omitempty" bson:"id"`
		Status string `json:"status,omitempty" bson:"status"`
	}
	MsgAddNodeRequest struct {
		From    string `json:"from,omitempty" bson:"from"`
		ID      uint64 `json:"id,omitempty" bson:"id"`
		Address string `json:"address,omitempty" bson:"address"`
	}
	MsgRemoveNodeRequest struct {
		From    string `json:"from,omitempty" bson:"from"`
		ID      uint64 `json:"id,omitempty" bson:"id"`
		Address string `json:"address,omitempty" bson:"address"`
	}
)

func NewMsgAddRequestFromRaw(v *plantypes.MsgAddRequest) *MsgAddRequest {
	return &MsgAddRequest{
		From:     v.From,
		Price:    commontypes.NewCoinsFromRaw(v.Price),
		Validity: v.Validity,
		Bytes:    v.Bytes.Int64(),
	}
}

func NewMsgSetStatusRequestFromRaw(v *plantypes.MsgSetStatusRequest) *MsgSetStatusRequest {
	return &MsgSetStatusRequest{
		From:   v.From,
		ID:     v.Id,
		Status: v.Status.String(),
	}
}

func NewMsgMsgAddNodeRequestFromRaw(v *plantypes.MsgAddNodeRequest) *MsgAddNodeRequest {
	return &MsgAddNodeRequest{
		From:    v.From,
		ID:      v.Id,
		Address: v.Address,
	}
}

func NewMsgRemoveNodeRequestFromRaw(v *plantypes.MsgRemoveNodeRequest) *MsgRemoveNodeRequest {
	return &MsgRemoveNodeRequest{
		From:    v.From,
		ID:      v.Id,
		Address: v.Address,
	}
}
