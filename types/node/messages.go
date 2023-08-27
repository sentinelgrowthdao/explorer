package node

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type MsgRegisterRequest struct {
	From      string
	Provider  string
	Price     types.Coins
	RemoteURL string
}

func NewMsgRegisterRequest(v bson.M) (*MsgRegisterRequest, error) {
	buf, err := json.Marshal(v["price"])
	if err != nil {
		return nil, err
	}

	var price sdk.Coins
	if err := json.Unmarshal(buf, &price); err != nil {
		return nil, err
	}

	return &MsgRegisterRequest{
		From:      v["from"].(string),
		Provider:  v["provider"].(string),
		Price:     types.NewCoins(price),
		RemoteURL: v["remote_url"].(string),
	}, nil
}

func (msg *MsgRegisterRequest) NodeAddress() hubtypes.NodeAddress {
	addr := utils.MustAccAddressFromBech32(msg.From)
	return addr.Bytes()
}

type MsgUpdateRequest struct {
	From      string
	Provider  string
	Price     types.Coins
	RemoteURL string
}

func NewMsgUpdateRequest(v bson.M) (*MsgUpdateRequest, error) {
	buf, err := json.Marshal(v["price"])
	if err != nil {
		return nil, err
	}

	var price sdk.Coins
	if err := json.Unmarshal(buf, &price); err != nil {
		return nil, err
	}

	return &MsgUpdateRequest{
		From:      v["from"].(string),
		Provider:  v["provider"].(string),
		Price:     types.NewCoins(price),
		RemoteURL: v["remote_url"].(string),
	}, nil
}

type MsgSetStatusRequest struct {
	From   string
	Status string
}

func NewMsgSetStatusRequest(v bson.M) (*MsgSetStatusRequest, error) {
	return &MsgSetStatusRequest{
		From:   v["from"].(string),
		Status: v["status"].(string),
	}, nil
}
