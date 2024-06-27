package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ShowCommitments(goCtx context.Context, req *types.QueryShowCommitmentsRequest) (*types.QueryShowCommitmentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	val := k.GetCommitments(ctx, req.Creator)
	return &types.QueryShowCommitmentsResponse{Commitments: val}, nil
}

func (k Keeper) NumberOfCommitments(goCtx context.Context, req *types.QueryNumberOfCommitmentsRequest) (*types.QueryNumberOfCommitmentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryNumberOfCommitmentsResponse{Number: k.TotalNumberOfCommitments(ctx)}, nil
}

func (k Keeper) CommittedTokensLocked(goCtx context.Context, req *types.QueryCommittedTokensLockedRequest) (*types.QueryCommittedTokensLockedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	commitments := k.GetCommitments(ctx, req.Address)
	totalLocked, totalCommitted := commitments.CommittedTokensLocked(ctx)
	return &types.QueryCommittedTokensLockedResponse{
		Address:         req.Address,
		LockedCommitted: totalLocked,
		TotalCommitted:  totalCommitted,
	}, nil
}
