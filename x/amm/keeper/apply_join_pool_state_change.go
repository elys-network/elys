package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) ApplyJoinPoolStateChange(
	ctx sdk.Context,
	pool types.Pool,
	joiner sdk.AccAddress,
	numShares math.Int,
	joinCoins sdk.Coins,
	weightBalanceBonus osmomath.BigDec,
	takerFees osmomath.BigDec,
	swapFee osmomath.BigDec,
	slippageCoins sdk.Coins,
) error {
	if err := k.bankKeeper.SendCoins(ctx, joiner, sdk.MustAccAddressFromBech32(pool.GetAddress()), joinCoins); err != nil {
		return err
	}

	if err := k.MintPoolShareToAccount(ctx, pool, joiner, numShares); err != nil {
		return err
	}

	k.SetPool(ctx, pool)

	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())

	swapFeeInCoins := sdk.Coins{}
	if swapFee.IsPositive() {
		swapFeeInCoins = PortionCoins(joinCoins, swapFee)
	}

	// send swap fee to rebalance treasury
	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err := k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
		if err != nil {
			return err
		}

		err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, math.ZeroInt(), swapFeeInCoins)
		if err != nil {
			return err
		}

		err = k.OnCollectFee(ctx, pool, swapFeeInCoins)
		if err != nil {
			return err
		}
	}

	weightRecoveryFeeAmount := math.ZeroInt()
	weightRecoveryFeeCoins := sdk.Coins{}
	// send half (weight breaking fee portion) of weight breaking fee to rebalance treasury
	if pool.PoolParams.UseOracle && weightBalanceBonus.IsNegative() {
		params := k.GetParams(ctx)
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		// we are multiplying here by params.WeightBreakingFeePortion as we didn't multiply in pool.Join/Exit for weight breaking fee
		weightRecoveryFee := weightBalanceBonus.Abs().Mul(params.GetBigDecWeightBreakingFeePortion())
		for _, coin := range joinCoins {
			weightRecoveryFeeAmount = osmomath.BigDecFromSDKInt(coin.Amount).Mul(weightRecoveryFee).Dec().RoundInt()

			if weightRecoveryFeeAmount.IsPositive() {
				// send weight recovery fee to rebalance treasury if weight recovery fee amount is positive¬
				netWeightBreakingFeeCoins := sdk.Coins{sdk.NewCoin(coin.Denom, weightRecoveryFeeAmount)}
				weightRecoveryFeeCoins = weightRecoveryFeeCoins.Add(netWeightBreakingFeeCoins...)

				err := k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, netWeightBreakingFeeCoins)
				if err != nil {
					return err
				}

				err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, math.ZeroInt(), netWeightBreakingFeeCoins)
				if err != nil {
					return err
				}

				// Track amount in pool
				weightRecoveryFeeForPool := weightBalanceBonus.Abs().Mul(osmomath.OneBigDec().Sub(params.GetBigDecWeightBreakingFeePortion()))
				k.TrackWeightBreakingSlippage(ctx, pool.PoolId, sdk.NewCoin(coin.Denom, weightRecoveryFeeForPool.Mul(osmomath.BigDecFromSDKInt(weightRecoveryFeeAmount)).Dec().TruncateInt()))
			}
		}
	}

	var weightBalanceBonusCoins sdk.Coins
	var otherAsset types.PoolAsset
	// Check treasury and update weightBalance
	if weightBalanceBonus.IsPositive() && joinCoins.Len() == 1 {
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		for _, asset := range pool.PoolAssets {
			if asset.Token.Denom == joinCoins[0].Denom {
				continue
			}
			otherAsset = asset
		}
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, otherAsset.Token.Denom).Amount

		bonusTokenAmount := osmomath.BigDecFromSDKInt(joinCoins[0].Amount).Mul(weightBalanceBonus).Dec().TruncateInt()

		if treasuryTokenAmount.LT(bonusTokenAmount) {
			bonusTokenAmount = treasuryTokenAmount
		}
		weightBalanceBonusCoins = sdk.Coins{sdk.NewCoin(otherAsset.Token.Denom, bonusTokenAmount)}

		// send bonus tokens to recipient if positive
		if weightBalanceBonusCoins.IsAllPositive() {
			if err := k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, joiner, weightBalanceBonusCoins); err != nil {
				return err
			}
		}
	}

	// Taker fees
	takerFeesInCoins := sdk.Coins{}
	if takerFees.IsPositive() {
		takerFeesInCoins = PortionCoins(joinCoins, takerFees)
	}

	// send taker fee to protocol treasury
	if takerFeesInCoins.IsAllPositive() {
		protocolAddress, err := sdk.AccAddressFromBech32(k.parameterKeeper.GetParams(ctx).TakerFeeCollectionAddress)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoins(ctx, poolAddr, protocolAddress, takerFeesInCoins)
		if err != nil {
			return err
		}

		err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, math.ZeroInt(), takerFeesInCoins)
		if err != nil {
			return err
		}
	}

	// convert the fees into USD
	swapFeeValueInUSD := k.CalculateCoinsUSDValue(ctx, swapFeeInCoins).String()
	slippageAmountInUSD := k.CalculateCoinsUSDValue(ctx, slippageCoins).String()
	weightRecoveryFeeAmountInUSD := k.CalculateCoinsUSDValue(ctx, weightRecoveryFeeCoins).String()
	bonusTokenAmountInUSD := k.CalculateCoinsUSDValue(ctx, weightBalanceBonusCoins).String()
	takerFeesAmountInUSD := k.CalculateCoinsUSDValue(ctx, takerFeesInCoins).String()

	// emit swap fees event
	types.EmitSwapFeesCollectedEvent(ctx, swapFeeValueInUSD, slippageAmountInUSD, weightRecoveryFeeAmountInUSD, bonusTokenAmountInUSD, takerFeesAmountInUSD)

	types.EmitAddLiquidityEvent(ctx, joiner, pool.GetPoolId(), joinCoins)
	if k.hooks != nil {
		err := k.hooks.AfterJoinPool(ctx, joiner, pool, joinCoins, numShares)
		if err != nil {
			return err
		}
	}
	return nil
}
