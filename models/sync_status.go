package models

import (
	"time"

	"github.com/sentinel-official/explorer/utils"
)

type SyncStatus struct {
	AppName   string    `json:"app_name,omitempty" bson:"app_name"`
	Height    int64     `json:"height,omitempty" bson:"height"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
}

func (ss *SyncStatus) String() string {
	return utils.MustMarshalIndentToString(ss)
}
