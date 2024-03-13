package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) JoinPoolEstimation(goCtx context.Context, req *types.QueryJoinPoolEstimationRequest) (*types.QueryJoinPoolEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	tokensIn, sharesOut, err := k.joinPoolEstimation(ctx, req.PoolId, req.AmountsIn)
	if err != nil {
		return nil, err
	}

	return &types.QueryJoinPoolEstimationResponse{
		ShareAmountOut: sharesOut,
		AmountsIn:      tokensIn,
	}, nil
}

func (k Keeper) joinPoolEstimation(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, sdk.ZeroInt(), types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, err
			}
		}

		sharesOut, err = pool.JoinPool(ctx, k.oracleKeeper, k.accountedPoolKeeper, tokensIn)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		return tokensIn, sharesOut, err
	}

	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	sharesOut, err = pool.JoinPool(ctx, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	return tokenInMaxs, sharesOut, err
}
