package common

import (
	hubtypes "github.com/sentinel-official/hub/types"
)

type Bandwidth struct {
	Upload   int64 `json:"upload,omitempty" bson:"upload"`
	Download int64 `json:"download,omitempty" bson:"download"`
}

func NewBandwidthFromRaw(v *hubtypes.Bandwidth) *Bandwidth {
	return &Bandwidth{
		Upload:   v.Upload.Int64(),
		Download: v.Download.Int64(),
	}
}
