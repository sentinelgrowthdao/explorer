package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseOperation func(ctx mongo.SessionContext) error

var (
	EncCfg = hub.MakeEncodingConfig()
)

type ABCIEvent struct {
	Type       string            `json:"type,omitempty" bson:"type"`
	Attributes map[string]string `json:"attributes,omitempty" bson:"attributes"`
}

func NewABCIEvent(v *abcitypes.Event) *ABCIEvent {
	item := &ABCIEvent{
		Type:       v.Type,
		Attributes: make(map[string]string),
	}

	for _, x := range v.Attributes {
		vLen := len(x.Value)
		if vLen >= 2 {
			if x.Value[0] == '"' && x.Value[vLen-1] == '"' {
				x.Value = x.Value[1 : vLen-1]
			}
		}

		item.Attributes[string(x.Key)] = string(x.Value)
	}

	return item
}

type ABCIEvents []*ABCIEvent

func NewABCIEvents(v []abcitypes.Event) ABCIEvents {
	items := make(ABCIEvents, 0, len(v))
	for _, item := range v {
		items = append(items, NewABCIEvent(&item))
	}

	return items
}

type StringEvent struct {
	Type       string            `json:"type,omitempty" bson:"type"`
	Attributes map[string]string `json:"attributes,omitempty" bson:"attributes"`
}

func NewStringEvent(v *sdk.StringEvent) *StringEvent {
	item := &StringEvent{
		Type:       v.Type,
		Attributes: make(map[string]string),
	}

	for _, x := range v.Attributes {
		item.Attributes[x.Key] = x.Value
	}

	return item
}

type StringEvents []*StringEvent

func NewStringEvents(v []sdk.StringEvent) StringEvents {
	items := make(StringEvents, 0, len(v))
	for _, item := range v {
		items = append(items, NewStringEvent(&item))
	}

	return items
}

func (se StringEvents) Get(s string) (*StringEvent, error) {
	for i := 0; i < len(se); i++ {
		if se[i].Type == s {
			return se[i], nil
		}
	}

	return nil, fmt.Errorf("event %s does not exist", s)
}

type ABCIMessageLog struct {
	Index  uint32       `json:"index,omitempty" bson:"index"`
	Log    string       `json:"log,omitempty" bson:"log"`
	Events StringEvents `json:"events,omitempty" bson:"events"`
}

func NewABCIMessageLog(v *sdk.ABCIMessageLog) *ABCIMessageLog {
	return &ABCIMessageLog{
		Index:  v.MsgIndex,
		Log:    replacer.Replace(v.Log),
		Events: NewStringEvents(v.Events),
	}
}

type ABCIMessageLogs []*ABCIMessageLog

func NewABCIMessageLogs(s string) ABCIMessageLogs {
	var v sdk.ABCIMessageLogs
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}

	items := make(ABCIMessageLogs, 0, len(v))
	for _, item := range v {
		items = append(items, NewABCIMessageLog(&item))
	}

	return items
}

type Coin struct {
	Denom  string `json:"denom,omitempty" bson:"denom"`
	Amount int64  `json:"amount,omitempty" bson:"amount"`
}

func NewCoin(v *sdk.Coin) *Coin {
	return &Coin{
		Denom:  v.Denom,
		Amount: v.Amount.Int64(),
	}
}

type Coins []*Coin

func NewCoins(v sdk.Coins) Coins {
	items := make(Coins, 0, v.Len())
	for _, c := range v {
		items = append(items, NewCoin(&c))
	}

	return items
}

type Bandwidth struct {
	Upload   int64 `json:"upload,omitempty" bson:"upload"`
	Download int64 `json:"download,omitempty" bson:"download"`
}

func NewBandwidth(v *hubtypes.Bandwidth) *Bandwidth {
	return &Bandwidth{
		Upload:   v.Upload.Int64(),
		Download: v.Download.Int64(),
	}
}
