package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimationExactAmountOut(goCtx context.Context, req *types.QuerySwapEstimationExactAmountOutRequest) (*types.QuerySwapEstimationExactAmountOutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	discountBigDec := osmomath.ZeroBigDec()
	if !req.Discount.IsNil() {
		discountBigDec = osmomath.BigDecFromDec(req.Discount)
	}

	spotPrice, _, tokenIn, swapFee, discount, availableLiquidity, slippage, weightBonus, err := k.CalcOutRouteSpotPrice(ctx, req.TokenOut, req.Routes, discountBigDec, osmomath.ZeroBigDec())
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationExactAmountOutResponse{
		SpotPrice:          spotPrice.Dec(),
		TokenIn:            tokenIn,
		SwapFee:            swapFee.Dec(),
		Discount:           discount.Dec(),
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage.Dec(),
		WeightBalanceRatio: weightBonus.Dec(),
	}, nil
}
