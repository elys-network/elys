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
	tokensIn, sharesOut, slippage, weightBalanceBonus, swapFee, takerFees, err := k.JoinPoolEst(ctx, req.PoolId, req.AmountsIn)
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
		TakerFee:           takerFees,
	}, nil
}

func (k Keeper) JoinPoolEst(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, slippage math.LegacyDec, weightBalanceBonus math.LegacyDec, swapFee math.LegacyDec, takerFeesFinal math.LegacyDec, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
			}
		}

		params := k.GetParams(ctx)
		takerFees := k.parameterKeeper.GetParams(ctx).TakerFees
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn, params, takerFees)
		if err != nil {
			return nil, math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}

		return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, err
	}

	params := k.GetParams(ctx)
	takerFees := k.parameterKeeper.GetParams(ctx).TakerFees
	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	cacheCtx, _ := ctx.CacheContext()
	tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, _, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs, params, takerFees)
	if err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	var otherAsset types.PoolAsset
	// Check treasury and update weightBalance
	if weightBalanceBonus.IsPositive() && tokensJoined.Len() == 1 {
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		for _, asset := range pool.PoolAssets {
			if asset.Token.Denom == tokensJoined[0].Denom {
				continue
			}
			otherAsset = asset
		}
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, otherAsset.Token.Denom).Amount

		bonusTokenAmount := tokensJoined[0].Amount.ToLegacyDec().Mul(weightBalanceBonus).TruncateInt()

		if treasuryTokenAmount.LT(bonusTokenAmount) {
			weightBalanceBonus = treasuryTokenAmount.ToLegacyDec().Quo(tokensJoined[0].Amount.ToLegacyDec())
		}
	}

	return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, err
}
