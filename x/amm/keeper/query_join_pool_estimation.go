package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
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
		Slippage:           slippage.String(),
		WeightBalanceRatio: weightBalanceBonus.String(),
	}, nil
}

func (k Keeper) JoinPoolEst(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, slippage elystypes.Dec34, weightBalanceBonus elystypes.Dec34, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, math.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
			}
		}

		params := k.GetParams(ctx)
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		tokensJoined, sharesOut, slippage, weightBalanceBonus, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn, params)
		if err != nil {
			return nil, math.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		return tokensJoined, sharesOut, slippage, weightBalanceBonus, err
	}

	params := k.GetParams(ctx)
	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	cacheCtx, _ := ctx.CacheContext()
	tokensJoined, sharesOut, slippage, weightBalanceBonus, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs, params)
	if err != nil {
		return nil, math.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	return tokensJoined, sharesOut, slippage, weightBalanceBonus, err
}
