package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) ApplyExitPoolStateChange(ctx sdk.Context, pool types.Pool, exiter sdk.AccAddress, numShares sdkmath.Int, exitCoins sdk.Coins, isLiquidation bool, weightBalanceBonus elystypes.Dec34) error {
	// Withdraw exit amount of token from commitment module to exiter's wallet.
	poolShareDenom := types.GetPoolShareDenom(pool.GetPoolId())

	// Withdraw committed LP tokens
	err := k.commitmentKeeper.UncommitTokens(ctx, exiter, poolShareDenom, numShares, isLiquidation)
	if err != nil {
		return err
	}

	if err = k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()), exiter, exitCoins); err != nil {
		return err
	}

	if err = k.BurnPoolShareFromAccount(ctx, pool, exiter, numShares); err != nil {
		return err
	}

	k.SetPool(ctx, pool)

	rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
	weightRecoveryFeeAmount := sdkmath.ZeroInt()
	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	// send half (weight breaking fee portion) of weight breaking fee to rebalance treasury
	if pool.PoolParams.UseOracle && weightBalanceBonus.IsNegative() {
		params := k.GetParams(ctx)
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		// we are multiplying here by params.WeightBreakingFeePortion as we didn't multiply in pool.Join/Exit for weight breaking fee
		weightRecoveryFee := weightBalanceBonus.Abs().MulLegacyDec(params.WeightBreakingFeePortion)

		for _, coin := range exitCoins {
			weightRecoveryFeeAmount = weightRecoveryFee.MulInt(coin.Amount).ToInt()

			if weightRecoveryFeeAmount.IsPositive() {
				// send weight recovery fee to rebalance treasury if weight recovery fee amount is positive¬
				netWeightBreakingFeeCoins := sdk.Coins{sdk.NewCoin(coin.Denom, weightRecoveryFeeAmount)}

				err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, netWeightBreakingFeeCoins)
				if err != nil {
					return err
				}

				err = k.RemoveFromPoolBalanceAndUpdateLiquidity(ctx, &pool, sdkmath.ZeroInt(), netWeightBreakingFeeCoins)
				if err != nil {
					return err
				}

				// Track amount in pool
				weightRecoveryFeeForPool := weightBalanceBonus.Abs().Mul(elystypes.OneDec34().SubLegacyDec(params.WeightBreakingFeePortion))
				k.TrackWeightBreakingSlippage(ctx, pool.PoolId, sdk.NewCoin(coin.Denom, weightRecoveryFeeForPool.MulInt(weightRecoveryFeeAmount).ToInt()))
			}
		}
	}

	if weightBalanceBonus.IsPositive() {
		// calculate treasury amounts to send as bonus
		weightBalanceBonusCoins := PortionCoins(exitCoins, weightBalanceBonus)
		for _, coin := range weightBalanceBonusCoins {
			treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, coin.Denom).Amount
			if treasuryTokenAmount.LT(coin.Amount) {
				// override coin amount by treasuryTokenAmount
				weightBalanceBonusCoins = weightBalanceBonusCoins.
					Sub(coin).                                        // remove the original coin
					Add(sdk.NewCoin(coin.Denom, treasuryTokenAmount)) // add the treasuryTokenAmount
			}
		}

		// send bonus tokens to recipient if positive
		if weightBalanceBonusCoins.IsAllPositive() {
			if err := k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, exiter, weightBalanceBonusCoins); err != nil {
				return err
			}
		}
	}

	types.EmitRemoveLiquidityEvent(ctx, exiter, pool.GetPoolId(), exitCoins)
	if k.hooks != nil {
		err = k.hooks.AfterExitPool(ctx, exiter, pool, numShares, exitCoins)
		if err != nil {
			return err
		}
	}
	return nil
}
