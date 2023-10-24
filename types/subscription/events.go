package subscription

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/explorer/types"
)

type EventUpdateStatus struct {
	ID     uint64
	Status string
}

func NewEventUpdateStatus(v *types.Event) (*EventUpdateStatus, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventUpdateStatus{
		ID:     id,
		Status: v.Attributes["status"],
	}, nil
}

type EventAllocate struct {
	ID            uint64
	Address       string
	GrantedBytes  string
	UtilisedBytes string
}

func NewEventAllocate(v *types.Event) (*EventAllocate, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &EventAllocate{
		ID:            id,
		Address:       v.Attributes["address"],
		GrantedBytes:  v.Attributes["granted_bytes"],
		UtilisedBytes: v.Attributes["utilised_bytes"],
	}, nil
}

func NewEventAllocateFromEvents(v types.Events) (int, *EventAllocate, error) {
	i, e, err := v.Get("sentinel.subscription.v2.EventAllocate")
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventAllocate(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
}

type EventPayForPayout struct {
	ID            uint64
	Address       string
	NodeAddress   string
	Payment       *types.Coin
	StakingReward *types.Coin
}

func NewEventPayForPayout(v *types.Event) (*EventPayForPayout, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	payment, err := sdk.ParseCoinNormalized(v.Attributes["payment"])
	if err != nil {
		return nil, err
	}

	stakingReward, err := sdk.ParseCoinNormalized(v.Attributes["staking_reward"])
	if err != nil {
		return nil, err
	}

	return &EventPayForPayout{
		ID:            id,
		Address:       v.Attributes["address"],
		NodeAddress:   v.Attributes["node_address"],
		Payment:       types.NewCoin(&payment),
		StakingReward: types.NewCoin(&stakingReward),
	}, nil
}

type EventPayForPlan struct {
	ID            uint64
	Payment       *types.Coin
	StakingReward *types.Coin
}

func NewEventPayForPlan(v *types.Event) (*EventPayForPlan, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	payment, err := sdk.ParseCoinNormalized(v.Attributes["payment"])
	if err != nil {
		return nil, err
	}

	stakingReward, err := sdk.ParseCoinNormalized(v.Attributes["staking_reward"])
	if err != nil {
		return nil, err
	}

	return &EventPayForPlan{
		ID:            id,
		Payment:       types.NewCoin(&payment),
		StakingReward: types.NewCoin(&stakingReward),
	}, nil
}

func NewEventPayForPlanFromEvents(v types.Events) (int, *EventPayForPlan, error) {
	i, e, err := v.Get("sentinel.subscription.v2.EventPayForPlan")
	if err != nil {
		return 0, nil, err
	}

	item, err := NewEventPayForPlan(e)
	if err != nil {
		return 0, nil, err
	}

	return i, item, nil
}

type EventPayForSession struct {
	ID            uint64
	Payment       *types.Coin
	StakingReward *types.Coin
}

func NewEventPayForSession(v *types.Event) (*EventPayForSession, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	payment, err := sdk.ParseCoinNormalized(v.Attributes["payment"])
	if err != nil {
		return nil, err
	}

	stakingReward, err := sdk.ParseCoinNormalized(v.Attributes["staking_reward"])
	if err != nil {
		return nil, err
	}

	return &EventPayForSession{
		ID:            id,
		Payment:       types.NewCoin(&payment),
		StakingReward: types.NewCoin(&stakingReward),
	}, nil
}

type EventRefund struct {
	ID     uint64
	Amount *types.Coin
}

func NewEventRefund(v *types.Event) (*EventRefund, error) {
	id, err := strconv.ParseUint(v.Attributes["id"], 10, 64)
	if err != nil {
		return nil, err
	}

	amount, err := sdk.ParseCoinNormalized(v.Attributes["amount"])
	if err != nil {
		return nil, err
	}

	return &EventRefund{
		ID:     id,
		Amount: types.NewCoin(&amount),
	}, nil
}
