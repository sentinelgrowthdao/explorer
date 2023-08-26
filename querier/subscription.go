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

func (q *Querier) QuerySubscription(id uint64, height int64) (*subscriptiontypes.Subscription, error) {
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

	return &res.Subscription, nil
}

func (q *Querier) QueryQuota(id uint64, addr sdk.AccAddress, height int64) (*subscriptiontypes.Quota, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryQuota", height, id, addr.String(), time.Since(now))
	}()

	res, err := q.subscription.QueryQuota(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&subscriptiontypes.QueryQuotaRequest{
			Id:      id,
			Address: addr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &res.Quota, nil
}
