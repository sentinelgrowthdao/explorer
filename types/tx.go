package types

import (
	"encoding/json"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	plantypes "github.com/sentinel-official/hub/x/plan/types"
	providertypes "github.com/sentinel-official/hub/x/provider/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	swaptypes "github.com/sentinel-official/hub/x/swap/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/sentinel-official/explorer/types/common"
	nodemessages "github.com/sentinel-official/explorer/types/messages/node"
	planmessages "github.com/sentinel-official/explorer/types/messages/plan"
	providermessages "github.com/sentinel-official/explorer/types/messages/provider"
	sessionmessages "github.com/sentinel-official/explorer/types/messages/session"
	subscriptionmessages "github.com/sentinel-official/explorer/types/messages/subscription"
	swapmessages "github.com/sentinel-official/explorer/types/messages/swap"
	explorerutils "github.com/sentinel-official/explorer/utils"
)

var (
	stringReplacer = strings.NewReplacer(`"\"`, `"`, `\""`, `"`)
)

type Message struct {
	Type string      `json:"type,omitempty" bson:"type"`
	Data interface{} `json:"data,omitempty" bson:"data"`
}

func NewMessageFromRaw(v sdk.Msg) *Message {
	item := &Message{
		Type: sdk.MsgTypeURL(v),
		Data: nil,
	}

	switch v := v.(type) {
	case *nodetypes.MsgRegisterRequest:
		item.Data = nodemessages.NewMsgRegisterRequestFromRaw(v)
	case *nodetypes.MsgUpdateDetailsRequest:
		item.Data = nodemessages.NewMsgUpdateDetailsRequestFromRaw(v)
	case *nodetypes.MsgUpdateStatusRequest:
		item.Data = nodemessages.NewMsgUpdateStatusRequestFromRaw(v)
	case *nodetypes.MsgSubscribeRequest:
		item.Data = nodemessages.NewMsgSubscribeRequestFromRaw(v)

	case *plantypes.MsgCreateRequest:
		item.Data = planmessages.NewMsgCreateRequestFromRaw(v)
	case *plantypes.MsgUpdateStatusRequest:
		item.Data = planmessages.NewMsgUpdateStatusRequestFromRaw(v)
	case *plantypes.MsgLinkNodeRequest:
		item.Data = planmessages.NewMsgLinkNodeRequestFromRaw(v)
	case *plantypes.MsgUnlinkNodeRequest:
		item.Data = planmessages.NewMsgUnlinkNodeRequestFromRaw(v)
	case *plantypes.MsgSubscribeRequest:
		item.Data = planmessages.NewMsgSubscribeRequestFromRaw(v)

	case *providertypes.MsgRegisterRequest:
		item.Data = providermessages.NewMsgRegisterRequestFromRaw(v)
	case *providertypes.MsgUpdateRequest:
		item.Data = providermessages.NewMsgMsgUpdateRequestFromRaw(v)

	case *sessiontypes.MsgStartRequest:
		item.Data = sessionmessages.NewMsgStartRequestFromRaw(v)
	case *sessiontypes.MsgUpdateDetailsRequest:
		item.Data = sessionmessages.NewMsgUpdateDetailsRequestFromRaw(v)
	case *sessiontypes.MsgEndRequest:
		item.Data = sessionmessages.NewMsgMsgEndRequestFromRaw(v)

	case *subscriptiontypes.MsgCancelRequest:
		item.Data = subscriptionmessages.NewMsgMsgCancelRequestFromRaw(v)
	case *subscriptiontypes.MsgAllocateRequest:
		item.Data = subscriptionmessages.NewMsgAllocateRequestFromRaw(v)

	case *swaptypes.MsgSwapRequest:
		item.Data = swapmessages.NewMsgSwapRequestFromRaw(v)
	default:

	}

	return item
}

type Messages []*Message

func NewMessagesFromRaw(v []sdk.Msg) Messages {
	items := make(Messages, 0, len(v))
	for _, item := range v {
		items = append(items, NewMessageFromRaw(item))
	}

	return items
}

type Event struct {
	Type       string            `json:"type,omitempty" bson:"type"`
	Attributes map[string]string `json:"attributes,omitempty" bson:"attributes"`
}

