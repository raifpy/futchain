package keeper

import (
	"context"

	"github.com/raifpy/futchain/x/futchain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Match(ctx context.Context, req *types.QueryMatchRequest) (*types.QueryMatchResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	match, err := q.k.GetMatch(ctx, int(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMatchResponse{
		Id:        int64(match.ID),
		LeagueId:  int64(match.LeagueID),
		Name:      match.Home.Name + " - " + match.Away.Name,
		Time:      match.Time,
		HomeId:    int64(match.Home.ID),
		AwayId:    int64(match.Away.ID),
		HomeScore: int64(match.Home.Score),
		AwayScore: int64(match.Away.Score),
		HomeName:  match.Home.Name,
		AwayName:  match.Away.Name,
		Started:   match.Status.Started,
		Finished:  match.Status.Finished,
		Cancelled: match.Status.Cancelled,
	}, nil
}
