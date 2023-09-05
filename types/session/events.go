package session

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/types"
)

type EventStartSession struct {
	ID uint64
}

func NewEventStartSession(v *types.Event) (*EventStartSession, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventStartSession{
		ID: id,
	}, nil
}

func NewEventStartSessionFromEvents(v types.Events) (*EventStartSession, error) {
	e, err := v.Get("sentinel.session.v1.EventStartSession")
	if err != nil {
		return nil, err
	}

	return NewEventStartSession(e)
}

type EventEndSession struct {
	ID     uint64
	Status string
}

func NewEventEndSession(v *types.Event) (*EventEndSession, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventEndSession{
		ID:     id,
		Status: v.Attributes["status"],
	}, nil
}

func NewEventEndSessionFromEvents(v types.Events) (*EventEndSession, error) {
	e, err := v.Get("sentinel.session.v1.EventEndSession")
	if err != nil {
		return nil, err
	}

	return NewEventEndSession(e)
}

type EventPay struct {
	ID      uint64
	Payment *types.Coin
}

func NewEventPay(v *types.Event) (*EventPay, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(v.Attributes["payment"])
	if err != nil {
		return nil, err
	}

	var payment sdk.Coin
	if err := json.Unmarshal(buf, &payment); err != nil {
		return nil, err
	}

	return &EventPay{
		ID:      id,
		Payment: types.NewCoin(&payment),
	}, nil
}

func NewEventPayFromEvents(v types.Events) (*EventPay, error) {
	e, err := v.Get("sentinel.session.v1.EventPay")
	if err != nil {
		return nil, err
	}

	return NewEventPay(e)
}
