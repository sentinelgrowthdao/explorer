package models

import (
	"math/big"
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type Event struct {
	Type      string    `json:"type,omitempty" bson:"type"`
	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	TxHash    string    `json:"tx_hash,omitempty" bson:"tx_hash"`

	AccAddress     string           `json:"acc_address,omitempty" bson:"acc_address,omitempty"`
	Allocated      *big.Int         `bson:"allocated,omitempty" bson:"allocated,omitempty"`
	Bandwidth      *types.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth,omitempty"`
	Coins          types.Coins      `json:"coins,omitempty" bson:"coins,omitempty"`
	Consumed       *big.Int         `json:"consumed,omitempty" bson:"consumed,omitempty"`
	Description    string           `json:"description,omitempty" bson:"description,omitempty"`
	Duration       int64            `json:"duration,omitempty" bson:"duration,omitempty"`
	Free           *big.Int         `json:"free,omitempty" bson:"free,omitempty"`
	Identity       string           `json:"identity,omitempty" bson:"identity,omitempty"`
	Name           string           `json:"name,omitempty" bson:"name,omitempty"`
	NodeAddress    string           `json:"node_address,omitempty" bson:"node_address,omitempty"`
	PlanID         uint64           `json:"plan_id,omitempty" bson:"plan_id,omitempty"`
	Price          types.Coins      `json:"price,omitempty" bson:"price,omitempty"`
	ProvAddress    string           `json:"prov_address,omitempty" bson:"prov_address,omitempty"`
	ReachError     string           `json:"reach_error,omitempty" bson:"reach_error,omitempty"`
	RemoteURL      string           `json:"remote_url,omitempty" bson:"remote_url,omitempty"`
	SessionID      uint64           `json:"session_id,omitempty" bson:"session_id,omitempty"`
	Status         string           `json:"status,omitempty" bson:"status,omitempty"`
	SubscriptionID uint64           `json:"subscription_id,omitempty" bson:"subscription_id,omitempty"`
	Website        string           `json:"website,omitempty" bson:"website,omitempty"`
}

func (e *Event) String() string {
	return utils.MustMarshalIndentToString(e)
}
