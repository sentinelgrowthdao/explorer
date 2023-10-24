package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type Event struct {
	Type      string    `json:"type,omitempty" bson:"type"`
	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`

	AccAddr        string           `json:"acc_addr,omitempty" bson:"acc_addr,omitempty"`
	Bandwidth      *types.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth,omitempty"`
	Coins          types.Coins      `json:"coins,omitempty" bson:"coins,omitempty"`
	Description    string           `json:"description,omitempty" bson:"description,omitempty"`
	Duration       int64            `json:"duration,omitempty" bson:"duration,omitempty"`
	GigabytePrices types.Coins      `json:"gigabyte_prices,omitempty" bson:"gigabyte_prices,omitempty"`
	GrantedBytes   string           `json:"granted_bytes,omitempty" bson:"granted_bytes,omitempty"`
	HourlyPrices   types.Coins      `json:"hourly_prices,omitempty" bson:"hourly_prices,omitempty"`
	Identity       string           `json:"identity,omitempty" bson:"identity,omitempty"`
	Name           string           `json:"name,omitempty" bson:"name,omitempty"`
	NodeAddr       string           `json:"node_addr,omitempty" bson:"node_addr,omitempty"`
	PlanID         uint64           `json:"plan_id,omitempty" bson:"plan_id,omitempty"`
	ProvAddr       string           `json:"prov_addr,omitempty" bson:"prov_addr,omitempty"`
	ReachError     string           `json:"reach_error,omitempty" bson:"reach_error,omitempty"`
	RemoteURL      string           `json:"remote_url,omitempty" bson:"remote_url,omitempty"`
	SessionID      uint64           `json:"session_id,omitempty" bson:"session_id,omitempty"`
	Status         string           `json:"status,omitempty" bson:"status,omitempty"`
	SubscriptionID uint64           `json:"subscription_id,omitempty" bson:"subscription_id,omitempty"`
	UtilisedBytes  string           `json:"utilised_bytes,omitempty" bson:"utilised_bytes,omitempty"`
	Website        string           `json:"website,omitempty" bson:"website,omitempty"`
}

func (e *Event) String() string {
	return utils.MustMarshalIndentToString(e)
}
