package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sentinel-official/explorer/utils"
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
	Amount string `json:"amount,omitempty" bson:"amount"`
}

func NewCoin(v *sdk.Coin) *Coin {
	return &Coin{
		Denom:  v.Denom,
		Amount: v.Amount.String(),
	}
}

func (c *Coin) String() string {
	return fmt.Sprintf("%s%s;", c.Amount, c.Denom)
}

func (c *Coin) Copy() *Coin {
	return &Coin{
		Denom:  c.Denom,
		Amount: c.Amount,
	}
}

func (c *Coin) Add(v string) *Coin {
	a1 := utils.MustIntFromString(c.Amount)
	a2 := utils.MustIntFromString(v)

	c.Amount = a1.Add(a2).String()
	return c
}

type Coins []*Coin

func (c Coins) Len() int           { return len(c) }
func (c Coins) Less(i, j int) bool { return c[i].Denom < c[j].Denom }
func (c Coins) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Coins) Sort() Coins        { sort.Sort(c); return c }
func (c Coins) IsSorted() bool     { return sort.IsSorted(c) }

func (c Coins) IndexOf(v string) int {
	return sort.Search(c.Len(), func(i int) bool {
		return c[i].Denom >= v
	})
}

func NewCoins(v sdk.Coins) Coins {
	items := make(Coins, 0, v.Len())
	for _, c := range v {
		items = append(items, NewCoin(&c))
	}

	return items.Sort()
}

func (c Coins) Copy() (n Coins) {
	for i := 0; i < len(c); i++ {
		n = append(n, c[i].Copy())
	}

	return n
}

func (c Coins) Add(v ...*Coin) (n Coins) {
	if !c.IsSorted() {
		panic("coins must be sorted")
	}

	n = c.Copy()
	for i := 0; i < len(v); i++ {
		index := n.IndexOf(v[i].Denom)
		if index < len(n) && n[index].Denom == v[i].Denom {
			n[index] = n[index].Add(v[i].Amount)
		} else {
			n = append(n, v[i]).Sort()
		}
	}

	return n
}

type Bandwidth struct {
	Upload   string `json:"upload,omitempty" bson:"upload"`
	Download string `json:"download,omitempty" bson:"download"`
}

func NewBandwidth(v *hubtypes.Bandwidth) *Bandwidth {
	if v == nil {
		return &Bandwidth{}
	}

	return &Bandwidth{
		Upload:   v.Upload.String(),
		Download: v.Download.String(),
	}
}

func (b *Bandwidth) Add(v *Bandwidth) *Bandwidth {
	bu := utils.MustIntFromString(b.Upload)
	bd := utils.MustIntFromString(b.Download)

	vu := utils.MustIntFromString(v.Upload)
	vd := utils.MustIntFromString(v.Download)

	b.Upload = bu.Add(vu).String()
	b.Download = bd.Add(vd).String()
	return b
}

func (b *Bandwidth) Sub(v *Bandwidth) *Bandwidth {
	bu := utils.MustIntFromString(b.Upload)
	bd := utils.MustIntFromString(b.Download)

	vu := utils.MustIntFromString(v.Upload)
	vd := utils.MustIntFromString(v.Download)

	b.Upload = bu.Sub(vu).String()
	b.Download = bd.Sub(vd).String()
	return b
}
