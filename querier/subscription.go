package querier

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	vpntypes "github.com/sentinel-official/hub/x/vpn/types"
	"google.golang.org/grpc/metadata"

	"github.com/sentinel-official/explorer/types"
)

func (q *Querier) ABCIQuerySubscription(id uint64, height int64) (*subscriptiontypes.Subscription, error) {
	now := time.Now()
	defer func() {
		log.Println("ABCIQuerySubscription", height, id, time.Since(now))
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

func (q *Querier) ABCIQueryQuota(id uint64, addr sdk.AccAddress, height int64) (*subscriptiontypes.Quota, error) {
	now := time.Now()
	defer func() {
		log.Println("ABCIQueryQuota", height, id, addr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(subscriptiontypes.ModuleName+"/"),
			subscriptiontypes.QuotaKey(id, addr)...,
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

func (q *Querier) GRPCQuerySubscription(id uint64, height int64) (*subscriptiontypes.Subscription, error) {
	now := time.Now()
	defer func() {
		log.Println("GRPCQuerySubscription", height, id, time.Since(now))
	}()

	qsc := subscriptiontypes.NewQueryServiceClient(q)
	res, err := qsc.QuerySubscription(
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

func (q *Querier) GRPCQueryQuota(id uint64, addr sdk.AccAddress, height int64) (*subscriptiontypes.Quota, error) {
	now := time.Now()
	defer func() {
		log.Println("GRPCQueryQuota", height, id, addr.String(), time.Since(now))
	}()

	qsc := subscriptiontypes.NewQueryServiceClient(q)
	res, err := qsc.QueryQuota(
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
