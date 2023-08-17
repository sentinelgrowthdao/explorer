package node

import (
	nodetypes "github.com/sentinel-official/hub/x/node/types"

	commontypes "github.com/sentinel-official/explorer/types/common"
)

type (
	MsgRegisterRequest struct {
		From           string            `json:"from,omitempty" bson:"from"`
		GigabytePrices commontypes.Coins `json:"gigabyte_prices,omitempty" bson:"gigabyte_prices"`
		HourlyPrices   commontypes.Coins `json:"hourly_prices,omitempty" bson:"hourly_prices"`
		RemoteURL      string            `json:"remote_url,omitempty" bson:"remote_url"`
	}
	MsgUpdateDetailsRequest struct {
		From           string            `json:"from,omitempty" bson:"from"`
		GigabytePrices commontypes.Coins `json:"gigabyte_prices,omitempty" bson:"gigabyte_prices"`
		HourlyPrices   commontypes.Coins `json:"hourly_prices,omitempty" bson:"hourly_prices"`
		RemoteURL      string            `json:"remote_url,omitempty" bson:"remote_url"`
	}
	MsgUpdateStatusRequest struct {
		From   string `json:"from,omitempty" bson:"from"`
		Status string `json:"status,omitempty" bson:"status"`
	}
	MsgSubscribeRequest struct {
		From        string `json:"from,omitempty" bson:"from"`
		NodeAddress string `json:"node_address,omitempty" bson:"node_address"`
		Gigabytes   int64  `json:"gigabytes,omitempty" bson:"gigabytes"`
		Hours       int64  `json:"hours,omitempty" bson:"hours"`
		Denom       string `json:"denom,omitempty" bson:"denom"`
	}
)

func NewMsgRegisterRequestFromRaw(v *nodetypes.MsgRegisterRequest) *MsgRegisterRequest {
	return &MsgRegisterRequest{
		From:           v.From,
		GigabytePrices: commontypes.NewCoinsFromRaw(v.GigabytePrices),
		HourlyPrices:   commontypes.NewCoinsFromRaw(v.HourlyPrices),
		RemoteURL:      v.RemoteURL,
	}
}

func NewMsgUpdateDetailsRequestFromRaw(v *nodetypes.MsgUpdateDetailsRequest) *MsgUpdateDetailsRequest {
	return &MsgUpdateDetailsRequest{
		From:           v.From,
		GigabytePrices: commontypes.NewCoinsFromRaw(v.GigabytePrices),
		HourlyPrices:   commontypes.NewCoinsFromRaw(v.HourlyPrices),
		RemoteURL:      v.RemoteURL,
	}
}

func NewMsgUpdateStatusRequestFromRaw(v *nodetypes.MsgUpdateStatusRequest) *MsgUpdateStatusRequest {
	return &MsgUpdateStatusRequest{
		From:   v.From,
		Status: v.Status.String(),
	}
}

func NewMsgSubscribeRequestFromRaw(v *nodetypes.MsgSubscribeRequest) *MsgSubscribeRequest {
	return &MsgSubscribeRequest{
		From:        v.From,
		NodeAddress: v.NodeAddress,
		Gigabytes:   v.Gigabytes,
		Hours:       v.Hours,
		Denom:       v.Denom,
	}
}
