package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
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
	swapFee sdkmath.LegacyDec,
	oracleInAmount sdkmath.Int,
	oracleOutAmount sdkmath.Int,
	weightBalanceBonus sdkmath.LegacyDec,
	givenOut bool,
) error {
	//if givenOut && oracleInAmount.IsZero() {
	//	return fmt.Errorf("oracleInAmount cannot be zero for SwapInGivenOut")
	//}
	//if !givenOut && tokenIn.IsZero() {
	//	return fmt.Errorf("tokenIn cannot be zero for SwapInGivenOut")
	//}
	//if givenOut && tokenOut.IsZero() {
	//	return fmt.Errorf("tokenOut cannot be zero for SwapOutGivenOut")
	//}
	//if !givenOut && oracleOutAmount.IsZero() {
	//	return fmt.Errorf("oracleOutAmount cannot be zero for SwapOutGivenOut")
	//}
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	k.SetPool(ctx, pool)

	// send tokensIn from sender to pool
	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err := k.bankKeeper.SendCoins(ctx, sender, poolAddr, tokensIn)
	if err != nil {
		return err
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeInCoins := sdk.Coins{}
	if !weightBalanceBonus.IsPositive() {
		if givenOut {
			// if swapInGivenOut, use oracleIn amount to get swap fees
			swapFeeInCoins = PortionCoins(sdk.Coins{sdk.NewCoin(tokenIn.Denom, oracleInAmount)}, swapFee)
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
		err = k.OnCollectFee(ctx, pool, swapFeeInCoins)
		if err != nil {
			return err
		}
	}

	// init weightRecoveryFeeAmount to zero
	weightRecoveryFeeAmount := sdkmath.ZeroInt()

	// send half (weight breaking fee portion) of weight breaking fee to rebalance treasury
	if !weightBalanceBonus.IsPositive() {
		params := k.GetParams(ctx)
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		weightRecoveryFee := weightBalanceBonus.Abs().Mul(params.WeightBreakingFeePortion)

		// weight recovery fee must be positive
		if weightRecoveryFee.IsPositive() {
			if givenOut {
				weightRecoveryFeeAmount = oracleInAmount.ToLegacyDec().Mul(weightRecoveryFee).RoundInt()
			} else {
				weightRecoveryFeeAmount = tokenIn.Amount.ToLegacyDec().Mul(weightRecoveryFee).RoundInt()
			}

			// send weight recovery fee to rebalance treasury if weight recovery fee amount is positive
			if weightRecoveryFeeAmount.IsPositive() {
				err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, sdk.Coins{sdk.NewCoin(tokenIn.Denom, weightRecoveryFeeAmount)})
				if err != nil {
					return err
				}
			}
		}
	}

	// Send coins to recipient
	err = k.bankKeeper.SendCoins(ctx, poolAddr, recipient, sdk.Coins{tokenOut})
	if err != nil {
		return err
	}

	// calculate bonus token amount if weightBalanceBonus is positive
	if weightBalanceBonus.IsPositive() {
		// get treasury balance
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

		bonusTokenAmount := sdkmath.ZeroInt()
		// bonus token amount is the tokenOut amount times weightBalanceBonus
		if givenOut {
			bonusTokenAmount = tokenOut.Amount.ToLegacyDec().Mul(weightBalanceBonus).TruncateInt()
		} else {
			bonusTokenAmount = oracleOutAmount.ToLegacyDec().Mul(weightBalanceBonus).TruncateInt()
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

	// emit swap event
	types.EmitSwapEvent(ctx, sender, recipient, pool.GetPoolId(), tokensIn, tokensOut)
	if k.hooks != nil {
		err = k.hooks.AfterSwap(ctx, sender, pool, tokensIn, tokensOut)
		if err != nil {
			return err
		}
	}

	// record tokenIn amount as total liquidity increase
	err = k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	if err != nil {
		return err
	}

	// record tokenOut amount as total liquidity decrease
	err = k.RecordTotalLiquidityDecrease(ctx, tokensOut)
	if err != nil {
		return err
	}

	// record swap fee as total liquidity decrease
	err = k.RecordTotalLiquidityDecrease(ctx, swapFeeInCoins)
	if err != nil {
		return err
	}

	// record weight recovery fee amount as total liquidity decrease
	weightRecoveryToken := sdk.NewCoin(tokenOut.Denom, weightRecoveryFeeAmount)
	err = k.RecordTotalLiquidityDecrease(ctx, sdk.Coins{weightRecoveryToken})
	if err != nil {
		return err
	}

	return nil
}
