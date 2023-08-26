package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

func MsgTypeURL(msg sdk.Msg) string {
	return "/" + proto.MessageName(msg)
}
