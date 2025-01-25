package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RewardsTotal(goCtx context.Context, req *types.QueryRewardsTotalRequest) (*types.QueryRewardsTotalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(req.User)
	total := k.RetrieveRewardsTotal(ctx, sender)

	return &types.QueryRewardsTotalResponse{
		Total: total.String(),
	}, nil
}
