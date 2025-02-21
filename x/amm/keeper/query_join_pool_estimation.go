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
	tokensIn, sharesOut, slippage, weightBalanceBonus, swapFee, err := k.JoinPoolEst(ctx, req.PoolId, req.AmountsIn)
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetPoolShareDenom(req.PoolId)
	return &types.QueryJoinPoolEstimationResponse{
		ShareAmountOut:     sdk.NewCoin(shareDenom, sharesOut),
		AmountsIn:          tokensIn,
		Slippage:           slippage,
		WeightBalanceRatio: weightBalanceBonus,
		SwapFee:            swapFee,
	}, nil
}

func (k Keeper) JoinPoolEst(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, slippage math.LegacyDec, weightBalanceBonus math.LegacyDec, swapFee math.LegacyDec, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
			}
		}

		params := k.GetParams(ctx)
		takerFees := k.parameterKeeper.GetParams(ctx).TakerFees
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn, params, takerFees)
		if err != nil {
			return nil, math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}

		return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, err
	}

	params := k.GetParams(ctx)
	takerFees := k.parameterKeeper.GetParams(ctx).TakerFees
	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	cacheCtx, _ := ctx.CacheContext()
	tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs, params, takerFees)
	if err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, err
}
