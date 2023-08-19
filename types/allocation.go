package types

import (
	explorerutils "github.com/sentinel-official/explorer/utils"
)

type Allocation struct {
	ID            uint64 `json:"id,omitempty" bson:"id"`
	Address       string `json:"address,omitempty" bson:"address"`
	UtilisedBytes int64  `json:"utilised_bytes,omitempty" bson:"utilised_bytes"`
	GrantedBytes  int64  `json:"granted_bytes,omitempty" bson:"granted_bytes"`
}

func (a *Allocation) String() string {
	return explorerutils.MustMarshalIndent(a)
}
