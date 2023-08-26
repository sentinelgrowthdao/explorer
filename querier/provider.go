package querier

import (
	"context"
	"log"
	"strconv"
	"time"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	hubtypes "github.com/sentinel-official/hub/types"
	providertypes "github.com/sentinel-official/hub/x/provider/types"
	"google.golang.org/grpc/metadata"
)

func (q *Querier) QueryProvider(addr hubtypes.ProvAddress, height int64) (*providertypes.Provider, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryProvider", height, addr.String(), time.Since(now))
	}()

	res, err := q.provider.QueryProvider(
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
