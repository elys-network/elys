package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RewardProgram(goCtx context.Context, req *types.QueryRewardProgramRequest) (*types.QueryRewardProgramResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	rewardProgram := k.GetRewardProgram(ctx, address)
	return &types.QueryRewardProgramResponse{
		Amount:  rewardProgram.Amount,
		Claimed: rewardProgram.Claimed,
	}, nil
}

func (k Keeper) TotalRewardProgramClaimed(goCtx context.Context, req *types.QueryTotalRewardProgramClaimedRequest) (*types.QueryTotalRewardProgramClaimedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	total := k.GetTotalRewardProgramClaimed(ctx)
	return &types.QueryTotalRewardProgramClaimedResponse{
		TotalEdenClaimed: total.TotalEdenClaimed,
	}, nil
}
