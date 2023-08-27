package types

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	DatabaseOperation func(ctx mongo.SessionContext) error
)

var (
	Replacer = strings.NewReplacer(`"\"`, `"`, `\""`, `"`)
	EncCfg   = hub.MakeEncodingConfig()
)

type Event struct {
	Type       string            `json:"type,omitempty" bson:"type"`
	Attributes map[string]string `json:"attributes,omitempty" bson:"attributes"`
}

func NewEventFromABCIEvent(v *abcitypes.Event) *Event {
	item := &Event{
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

func NewEventFromStringEvent(v *sdk.StringEvent) *Event {
	item := &Event{
		Type:       v.Type,
		Attributes: make(map[string]string),
	}

	for _, x := range v.Attributes {
		item.Attributes[x.Key] = x.Value
	}

	return item
}

type Events []*Event

func NewEventsFromABCIEvents(v []abcitypes.Event) Events {
	items := make(Events, 0, len(v))
	for _, item := range v {
		items = append(items, NewEventFromABCIEvent(&item))
	}

	return items
}

func NewEventsFromStringEvents(v []sdk.StringEvent) Events {
	items := make(Events, 0, len(v))
	for _, item := range v {
		items = append(items, NewEventFromStringEvent(&item))
	}

	return items
}

func (e Events) Get(s string) (*Event, error) {
	for i := 0; i < len(e); i++ {
		if e[i].Type == s {
			return e[i], nil
		}
	}

	return nil, fmt.Errorf("event %s does not exist", s)
}

type ABCIMessageLog struct {
	Index  uint32 `json:"index,omitempty" bson:"index"`
	Log    string `json:"log,omitempty" bson:"log"`
	Events Events `json:"events,omitempty" bson:"events"`
}

func NewABCIMessageLog(v *sdk.ABCIMessageLog) *ABCIMessageLog {
	return &ABCIMessageLog{
		Index:  v.MsgIndex,
		Log:    Replacer.Replace(v.Log),
		Events: NewEventsFromStringEvents(v.Events),
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
