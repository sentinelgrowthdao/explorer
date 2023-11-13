package querier

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Querier struct {
	codectypes.InterfaceRegistry
	*tmhttp.HTTP
}

func NewQuerier(ir codectypes.InterfaceRegistry, remote, wsEndpoint string) (q *Querier, err error) {
	http, err := tmhttp.New(remote, wsEndpoint)
	if err != nil {
		return nil, err
	}

	return &Querier{
		InterfaceRegistry: ir,
		HTTP:              http,
	}, nil
}

func (q *Querier) queryABCI(req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	opts := client.ABCIQueryOptions{
		Height: req.GetHeight(),
		Prove:  req.Prove,
	}

	result, err := q.ABCIQueryWithOptions(context.TODO(), req.Path, req.Data, opts)
	log.Println("ABCIQueryWithOptions", err)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return q.queryABCI(req)
		}
		if strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
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
