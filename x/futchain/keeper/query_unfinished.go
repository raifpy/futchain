package keeper

import (
	"context"

	"github.com/raifpy/futchain/x/futchain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) UnfinishedMatches(ctx context.Context, req *types.QueryUnfinishedMatchesRequest) (*types.QueryUnfinishedMatchesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ids, err := q.k.ListUnfinishedMatches(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	int64s := make([]int64, len(ids))
	for i, id := range ids {
		int64s[i] = int64(id)
	}

	return &types.QueryUnfinishedMatchesResponse{
		Ids: int64s,
	}, nil

}
