package querier

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	vpntypes "github.com/sentinel-official/hub/x/vpn/types"
	"google.golang.org/grpc/metadata"

	"github.com/sentinel-official/explorer/types"
)

func (q *Querier) ABCIQuerySession(id uint64, height int64) (*sessiontypes.Session, error) {
	now := time.Now()
	defer func() {
		log.Println("ABCIQuerySession", height, id, time.Since(now))
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
	if err := types.EncCfg.Marshaler.UnmarshalBinaryBare(value, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (q *Querier) GRPCQuerySession(id uint64, height int64) (*sessiontypes.Session, error) {
	now := time.Now()
	defer func() {
		log.Println("GRPCQuerySession", height, id, time.Since(now))
	}()

	qsc := sessiontypes.NewQueryServiceClient(q)
	res, err := qsc.QuerySession(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&sessiontypes.QuerySessionRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &res.Session, nil
}
