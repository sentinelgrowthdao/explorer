package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	m := make(map[string]sdk.Coins)

	m["hi"] = m["hi"].Add(sdk.Coin{Denom: "hello", Amount: sdk.NewInt(100)})
	fmt.Print(m)
}
