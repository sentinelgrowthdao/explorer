package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/tendermint/tendermint/libs/bytes"
)

func MustHexFromBech32AccAddress(s string) string {
	if s == "" {
		return ""
	}

	addr, err := sdk.AccAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return bytes.HexBytes(addr.Bytes()).String()
}

func MustHexFromBech32ConsAddress(s string) string {
	if s == "" {
		return ""
	}

	addr, err := sdk.ConsAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return bytes.HexBytes(addr.Bytes()).String()
}

func MustHexFromBech32ProvAddress(s string) string {
	if s == "" {
		return ""
	}

	addr, err := hubtypes.ProvAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return bytes.HexBytes(addr.Bytes()).String()
}

func MustHexFromBech32NodeAddress(s string) string {
	if s == "" {
		return ""
	}

	addr, err := hubtypes.NodeAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return bytes.HexBytes(addr.Bytes()).String()
}

func MustHexFromBech32ValAddress(s string) string {
	if s == "" {
		return ""
	}

	addr, err := sdk.ValAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return bytes.HexBytes(addr.Bytes()).String()
}
