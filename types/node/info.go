package node

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type (
	HandshakeDNS struct {
		Enable bool   `json:"enable,omitempty" bson:"enable"`
		Peers  uint64 `json:"peers,omitempty" bson:"peers"`
	}
	Health struct {
		ClientConfig            []byte    `json:"client_config,omitempty" bson:"client_config"`
		ConfigExchangeError     string    `json:"config_exchange_error,omitempty" bson:"config_exchange_error"`
		ConfigExchangeTimestamp time.Time `json:"config_exchange_timestamp,omitempty" bson:"config_exchange_timestamp"`
		LocationFetchError      string    `json:"location_fetch_error,omitempty" bson:"location_fetch_error"`
		LocationFetchTimestamp  time.Time `json:"location_fetch_timestamp,omitempty" bson:"location_fetch_timestamp"`
		ServerConfig            []byte    `json:"server_config,omitempty" bson:"server_config"`
		SessionID               uint64    `json:"session_id,omitempty" bson:"session_id"`
		StatusFetchError        string    `json:"status_fetch_error,omitempty" bson:"status_fetch_error"`
		StatusFetchTimestamp    time.Time `json:"status_fetch_timestamp,omitempty" bson:"status_fetch_timestamp"`
		SubscriptionID          uint64    `json:"subscription_id,omitempty" bson:"subscription_id"`
	}
	Location struct {
		City      string  `json:"city,omitempty" bson:"city"`
		Country   string  `json:"country,omitempty" bson:"country"`
		Latitude  float64 `json:"latitude,omitempty" bson:"latitude"`
		Longitude float64 `json:"longitude,omitempty" bson:"longitude"`
	}
	QOS struct {
		MaxPeers int `json:"max_peers,omitempty" bson:"max_peers"`
	}
)

type Info struct {
	Bandwidth struct {
		Download int64 `json:"download"`
		Upload   int64 `json:"upload"`
	} `json:"bandwidth"`
	Handshake              *HandshakeDNS `json:"handshake,omitempty"`
	IntervalSetSessions    int64         `json:"interval_set_sessions,omitempty"`
	IntervalUpdateSessions int64         `json:"interval_update_sessions,omitempty"`
	IntervalUpdateStatus   int64         `json:"interval_update_status,omitempty"`
	Location               *Location     `json:"location,omitempty"`
	Moniker                string        `json:"moniker,omitempty"`
	Peers                  int           `json:"peers,omitempty"`
	QOS                    *QOS          `json:"qos,omitempty"`
	Type                   uint64        `json:"type,omitempty"`
	Version                string        `json:"version,omitempty"`
}

func FetchNewInfo(remoteURL string, timeout time.Duration) (*Info, error) {
	urlPath, err := url.JoinPath(remoteURL, "status")
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}

	resp, err := client.Get(urlPath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var m map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	buf, err := json.Marshal(m["result"])
	if err != nil {
		return nil, err
	}

	var v Info
	if err := json.Unmarshal(buf, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
