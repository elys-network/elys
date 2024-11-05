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
	tokensIn, sharesOut, slippage, weightBalanceBonus, err := k.JoinPoolEst(ctx, req.PoolId, req.AmountsIn)
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetPoolShareDenom(req.PoolId)
	return &types.QueryJoinPoolEstimationResponse{
		ShareAmountOut:     sdk.NewCoin(shareDenom, sharesOut),
		AmountsIn:          tokensIn,
		Slippage:           slippage,
		WeightBalanceRatio: weightBalanceBonus,
	}, nil
}

func (k Keeper) JoinPoolEst(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, slippage sdk.Dec, weightBalanceBonus sdk.Dec, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, sdk.ZeroDec(), sdk.ZeroDec(), err
			}
		}

		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		sharesOut, slippage, weightBalanceBonus, err = pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn)
		if err != nil {
			return nil, sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		return tokensIn, sharesOut, slippage, weightBalanceBonus, err
	}

	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	cacheCtx, _ := ctx.CacheContext()
	sharesOut, slippage, weightBalanceBonus, err = pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs)
	if err != nil {
		return nil, sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}

	return tokenInMaxs, sharesOut, slippage, weightBalanceBonus, err
}
