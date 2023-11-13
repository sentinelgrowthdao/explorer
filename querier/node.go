package querier

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	vpntypes "github.com/sentinel-official/hub/x/vpn/types"
	"google.golang.org/grpc/metadata"

	"github.com/sentinel-official/explorer/types"
)

func (q *Querier) ABCIQueryNode(addr hubtypes.NodeAddress, height int64) (*nodetypes.Node, error) {
	now := time.Now()
	defer func() {
		log.Println("ABCIQueryNode", height, addr.String(), time.Since(now))
	}()

	value, err := q.queryKey(
		vpntypes.ModuleName,
		append(
			[]byte(nodetypes.ModuleName+"/"),
			nodetypes.NodeKey(addr)...,
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

func (q *Querier) GRPCQueryNode(addr hubtypes.NodeAddress, height int64) (*nodetypes.Node, error) {
	now := time.Now()
	defer func() {
		log.Println("GRPCQueryNode", height, addr.String(), time.Since(now))
	}()

	qsc := nodetypes.NewQueryServiceClient(q)
	res, err := qsc.QueryNode(
		metadata.AppendToOutgoingContext(
			context.TODO(),
			grpctypes.GRPCBlockHeightHeader,
			strconv.FormatInt(height, 10),
		),
		&nodetypes.QueryNodeRequest{
			Address: addr.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &res.Node, nil
}
