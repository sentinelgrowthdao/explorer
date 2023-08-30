package querier

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	hubtypes "github.com/sentinel-official/hub/types"
	providertypes "github.com/sentinel-official/hub/x/provider/types"
	vpntypes "github.com/sentinel-official/hub/x/vpn/types"
	"google.golang.org/grpc/metadata"

	"github.com/sentinel-official/explorer/types"
)

func (q *Querier) ABCIQueryProvider(addr hubtypes.ProvAddress, height int64) (*providertypes.Provider, error) {
	now := time.Now()
	defer func() {
		log.Println("ABCIQueryProvider", height, addr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(providertypes.ModuleName+"/"),
			providertypes.ProviderKey(addr)...,
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
	if err := types.EncCfg.Marshaler.UnmarshalBinaryBare(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) GRPCQueryProvider(addr hubtypes.ProvAddress, height int64) (*providertypes.Provider, error) {
	now := time.Now()
	defer func() {
		log.Println("GRPCQueryProvider", height, addr.String(), time.Since(now))
	}()

	qsc := providertypes.NewQueryServiceClient(q)
	res, err := qsc.QueryProvider(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&providertypes.QueryProviderRequest{
			Address: addr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &res.Provider, nil
}
