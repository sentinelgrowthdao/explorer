package plan

import (
	"encoding/json"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
)

type MsgCreateRequest struct {
	From      string
	Duration  int64
	Gigabytes int64
	Prices    types.Coins
}

func NewMsgCreateRequest(v bson.M) (*MsgCreateRequest, error) {
	duration, err := time.ParseDuration(v["duration"].(string))
	if err != nil {
		return nil, err
	}

	gigabytes, err := strconv.ParseInt(v["gigabytes"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(v["prices"])
	if err != nil {
		return nil, err
	}

	var prices sdk.Coins
	if err := json.Unmarshal(buf, &prices); err != nil {
		return nil, err
	}

	return &MsgCreateRequest{
		From:      v["from"].(string),
		Gigabytes: gigabytes,
		Duration:  duration.Nanoseconds(),
		Prices:    types.NewCoins(prices),
	}, nil
}

type MsgUpdateStatusRequest struct {
	From   string
	ID     uint64
	Status string
}

func NewMsgUpdateStatusRequest(v bson.M) (*MsgUpdateStatusRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgUpdateStatusRequest{
		From:   v["from"].(string),
		ID:     id,
		Status: v["status"].(string),
	}, nil
}

type MsgLinkNodeRequest struct {
	From        string
	ID          uint64
	NodeAddress string
}

func NewMsgLinkNodeRequest(v bson.M) (*MsgLinkNodeRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgLinkNodeRequest{
		From:        v["from"].(string),
		ID:          id,
		NodeAddress: v["node_address"].(string),
	}, nil
}

type MsgUnlinkNodeRequest struct {
	From        string
	ID          uint64
	NodeAddress string
}

func NewMsgUnlinkNodeRequest(v bson.M) (*MsgUnlinkNodeRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgUnlinkNodeRequest{
		From:        v["from"].(string),
		ID:          id,
		NodeAddress: v["node_address"].(string),
	}, nil
}

type MsgSubscribeRequest struct {
	From  string
	ID    uint64
	Denom string
}

func NewMsgSubscribeRequest(v bson.M) (*MsgSubscribeRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgSubscribeRequest{
		From:  v["from"].(string),
		ID:    id,
		Denom: v["denom"].(string),
	}, nil
}
