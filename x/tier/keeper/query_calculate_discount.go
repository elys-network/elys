package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CalculateDiscount(goCtx context.Context, req *types.QueryCalculateDiscountRequest) (*types.QueryCalculateDiscountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	user, err := sdk.AccAddressFromBech32(req.User)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	portfolio, tier, discount := k.GetMembershipTier(ctx, user)

	return &types.QueryCalculateDiscountResponse{
		Discount:  discount.String(),
		Portfolio: portfolio.String(),
		Tier:      tier,
	}, nil
}
