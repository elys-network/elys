package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimation(goCtx context.Context, req *types.QuerySwapEstimationRequest) (*types.QuerySwapEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	spotPrice, tokenOut, err := k.CalcInRouteSpotPrice(ctx, req.TokenIn, req.Routes)
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationResponse{
		SpotPrice: spotPrice,
		TokenOut:  tokenOut,
	}, nil
}
