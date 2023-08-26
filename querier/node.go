package querier

import (
	"context"
	"log"
	"strconv"
	"time"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	"google.golang.org/grpc/metadata"
)

func (q *Querier) QueryNode(addr hubtypes.NodeAddress, height int64) (*nodetypes.Node, error) {
	now := time.Now()
	defer func() {
		log.Println("QueryNode", height, addr.String(), time.Since(now))
	}()

	res, err := q.node.QueryNode(
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
