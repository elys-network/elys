package keeper

import (
	"context"

	"cosmossdk.io/math"
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

	// Start with the initial token input
	tokensIn := sdk.Coins{req.TokenIn}

	for _, route := range req.Routes {
		poolId := route.PoolId
		tokenOutDenom := route.TokenOutDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return nil, status.Error(codes.NotFound, "pool not found")
		}

		// Estimate swap
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		swapResult, err := k.CalcOutAmtGivenIn(ctx, pool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, sdk.ZeroDec())

		if err != nil {
			return nil, err
		}

		if swapResult.IsZero() {
			return nil, status.Error(codes.InvalidArgument, "amount too low")
		}

		// Use the current swap result as the input for the next iteration
		tokensIn = sdk.Coins{swapResult}
	}

	// Calculate the spot price given the initial token in and the final token out
	spotPrice := math.LegacyNewDecFromInt(tokensIn[0].Amount).Quo(math.LegacyNewDecFromInt(req.TokenIn.Amount))

	return &types.QuerySwapEstimationResponse{
		SpotPrice: spotPrice,
		TokenOut:  tokensIn[0],
	}, nil
}
