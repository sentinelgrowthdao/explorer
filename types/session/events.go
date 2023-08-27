package session

import (
	"strconv"

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
