package keeper

import (
	"context"

	"github.com/raifpy/futchain/x/futchain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Team(ctx context.Context, req *types.QueryTeamRequest) (*types.QueryTeamResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	team, err := q.k.GetTeam(ctx, int(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTeamResponse{
		Id:   int64(team.ID),
		Name: team.Name,
	}, nil
}
