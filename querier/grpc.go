package querier

import (
	"context"
	"fmt"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	gogogrpc "github.com/gogo/protobuf/grpc"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/metadata"
)

var (
	_ gogogrpc.ClientConn = &Querier{}
)

func (q *Querier) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	data, err := encoding.GetCodec(proto.Name).Marshal(args)
	if err != nil {
		return err
	}

	req := &abcitypes.RequestQuery{
		Path: method,
		Data: data,
	}

	md, _ := metadata.FromOutgoingContext(ctx)
	if heights := md.Get(grpctypes.GRPCBlockHeightHeader); len(heights) > 0 {
		req.Height, err = strconv.ParseInt(heights[0], 10, 64)
		if err != nil {
			return err
		}
		if req.Height < 0 {
			return fmt.Errorf("request height cannot be negative")
		}
	}

	res, err := q.queryABCI(req)
	if err != nil {
		return err
	}

	if err = encoding.GetCodec(proto.Name).Unmarshal(res.Value, reply); err != nil {
		return err
	}

	md = metadata.Pairs(grpctypes.GRPCBlockHeightHeader, strconv.FormatInt(res.Height, 10))
	for _, opt := range opts {
		header, ok := opt.(grpc.HeaderCallOption)
		if !ok {
			continue
		}

		*header.HeaderAddr = md
	}

	if q.InterfaceRegistry != nil {
		return codectypes.UnpackInterfaces(reply, q.InterfaceRegistry)
	}

	return nil
}

func (_ *Querier) NewStream(_ context.Context, _ *grpc.StreamDesc, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, fmt.Errorf("not supported")
}
