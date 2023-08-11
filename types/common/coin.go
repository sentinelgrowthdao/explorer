package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Coin struct {
	Denom  string `json:"denom,omitempty" bson:"denom"`
	Amount int64  `json:"amount,omitempty" bson:"amount"`
}

func NewCoinFromRaw(v *sdk.Coin) *Coin {
	return &Coin{
		Denom:  v.Denom,
		Amount: v.Amount.Int64(),
	}
}

func (c Coin) Raw() sdk.Coin {
	return sdk.Coin{
		Denom:  c.Denom,
		Amount: sdk.NewInt(c.Amount),
	}
}

type Coins []*Coin

func NewCoinsFromRaw(v sdk.Coins) Coins {
	items := make(Coins, 0, v.Len())
	for _, c := range v {
		items = append(items, NewCoinFromRaw(&c))
	}

	return items
}

func (c Coins) Raw() sdk.Coins {
	items := make(sdk.Coins, 0, len(c))
	for _, item := range c {
		items = append(items, item.Raw())
	}

	return items
}
