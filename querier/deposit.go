package querier

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	deposittypes "github.com/sentinel-official/hub/x/deposit/types"
	vpntypes "github.com/sentinel-official/hub/x/vpn/types"
	"google.golang.org/grpc/metadata"

	"github.com/sentinel-official/explorer/types"
)

func (q *Querier) ABCIQueryDeposit(addr sdk.AccAddress, height int64) (*deposittypes.Deposit, error) {
	now := time.Now()
	defer func() {
		log.Println("ABCIQueryDeposit", height, addr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(deposittypes.ModuleName+"/"),
			deposittypes.DepositKey(addr)...,
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

func (q *Querier) GRPCQueryDeposit(addr sdk.AccAddress, height int64) (*deposittypes.Deposit, error) {
	now := time.Now()
	defer func() {
		log.Println("GRPCQueryDeposit", height, addr.String(), time.Since(now))
	}()

	qsc := deposittypes.NewQueryServiceClient(q)
	res, err := qsc.QueryDeposit(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&deposittypes.QueryDepositRequest{
			Address: addr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &res.Deposit, nil
}
