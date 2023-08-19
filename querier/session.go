package querier

import (
	"context"
	"log"
	"strconv"
	"time"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	"google.golang.org/grpc/metadata"
)

func (q *Querier) QuerySession(id uint64, height int64) (*sessiontypes.Session, error) {
	now := time.Now()
	defer func() {
		log.Println("QuerySession", height, id, time.Since(now))
	}()

	res, err := q.session.QuerySession(
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
