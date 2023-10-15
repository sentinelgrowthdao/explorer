package subscription

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
)

type MsgSubscribeToNodeRequest struct {
	From    string
	Address string
	Deposit *types.Coin
}

func NewMsgSubscribeToNodeRequest(v bson.M) (*MsgSubscribeToNodeRequest, error) {
	buf, err := json.Marshal(v["deposit"])
	if err != nil {
		return nil, err
	}

	var deposit sdk.Coin
	if err := json.Unmarshal(buf, &deposit); err != nil {
		return nil, err
	}

	return &MsgSubscribeToNodeRequest{
		From:    v["from"].(string),
		Address: v["address"].(string),
		Deposit: types.NewCoin(&deposit),
	}, nil
}

type MsgCancelRequest struct {
	ID uint64
}

func NewMsgCancelRequest(v bson.M) (*MsgCancelRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgCancelRequest{
		ID: id,
	}, nil
}

type MsgAddQuotaRequest struct {
	From    string
	ID      uint64
	Address string
	Bytes   string
}

func NewMsgAddQuotaRequest(v bson.M) (*MsgAddQuotaRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgAddQuotaRequest{
		From:    v["from"].(string),
		ID:      id,
		Address: v["address"].(string),
		Bytes:   v["bytes"].(string),
	}, nil
}

type MsgUpdateQuotaRequest struct {
	From    string
	ID      uint64
	Address string
	Bytes   string
}

func NewMsgUpdateQuotaRequest(v bson.M) (*MsgUpdateQuotaRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgUpdateQuotaRequest{
		From:    v["from"].(string),
		ID:      id,
		Address: v["address"].(string),
		Bytes:   v["bytes"].(string),
	}, nil
}
