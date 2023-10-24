package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	"github.com/sentinel-official/explorer/utils"
)

type (
	NodeHandshake struct {
		Enable bool   `json:"enable,omitempty" bson:"enable"`
		Peers  uint64 `json:"peers,omitempty" bson:"peers"`
	}
	NodeLocation struct {
		City      string  `json:"city,omitempty" bson:"city"`
		Country   string  `json:"country,omitempty" bson:"country"`
		Latitude  float64 `json:"latitude,omitempty" bson:"latitude"`
		Longitude float64 `json:"longitude,omitempty" bson:"longitude"`
	}
	NodeQOS struct {
		MaxPeers int `json:"max_peers,omitempty" bson:"max_peers"`
	}
	NodeReachStatus struct {
		ErrorMessage string    `json:"error_message,omitempty" bson:"error_message"`
		Timestamp    time.Time `json:"timestamp,omitempty" bson:"timestamp"`
	}
)

type Node struct {
	Addr           string      `json:"addr,omitempty" bson:"addr"`
	GigabytePrices types.Coins `json:"gigabyte_prices,omitempty" bson:"gigabyte_prices"`
	HourlyPrices   types.Coins `json:"hourly_prices,omitempty" bson:"hourly_prices"`
	RemoteURL      string      `json:"remote_url,omitempty" bson:"remote_url"`

	RegisterHeight    int64     `json:"register_height,omitempty" bson:"register_height"`
	RegisterTimestamp time.Time `json:"register_timestamp,omitempty" bson:"register_timestamp"`
	RegisterTxHash    string    `json:"register_tx_hash,omitempty" bson:"register_tx_hash"`

	Bandwidth              *types.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth"`
	Handshake              *NodeHandshake   `json:"handshake,omitempty" bson:"handshake"`
	IntervalSetSessions    int64            `json:"interval_set_sessions,omitempty" bson:"interval_set_sessions"`
	IntervalUpdateSessions int64            `json:"interval_update_sessions,omitempty" bson:"interval_update_sessions"`
	IntervalUpdateStatus   int64            `json:"interval_update_status,omitempty" bson:"interval_update_status"`
	Location               *NodeLocation    `json:"location,omitempty" bson:"location"`
	Moniker                string           `json:"moniker,omitempty" bson:"moniker"`
	Peers                  int              `json:"peers,omitempty" bson:"peers"`
	QOS                    *NodeQOS         `json:"qos,omitempty" bson:"qos"`
	Type                   uint64           `json:"type,omitempty" bson:"type"`
	Version                string           `json:"version,omitempty" bson:"version"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`

	ReachStatus *NodeReachStatus `json:"reach_status,omitempty" bson:"reach_status"`
}

func (n *Node) String() string {
	return utils.MustMarshalIndentToString(n)
}
