package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"

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

	// Even when multiple routes, taker Fee per route is params.TakerFee/TotalRoutes, so net taker fee will be params.TakerFee
	takerFees := k.parameterKeeper.GetParams(ctx).TakerFees

	spotPrice, _, tokenOut, swapFee, discount, availableLiquidity, slippage, weightBonus, err := k.CalcInRouteSpotPrice(ctx, req.TokenIn, req.Routes, req.Discount, sdkmath.LegacyZeroDec())
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationResponse{
		SpotPrice:          spotPrice,
		TokenOut:           tokenOut,
		SwapFee:            swapFee,
		TakerFee:           takerFees,
		Discount:           discount,
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage,
		WeightBalanceRatio: weightBonus,
	}, nil
}
