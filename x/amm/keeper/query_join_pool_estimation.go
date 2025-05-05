package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) JoinPoolEstimation(goCtx context.Context, req *types.QueryJoinPoolEstimationRequest) (*types.QueryJoinPoolEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	tokensIn, sharesOut, slippage, weightBalanceBonus, swapFee, takerFees, weightRewardAmount, err := k.JoinPoolEst(ctx, req.PoolId, req.AmountsIn)
	if err != nil {
		return nil, err
	}

	takerFeesLegacyDec := math.LegacyZeroDec()
	if !takerFees.IsNil() {
		takerFeesLegacyDec = takerFees.Dec()
	}

	shareDenom := types.GetPoolShareDenom(req.PoolId)
	return &types.QueryJoinPoolEstimationResponse{
		ShareAmountOut:            sdk.NewCoin(shareDenom, sharesOut),
		AmountsIn:                 tokensIn,
		Slippage:                  slippage.Dec(),
		WeightBalanceRatio:        weightBalanceBonus.Dec(),
		SwapFee:                   swapFee.Dec(),
		TakerFee:                  takerFeesLegacyDec,
		WeightBalanceRewardAmount: weightRewardAmount,
	}, nil
}

func (k Keeper) JoinPoolEst(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, slippage osmomath.BigDec, weightBalanceBonus osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, weightRewardAmount sdk.Coin, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, err
			}
		}

		params := k.GetParams(ctx)
		takerFees := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()
		snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)
		cacheCtx, _ := ctx.CacheContext()
		tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn, params, takerFees)
		if err != nil {
			return nil, math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, err
		}

		return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, sdk.Coin{}, nil
	}

	params := k.GetParams(ctx)
	takerFees := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()
	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)
	cacheCtx, _ := ctx.CacheContext()
	tokensJoined := sdk.Coins{}
	tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, err = pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs, params, takerFees)
	if err != nil {
		return nil, math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, err
	}

	var otherAsset types.PoolAsset
	bonusTokenAmount := math.ZeroInt()
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
		// bonusTokenAmount = osmomath.BigDecFromSDKInt(tokensJoined[0].Amount).Mul(weightBalanceBonus).Dec().TruncateInt()
		// ensure token prices for in/out tokens set properly
		inTokenPrice := k.oracleKeeper.GetDenomPrice(ctx, tokensJoined[0].Denom)
		if inTokenPrice.IsZero() {
			return nil, math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, fmt.Errorf("price for inToken not set: %s", tokensJoined[0].Denom)
		}
		outTokenPrice := k.oracleKeeper.GetDenomPrice(ctx, otherAsset.Token.Denom)
		if outTokenPrice.IsZero() {
			return nil, math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, fmt.Errorf("price for outToken not set: %s", otherAsset.Token.Denom)
		}
		bonusTokenAmount = (osmomath.BigDecFromSDKInt(tokensJoined[0].Amount).Mul(weightBalanceBonus).Mul(inTokenPrice).Quo(outTokenPrice)).Dec().TruncateInt()

		if treasuryTokenAmount.LT(bonusTokenAmount) {
			bonusTokenAmount = treasuryTokenAmount
		}
	}
	rewards := sdk.Coin{}
	if otherAsset.Token.Denom != "" && bonusTokenAmount.IsPositive() {
		rewards = sdk.NewCoin(otherAsset.Token.Denom, bonusTokenAmount)
	}

	return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, rewards, nil
}
