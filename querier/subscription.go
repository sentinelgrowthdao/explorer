package querier

import (
	"context"
	"log"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	"google.golang.org/grpc/metadata"
)

func (q *Querier) QuerySubscription(id uint64, height int64) (subscriptiontypes.Subscription, error) {
	now := time.Now()
	defer func() {
		log.Println("QuerySubscription", height, id, time.Since(now))
	}()

	res, err := q.subscription.QuerySubscription(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&subscriptiontypes.QuerySubscriptionRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	var item subscriptiontypes.Subscription
	if err = q.InterfaceRegistry.UnpackAny(res.Subscription, &item); err != nil {
		return nil, err
	}

	return item, nil
}

func (q *Querier) QueryAllocation(id uint64, addr sdk.AccAddress, height int64) (*subscriptiontypes.Allocation, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryAllocation", height, id, addr.String(), time.Since(now))
	}()

	res, err := q.subscription.QueryAllocation(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&subscriptiontypes.QueryAllocationRequest{
			Id:      id,
			Address: addr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &res.Allocation, nil
}
