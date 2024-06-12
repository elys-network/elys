package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/membershiptier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CalculateDiscount(goCtx context.Context, req *types.QueryCalculateDiscountRequest) (*types.QueryCalculateDiscountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	portfolio, tier, discount := k.GetMembershipTier(ctx, req.User)

	return &types.QueryCalculateDiscountResponse{
		Discount:  strconv.FormatUint(discount, 10),
		Portfolio: portfolio.String(),
		Tier:      tier,
	}, nil
}
