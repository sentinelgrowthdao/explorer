package querier

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	deposittypes "github.com/sentinel-official/hub/x/deposit/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	providertypes "github.com/sentinel-official/hub/x/provider/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	vpntypes "github.com/sentinel-official/hub/x/vpn/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/sentinel-official/explorer/types"
)

type Querier struct {
	*http.HTTP
}

func NewQuerier(remote, wsEndpoint string) (*Querier, error) {
	h, err := http.New(remote, wsEndpoint)
	if err != nil {
		return nil, err
	}

	return &Querier{HTTP: h}, nil
}

func (q *Querier) queryABCI(req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	opts := client.ABCIQueryOptions{
		Height: req.GetHeight(),
		Prove:  req.Prove,
	}

	result, err := q.ABCIQueryWithOptions(context.Background(), req.Path, req.Data, opts)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return q.queryABCI(req)
		}

		return nil, err
	}

	if !result.Response.IsOK() {
		return nil, fmt.Errorf(result.Response.Log)
	}

	return &result.Response, nil
}

func (q *Querier) queryKey(store string, data bytes.HexBytes, height int64) ([]byte, error) {
	req := &abcitypes.RequestQuery{
		Data:   data,
		Path:   fmt.Sprintf("/store/%s/key", store),
		Height: height,
		Prove:  false,
	}

	res, err := q.queryABCI(req)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func (q *Querier) QueryBlock(ctx context.Context, height int64) (*coretypes.ResultBlock, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryBlock", height, time.Since(now))
	}()

	return q.Block(ctx, &height)
}

func (q *Querier) QueryBlockResults(ctx context.Context, height int64) (*coretypes.ResultBlockResults, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryBlockResults", height, time.Since(now))
	}()

	return q.BlockResults(ctx, &height)
}

func (q *Querier) QueryNode(nodeAddr hubtypes.NodeAddress, height int64) (*nodetypes.Node, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryNode", height, nodeAddr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(nodetypes.ModuleName+"/"),
			nodetypes.NodeKey(nodeAddr)...,
		),
		height,
	)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var item nodetypes.Node
	if err := types.EncCfg.Marshaler.Unmarshal(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) QuerySubscription(id uint64, height int64) (*subscriptiontypes.Subscription, error) {
	now := time.Now()
	defer func() {
		log.Println("QuerySubscription", height, id, time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(subscriptiontypes.ModuleName+"/"),
			subscriptiontypes.SubscriptionKey(id)...,
		),
		height,
	)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var item subscriptiontypes.Subscription
	if err := types.EncCfg.Marshaler.Unmarshal(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) QuerySubscriptionQuota(id uint64, accAddr sdk.AccAddress, height int64) (*subscriptiontypes.Quota, error) {
	now := time.Now()
	defer func() {
		log.Println("QuerySubscriptionQuota", height, id, accAddr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(subscriptiontypes.ModuleName+"/"),
			subscriptiontypes.QuotaKey(id, accAddr)...,
		),
		height,
	)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var item subscriptiontypes.Quota
	if err := types.EncCfg.Marshaler.Unmarshal(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) QuerySession(id uint64, height int64) (*sessiontypes.Session, error) {
	now := time.Now()
	defer func() {
		log.Println("QuerySession", height, id, time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(sessiontypes.ModuleName+"/"),
			sessiontypes.SessionKey(id)...,
		),
		height,
	)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var item sessiontypes.Session
	if err := types.EncCfg.Marshaler.Unmarshal(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) QueryProvider(provAddr hubtypes.ProvAddress, height int64) (*providertypes.Provider, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryProvider", height, provAddr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(providertypes.ModuleName+"/"),
			providertypes.ProviderKey(provAddr)...,
		),
		height,
	)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var item providertypes.Provider
	if err := types.EncCfg.Marshaler.Unmarshal(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) QueryDeposit(accAddr sdk.AccAddress, height int64) (*deposittypes.Deposit, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryDeposit", height, accAddr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(deposittypes.ModuleName+"/"),
			deposittypes.DepositKey(accAddr)...,
		),
		height,
	)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	var item deposittypes.Deposit
	if err := types.EncCfg.Marshaler.Unmarshal(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}
