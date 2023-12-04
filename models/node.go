package models

import (
	"time"

	"github.com/sentinel-official/explorer/types"
	nodetypes "github.com/sentinel-official/explorer/types/node"
	"github.com/sentinel-official/explorer/utils"
)

type Node struct {
	Addr           string      `json:"addr,omitempty" bson:"addr"`
	GigabytePrices types.Coins `json:"gigabyte_prices,omitempty" bson:"gigabyte_prices"`
	HourlyPrices   types.Coins `json:"hourly_prices,omitempty" bson:"hourly_prices"`
	RemoteURL      string      `json:"remote_url,omitempty" bson:"remote_url"`

	RegisterHeight    int64     `json:"register_height,omitempty" bson:"register_height"`
	RegisterTimestamp time.Time `json:"register_timestamp,omitempty" bson:"register_timestamp"`
	RegisterTxHash    string    `json:"register_tx_hash,omitempty" bson:"register_tx_hash"`

	InternetSpeed          types.Bandwidth        `json:"internet_speed,omitempty" bson:"internet_speed"`
	HandshakeDNS           nodetypes.HandshakeDNS `json:"handshake_dns,omitempty" bson:"handshake_dns"`
	IntervalSetSessions    int64                  `json:"interval_set_sessions,omitempty" bson:"interval_set_sessions"`
	IntervalUpdateSessions int64                  `json:"interval_update_sessions,omitempty" bson:"interval_update_sessions"`
	IntervalUpdateStatus   int64                  `json:"interval_update_status,omitempty" bson:"interval_update_status"`
	Location               nodetypes.Location     `json:"location,omitempty" bson:"location"`
	Moniker                string                 `json:"moniker,omitempty" bson:"moniker"`
	Peers                  int                    `json:"peers,omitempty" bson:"peers"`
	QOS                    nodetypes.QOS          `json:"qos,omitempty" bson:"qos"`
	Type                   uint64                 `json:"type,omitempty" bson:"type"`
	Version                string                 `json:"version,omitempty" bson:"version"`

	Status          string    `json:"status,omitempty" bson:"status"`
	StatusHeight    int64     `json:"status_height,omitempty" bson:"status_height"`
	StatusTimestamp time.Time `json:"status_timestamp,omitempty" bson:"status_timestamp"`
	StatusTxHash    string    `json:"status_tx_hash,omitempty" bson:"status_tx_hash"`

	Health nodetypes.Health `json:"health,omitempty" bson:"health"`
}

func (n *Node) String() string {
	return utils.MustMarshalIndentToString(n)
}
