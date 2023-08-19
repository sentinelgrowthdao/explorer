package types

import (
	"time"

	commontypes "github.com/sentinel-official/explorer/types/common"
	explorerutils "github.com/sentinel-official/explorer/utils"
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
	Address        string            `json:"address,omitempty" bson:"address"`
	GigabytePrices commontypes.Coins `json:"gigabyte_prices,omitempty" bson:"gigabyte_prices"`
	HourlyPrices   commontypes.Coins `json:"hourly_prices,omitempty" bson:"hourly_prices"`
	RemoteURL      string            `json:"remote_url,omitempty" bson:"remote_url"`

	JoinHeight    int64     `json:"join_height,omitempty" bson:"join_height"`
	JoinTimestamp time.Time `json:"join_timestamp,omitempty" bson:"join_timestamp"`
	JoinTxHash    string    `json:"join_tx_hash,omitempty" bson:"join_tx_hash"`

	Bandwidth              *commontypes.Bandwidth `json:"bandwidth,omitempty" bson:"bandwidth"`
	Handshake              *NodeHandshake         `json:"handshake,omitempty" bson:"handshake"`
	IntervalSetSessions    int64                  `json:"interval_set_sessions,omitempty" bson:"interval_set_sessions"`
	IntervalUpdateSessions int64                  `json:"interval_update_sessions,omitempty" bson:"interval_update_sessions"`
	IntervalUpdateStatus   int64                  `json:"interval_update_status,omitempty" bson:"interval_update_status"`
	Location               *NodeLocation          `json:"location,omitempty" bson:"location"`
	Moniker                string                 `json:"moniker,omitempty" bson:"moniker"`
	Peers                  int                    `json:"peers,omitempty" bson:"peers"`
	QOS                    *NodeQOS               `json:"qos,omitempty" bson:"qos"`
	Type                   uint64                 `json:"type,omitempty" bson:"type"`
	Version                string                 `json:"version,omitempty" bson:"version"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`

	ReachStatus *NodeReachStatus `json:"reach_status,omitempty" bson:"reach_status"`
}

func (n *Node) String() string {
	return explorerutils.MustMarshalIndent(n)
}
