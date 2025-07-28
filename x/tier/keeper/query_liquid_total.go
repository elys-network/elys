package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LiquidTotal(goCtx context.Context, req *types.QueryLiquidTotalRequest) (*types.QueryLiquidTotalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(req.User)
	total := k.RetrieveLiquidAssetsTotal(ctx, sender)

	return &types.QueryLiquidTotalResponse{
		Total: total.Dec(),
	}, nil
}
