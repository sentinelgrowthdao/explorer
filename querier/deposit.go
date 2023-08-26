package querier

import (
	"context"
	"log"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	deposittypes "github.com/sentinel-official/hub/x/deposit/types"
	"google.golang.org/grpc/metadata"
)

func (q *Querier) QueryDeposit(addr sdk.AccAddress, height int64) (*deposittypes.Deposit, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryDeposit", height, addr.String(), time.Since(now))
	}()

	res, err := q.deposit.QueryDeposit(
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
