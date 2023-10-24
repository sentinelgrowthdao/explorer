package session

import (
	"encoding/json"
	"strconv"
	"time"

	hubtypes "github.com/sentinel-official/hub/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sentinel-official/explorer/types"
)

type MsgStartRequest struct {
	From        string
	ID          uint64
	NodeAddress string
}

func NewMsgStartRequest(v bson.M) (*MsgStartRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgStartRequest{
		From:        v["from"].(string),
		ID:          id,
		NodeAddress: v["address"].(string),
	}, nil
}

type MsgUpdateDetailsRequest struct {
	From      string
	ID        uint64
	Duration  int64
	Bandwidth *types.Bandwidth
}

func NewMsgUpdateDetailsRequest(v bson.M) (*MsgUpdateDetailsRequest, error) {
	id, err := strconv.ParseUint(v["proof"].(bson.M)["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(v["proof"].(bson.M)["duration"].(string))
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(v["proof"].(bson.M)["bandwidth"])
	if err != nil {
		return nil, err
	}

	var bandwidth hubtypes.Bandwidth
	if err := json.Unmarshal(buf, &bandwidth); err != nil {
		return nil, err
	}

	return &MsgUpdateDetailsRequest{
		From:      v["from"].(string),
		ID:        id,
		Duration:  duration.Nanoseconds(),
		Bandwidth: types.NewBandwidth(&bandwidth),
	}, nil
}

type MsgEndRequest struct {
	ID     uint64
	Rating int64
}

func NewMsgEndRequest(v bson.M) (*MsgEndRequest, error) {
	id, err := strconv.ParseUint(v["id"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	rating, err := strconv.ParseInt(v["rating"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	return &MsgEndRequest{
		ID:     id,
		Rating: rating,
	}, nil
}
