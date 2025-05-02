package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// UpdatePoolForSwap takes a pool, sender, and tokenIn, tokenOut amounts
// It then updates the pool's balances to the new reserve amounts, and
// sends the in tokens from the sender to the pool, and the out tokens from the pool to the sender.
// Swap Fees and Weight Breaking Fees is always applied on In Amount. When SwapInGivenOut, oracleInAmount should be used for it otherwise tokenInAmount
// Weight Recovery Bonus is always applied on Out Amount. When SwapInGivenOut, tokenOutAmount should be used for it otherwise oracleOutAmount
func (k Keeper) UpdatePoolForSwap(
	ctx sdk.Context,
	pool types.Pool,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
	swapFee osmomath.BigDec,
	slippageAmount osmomath.BigDec,
	oracleInAmount sdkmath.Int,
	oracleOutAmount sdkmath.Int,
	weightBalanceBonus osmomath.BigDec,
	takerFees osmomath.BigDec,
	givenOut bool,
) error {
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	// send tokensIn from sender to pool
	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err := k.bankKeeper.SendCoins(ctx, sender, poolAddr, tokensIn)
	if err != nil {
		return err
	}
	err = k.AddToPoolBalanceAndUpdateLiquidity(ctx, &pool, sdkmath.ZeroInt(), tokensIn)
	if err != nil {
		return err
	}

	// send tokensOut from pool to sender
	err = k.bankKeeper.SendCoins(ctx, poolAddr, recipient, tokensOut)
	if err != nil {
		return err
	}

	err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, sdkmath.ZeroInt(), tokensOut)
	if err != nil {
		return err
	}

	// Taker fees
	takerFeesInCoins := sdk.Coins{}
	if takerFees.IsPositive() {
		if givenOut {
			takeFeesFrom := sdk.NewCoins(sdk.NewCoin(tokenIn.Denom, oracleInAmount))
			if !pool.PoolParams.UseOracle {
				takeFeesFrom = tokensIn
			}
			// if swapInGivenOut, use oracleIn amount to get taker fees
			takerFeesInCoins = PortionCoins(takeFeesFrom, takerFees)
		} else {
			takerFeesInCoins = PortionCoins(tokensIn, takerFees)
		}
	}

	// send taker fee to taker collection address
	if takerFeesInCoins.IsAllPositive() {
		protocolAddress, err := sdk.AccAddressFromBech32(k.parameterKeeper.GetParams(ctx).TakerFeeCollectionAddress)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoins(ctx, poolAddr, protocolAddress, takerFeesInCoins)
		if err != nil {
			return err
		}

		err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, sdkmath.ZeroInt(), takerFeesInCoins)
		if err != nil {
			return err
		}
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeInCoins := sdk.Coins{}
	if swapFee.IsPositive() {
		if givenOut {
			takeFeesFrom := sdk.NewCoins(sdk.NewCoin(tokenIn.Denom, oracleInAmount))
			if !pool.PoolParams.UseOracle {
				takeFeesFrom = tokensIn
			}
			// if swapInGivenOut, use oracleIn amount to get swap fees
			swapFeeInCoins = PortionCoins(takeFeesFrom, swapFee)
		} else {
			swapFeeInCoins = PortionCoins(tokensIn, swapFee)
		}
	}

	// send swap fee to rebalance treasury
	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
		if err != nil {
			return err
		}

		err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, sdkmath.ZeroInt(), swapFeeInCoins)
		if err != nil {
			return err
		}

		err = k.OnCollectFee(ctx, pool, swapFeeInCoins)
		if err != nil {
			return err
		}
	}

	// init weightRecoveryFeeAmount to zero
	weightRecoveryFeeAmount := sdkmath.ZeroInt()

	// send half (weight breaking fee portion) of weight breaking fee to rebalance treasury
	if pool.PoolParams.UseOracle && weightBalanceBonus.IsNegative() {
		params := k.GetParams(ctx)
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		// we are multiplying here by params.WeightBreakingFeePortion as we didn't multiply in pool.SwapIn/OutGiveOut/In for weight breaking fee
		weightRecoveryFee := weightBalanceBonus.Abs().Mul(params.GetBigDecWeightBreakingFeePortion())

		if givenOut {
			weightRecoveryFeeAmount = osmomath.BigDecFromSDKInt(oracleInAmount).Mul(weightRecoveryFee).Dec().RoundInt()
		} else {
			weightRecoveryFeeAmount = osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(weightRecoveryFee).Dec().RoundInt()
		}

		if weightRecoveryFeeAmount.IsPositive() {
			// send weight recovery fee to rebalance treasury if weight recovery fee amount is positive¬
			netWeightBreakingFeeCoins := sdk.Coins{sdk.NewCoin(tokenIn.Denom, weightRecoveryFeeAmount)}

			err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, netWeightBreakingFeeCoins)
			if err != nil {
				return err
			}

			err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, sdkmath.ZeroInt(), netWeightBreakingFeeCoins)
			if err != nil {
				return err
			}

			// Track amount in pool
			weightRecoveryFeeAmountForPool := sdkmath.ZeroInt()
			weightRecoveryFeeForPool := weightBalanceBonus.Abs().Mul(osmomath.OneBigDec().Sub(params.GetBigDecWeightBreakingFeePortion()))
			if givenOut {
				weightRecoveryFeeAmountForPool = osmomath.BigDecFromSDKInt(oracleInAmount).Mul(weightRecoveryFeeForPool).Dec().RoundInt()
			} else {
				weightRecoveryFeeAmountForPool = osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(weightRecoveryFeeForPool).Dec().RoundInt()
			}
			k.TrackWeightBreakingSlippage(ctx, pool.PoolId, sdk.NewCoin(tokenIn.Denom, weightRecoveryFeeAmountForPool))
		}

	}

	bonusTokenAmount := sdkmath.ZeroInt()
	// calculate bonus token amount if weightBalanceBonus is positive
	if pool.PoolParams.UseOracle && weightBalanceBonus.IsPositive() {
		// get treasury balance
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

		// bonus token amount is the tokenOut amount times weightBalanceBonus
		if givenOut {
			bonusTokenAmount = osmomath.BigDecFromSDKInt(tokenOut.Amount).Mul(weightBalanceBonus).Dec().TruncateInt()
		} else {
			bonusTokenAmount = osmomath.BigDecFromSDKInt(oracleOutAmount).Mul(weightBalanceBonus).Dec().TruncateInt()
		}

		// if treasury balance is less than bonusTokenAmount, set bonusTokenAmount to treasury balance
		if treasuryTokenAmount.LT(bonusTokenAmount) {
			bonusTokenAmount = treasuryTokenAmount
		}

		// send bonusTokenAmount from pool addr to recipient addr, we are shortcutting the rebalance treasury address to optimize gas
		if bonusTokenAmount.IsPositive() {
			bonusToken := sdk.NewCoin(tokenOut.Denom, bonusTokenAmount)
			err = k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, recipient, sdk.Coins{bonusToken})
			if err != nil {
				return err
			}
		}
	}

	k.SetPool(ctx, pool)

	// convert the fees into USD
	swapFeeValueInUSD := k.CalculateCoinsUSDValue(ctx, swapFeeInCoins).String()
	slippageAmountInUSD := k.CalculateUSDValue(ctx, tokenIn.Denom, slippageAmount.Dec().TruncateInt()).String()
	weightRecoveryFeeAmountInUSD := k.CalculateUSDValue(ctx, tokenIn.Denom, weightRecoveryFeeAmount).String()
	bonusTokenAmountInUSD := k.CalculateUSDValue(ctx, tokenOut.Denom, bonusTokenAmount).String()
	takerFeesAmountInUSD := k.CalculateCoinsUSDValue(ctx, takerFeesInCoins).String()

	// emit swap fees event
	types.EmitSwapFeesCollectedEvent(ctx, swapFeeValueInUSD, slippageAmountInUSD, weightRecoveryFeeAmountInUSD, bonusTokenAmountInUSD, takerFeesAmountInUSD)

	// emit swap event
	types.EmitSwapEvent(ctx, sender, recipient, pool.GetPoolId(), tokensIn, tokensOut)
	if k.hooks != nil {
		err = k.hooks.AfterSwap(ctx, sender, pool, tokensIn, tokensOut)
		if err != nil {
			return err
		}
	}

	return nil
}
