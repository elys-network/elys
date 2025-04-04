package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimation(goCtx context.Context, req *types.QuerySwapEstimationRequest) (*types.QuerySwapEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	discountBigDec := osmomath.ZeroBigDec()
	if !req.Discount.IsNil() {
		discountBigDec = osmomath.BigDecFromDec(req.Discount)
	}
	spotPrice, _, tokenOut, swapFee, discount, availableLiquidity, slippage, weightBonus, err := k.CalcInRouteSpotPrice(ctx, req.TokenIn, req.Routes, discountBigDec, osmomath.ZeroBigDec())
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationResponse{
		SpotPrice:          spotPrice.Dec(),
		TokenOut:           tokenOut,
		SwapFee:            swapFee.Dec(),
		Discount:           discount.Dec(),
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage.Dec(),
		WeightBalanceRatio: weightBonus.Dec(),
	}, nil
}
