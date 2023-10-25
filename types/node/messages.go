package node

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type MsgRegisterRequest struct {
	From           string
	GigabytePrices types.Coins
	HourlyPrices   types.Coins
	RemoteURL      string
}

func NewMsgRegisterRequest(v bson.M) (*MsgRegisterRequest, error) {
	buf, err := json.Marshal(v["gigabyte_prices"])
	if err != nil {
		return nil, err
	}

	var gigabytePrices sdk.Coins
	if err := json.Unmarshal(buf, &gigabytePrices); err != nil {
		return nil, err
	}

	buf, err = json.Marshal(v["hourly_prices"])
	if err != nil {
		return nil, err
	}

	var hourlyPrices sdk.Coins
	if err := json.Unmarshal(buf, &hourlyPrices); err != nil {
		return nil, err
	}

	return &MsgRegisterRequest{
		From:           v["from"].(string),
		GigabytePrices: types.NewCoins(gigabytePrices),
		HourlyPrices:   types.NewCoins(hourlyPrices),
		RemoteURL:      v["remote_url"].(string),
	}, nil
}

func (msg *MsgRegisterRequest) NodeAddr() hubtypes.NodeAddress {
	addr := utils.MustAccAddressFromBech32(msg.From)
	return addr.Bytes()
}

type MsgUpdateDetailsRequest struct {
	From           string
	GigabytePrices types.Coins
	HourlyPrices   types.Coins
	RemoteURL      string
}

func NewMsgUpdateDetailsRequest(v bson.M) (*MsgUpdateDetailsRequest, error) {
	buf, err := json.Marshal(v["gigabyte_prices"])
	if err != nil {
		return nil, err
	}

	var gigabytePrices sdk.Coins
	if err := json.Unmarshal(buf, &gigabytePrices); err != nil {
		return nil, err
	}

	buf, err = json.Marshal(v["hourly_prices"])
	if err != nil {
		return nil, err
	}

	var hourlyPrices sdk.Coins
	if err := json.Unmarshal(buf, &hourlyPrices); err != nil {
		return nil, err
	}

	return &MsgUpdateDetailsRequest{
		From:           v["from"].(string),
		GigabytePrices: types.NewCoins(gigabytePrices),
		HourlyPrices:   types.NewCoins(hourlyPrices),
		RemoteURL:      v["remote_url"].(string),
	}, nil
}

type MsgUpdateStatusRequest struct {
	From   string
	Status string
}

func NewMsgUpdateStatusRequest(v bson.M) (*MsgUpdateStatusRequest, error) {
	return &MsgUpdateStatusRequest{
		From:   v["from"].(string),
		Status: v["status"].(string),
	}, nil
}

type MsgSubscribeRequest struct {
	From        string
	NodeAddress string
	Gigabytes   int64
	Hours       int64
	Denom       string
}

func NewMsgSubscribeRequest(v bson.M) (*MsgSubscribeRequest, error) {
	gigabytes, err := strconv.ParseInt(v["gigabytes"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	hours, err := strconv.ParseInt(v["hours"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgSubscribeRequest{
		From:        v["from"].(string),
		NodeAddress: v["node_address"].(string),
		Gigabytes:   gigabytes,
		Hours:       hours,
		Denom:       v["denom"].(string),
	}, nil
}
