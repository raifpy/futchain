package keeper

import (
	"context"

	"github.com/raifpy/futchainevm/x/futchain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) League(ctx context.Context, req *types.QueryLeagueRequest) (*types.QueryLeagueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	leag, err := q.k.GetLeague(ctx, int(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLeagueResponse{
		Id:        int64(leag.ID),
		Name:      leag.Name,
		GroupName: leag.GroupName,
	}, nil
}
