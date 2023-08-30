package plan

import (
	"encoding/json"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
)

type MsgAddRequest struct {
	From     string
	Price    types.Coins
	Validity int64
	Bytes    int64
}

func NewMsgAddRequest(v bson.M) (*MsgAddRequest, error) {
	buf, err := json.Marshal(v["price"])
	if err != nil {
		return nil, err
	}

	var price sdk.Coins
	if err := json.Unmarshal(buf, &price); err != nil {
		return nil, err
	}

	validity, err := time.ParseDuration(v["validity"].(string))
	if err != nil {
		return nil, err
	}

	bytes, err := strconv.ParseInt(v["bytes"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgAddRequest{
		From:     v["from"].(string),
		Price:    types.NewCoins(price),
		Validity: validity.Nanoseconds(),
		Bytes:    bytes,
	}, nil
}

type MsgSetStatusRequest struct {
	From   string
	ID     uint64
	Status string
}

func NewMsgSetStatusRequest(v bson.M) (*MsgSetStatusRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgSetStatusRequest{
		From:   v["from"].(string),
		ID:     id,
		Status: v["status"].(string),
	}, nil
}

type MsgAddNodeRequest struct {
	From    string
	ID      uint64
	Address string
}

func NewMsgAddNodeRequest(v bson.M) (*MsgAddNodeRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgAddNodeRequest{
		From:    v["from"].(string),
		ID:      id,
		Address: v["address"].(string),
	}, nil
}

type MsgRemoveNodeRequest struct {
	From    string
	ID      uint64
	Address string
}

func NewMsgRemoveNodeRequest(v bson.M) (*MsgRemoveNodeRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgRemoveNodeRequest{
		From:    v["from"].(string),
		ID:      id,
		Address: v["address"].(string),
	}, nil
}
