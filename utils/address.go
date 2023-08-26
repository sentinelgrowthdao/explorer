package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
)

func MustAccAddressFromBech32(s string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return addr
}

func MustNodeAddressFromBech32(s string) hubtypes.NodeAddress {
	addr, err := hubtypes.NodeAddressFromBech32(s)
	if err != nil {
		panic(err)
	}

	return addr
}