func NewEventFromRaw(v *abcitypes.Event) *Event {
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

type Events []*Event

func NewEventsFromRaw(v []abcitypes.Event) Events {
	items := make(Events, 0, len(v))
	for _, item := range v {
		items = append(items, NewEventFromRaw(&item))
	}

	return items
}

type StringEvent struct {
	Type       string            `json:"type,omitempty" bson:"type"`
	Attributes map[string]string `json:"attributes,omitempty" bson:"attributes"`
}

func NewStringEventFromRaw(v *sdk.StringEvent) *StringEvent {
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

func NewStringEventsFromRaw(v []sdk.StringEvent) StringEvents {
	items := make(StringEvents, 0, len(v))
	for _, item := range v {
		items = append(items, NewStringEventFromRaw(&item))
	}

	return items
}

type TxResultABCIMessageLog struct {
	Index  uint32       `json:"index,omitempty" bson:"index"`
	Log    string       `json:"log,omitempty" bson:"log"`
	Events StringEvents `json:"events,omitempty" bson:"events"`
}

func NewTxResultABCIMessageLogFromRaw(v *sdk.ABCIMessageLog) *TxResultABCIMessageLog {
	return &TxResultABCIMessageLog{
		Index:  v.MsgIndex,
		Log:    v.Log,
		Events: NewStringEventsFromRaw(v.Events),
	}
}

type TxResultABCIMessageLogs []*TxResultABCIMessageLog

func NewTxResultABCIMessageLogsFromRaw(s string) TxResultABCIMessageLogs {
	var v sdk.ABCIMessageLogs
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		panic(err)
	}

	items := make(TxResultABCIMessageLogs, 0, len(v))
	for _, item := range v {
		items = append(items, NewTxResultABCIMessageLogFromRaw(&item))
	}

	return items
}

type TxResult struct {
	Code      uint32 `json:"code,omitempty" bson:"code"`
	Codespace string `json:"codespace,omitempty" bson:"codespace"`
	GasWanted int64  `json:"gas_wanted,omitempty" bson:"gas_wanted"`
	GasUsed   int64  `json:"gas_used,omitempty" bson:"gas_used"`
	Info      string `json:"info,omitempty" bson:"info"`
	Logs      string `json:"logs,omitempty" bson:"logs"`
	Events    Events `json:"events,omitempty" bson:"events"`
}

func NewTxResultFromRaw(v *abcitypes.ResponseDeliverTx) *TxResult {
	return &TxResult{
		Code:      v.Code,
		Codespace: v.Codespace,
		GasWanted: v.GasWanted,
		GasUsed:   v.GasUsed,
		Info:      v.Info,
		Logs:      stringReplacer.Replace(v.Log),
		Events:    NewEventsFromRaw(v.Events),
	}
}

type TxSignerInfo struct {
	Address   string `json:"address,omitempty" bson:"address"`
	PublicKey string `json:"public_key,omitempty" bson:"public_key"`
	Sequence  uint64 `json:"sequence,omitempty" bson:"sequence"`
	Mode      string `json:"mode,omitempty" bson:"mode"`
	Signature string `json:"signature,omitempty" bson:"signature"`
}

type TxSignerInfos []*TxSignerInfo

func NewTxSignerInfosFromTx(v authsigning.Tx) TxSignerInfos {
	signatures, err := v.GetSignaturesV2()
	if err != nil {
		panic(err)
	}

	var (
		signers = v.GetSigners()
		items   = make(TxSignerInfos, 0, len(signers))
	)

	for i := 0; i < len(signers); i++ {
		_, signature := authtx.SignatureDataToModeInfoAndSig(signatures[i].Data)

		items = append(
			items,
			&TxSignerInfo{
				Address:   bytes.HexBytes(signers[i].Bytes()).String(),
				PublicKey: bytes.HexBytes(signatures[i].PubKey.Bytes()).String(),
				Sequence:  signatures[i].Sequence,
				Mode:      "",
				Signature: bytes.HexBytes(signature).String(),
			},
		)
	}

	return items
}

type Tx struct {
	Hash          string        `json:"hash,omitempty" bson:"hash"`
	Height        int64         `json:"height,omitempty" bson:"height"`
	Index         int           `json:"index,omitempty" bson:"index"`
	SignerInfos   TxSignerInfos `json:"signer_infos,omitempty" bson:"signer_infos"`
	Fee           common.Coins  `json:"fee,omitempty" bson:"fee"`
	GasLimit      uint64        `json:"gas_limit,omitempty" bson:"gas_limit"`
	Payer         string        `json:"payer,omitempty" bson:"payer"`
	Granter       string        `json:"granter,omitempty" bson:"granter"`
	Messages      Messages      `json:"messages,omitempty" bson:"messages"`
	Memo          string        `json:"memo,omitempty" bson:"memo"`
	TimeoutHeight uint64        `json:"timeout_height,omitempty" bson:"timeout_height"`
	Result        *TxResult     `json:"result,omitempty" bson:"result"`
}

func NewTxFromRaw(v tmtypes.Tx) *Tx {
	tx, err := DecodeTx(v)
	if err != nil {
		return &Tx{
			Hash:          bytes.HexBytes(v.Hash()).String(),
			Height:        0,
			Index:         0,
			SignerInfos:   nil,
			Fee:           nil,
			GasLimit:      0,
			Payer:         "",
			Granter:       "",
			Messages:      nil,
			Memo:          "",
			TimeoutHeight: 0,
			Result:        nil,
		}
	}

	return &Tx{
		Hash:          bytes.HexBytes(v.Hash()).String(),
		Height:        0,
		Index:         0,
		SignerInfos:   NewTxSignerInfosFromTx(tx),
		Fee:           common.NewCoinsFromRaw(tx.GetFee()),
		GasLimit:      tx.GetGas(),
		Payer:         bytes.HexBytes(tx.FeePayer().Bytes()).String(),
		Granter:       bytes.HexBytes(tx.FeeGranter().Bytes()).String(),
		Messages:      NewMessagesFromRaw(tx.GetMsgs()),
		Memo:          tx.GetMemo(),
		TimeoutHeight: tx.GetTimeoutHeight(),
	}
}

func (t *Tx) String() string {
	return explorerutils.MustMarshalIndent(t)
}

func (t *Tx) WithHeight(v int64) *Tx { t.Height = v; return t }
func (t *Tx) WithIndex(v int) *Tx    { t.Index = v; return t }

func (t *Tx) WithResultRaw(v *abcitypes.ResponseDeliverTx) *Tx {
	t.Result = NewTxResultFromRaw(v)
	return t
}

func DecodeTx(v tmtypes.Tx) (authsigning.Tx, error) {
	tx, err := EncCfg.TxConfig.TxDecoder()(v)
	if err != nil {
		return nil, err
	}

	return tx.(authsigning.Tx), nil
}
