package subscription

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

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

type MsgAllocateRequest struct {
	From    string
	ID      uint64
	Address string
	Bytes   string
}

func NewMsgAllocateRequest(v bson.M) (*MsgAllocateRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgAllocateRequest{
		From:    v["from"].(string),
		ID:      id,
		Address: v["address"].(string),
		Bytes:   v["bytes"].(string),
	}, nil
}
