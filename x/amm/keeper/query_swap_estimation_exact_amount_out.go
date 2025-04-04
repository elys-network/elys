package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimationExactAmountOut(goCtx context.Context, req *types.QuerySwapEstimationExactAmountOutRequest) (*types.QuerySwapEstimationExactAmountOutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	spotPrice, _, tokenIn, swapFee, discount, availableLiquidity, slippage, weightBonus, err := k.CalcOutRouteSpotPrice(ctx, req.TokenOut, req.Routes, req.Discount, sdkmath.LegacyZeroDec())
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationExactAmountOutResponse{
		SpotPrice:          spotPrice.String(),
		TokenIn:            tokenIn,
		SwapFee:            swapFee,
		Discount:           discount,
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage.String(),
		WeightBalanceRatio: weightBonus.String(),
	}, nil
}
