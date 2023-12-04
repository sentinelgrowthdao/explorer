package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MustIntFromString(v string) sdk.Int {
	if v == "" {
		v = "0"
	}

	i, ok := sdk.NewIntFromString(v)
	if !ok {
		panic("not ok")
	}

	return i
}

func MustStringFromInt64(v int64) string {
	return fmt.Sprintf("%d", v)
}
