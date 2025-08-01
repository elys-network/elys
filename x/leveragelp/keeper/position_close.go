package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v7/x/amm/keeper"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) CheckHealthStopLossThenRepayAndClose(ctx sdk.Context, position *types.Position, pool *types.Pool, closingRatio osmomath.BigDec, isLiquidation bool) (osmomath.BigDec, math.Int, sdk.Coins, math.Int, sdk.Coins, osmomath.BigDec, bool, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, error) {

	positionHealth, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), false, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	safetyFactor := osmomath.BigDecFromDec(k.GetSafetyFactor(ctx))

	if positionHealth.LTE(safetyFactor) {
		closingRatio = osmomath.OneBigDec()
	}

	ammPool, found := k.amm.GetPool(ctx, position.AmmPoolId)
	if !found {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), false, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrAmmPoolNotFound
	}

	// Note this price of One Share (multiplied by 1^18)
	lpTokenPrice, err := ammPool.LpTokenPriceForShare(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), false, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	stopLossReached := position.CheckStopLossReached(lpTokenPrice)
	if stopLossReached {
		closingRatio = osmomath.OneBigDec()
	}

	if closingRatio.IsNil() || closingRatio.IsNegative() || closingRatio.IsZero() || closingRatio.GT(osmomath.OneBigDec()) {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrInvalidClosingRatio
	}

	// This will be the net amount that will be closed
	totalLpAmountToClose := closingRatio.Mul(osmomath.BigDecFromSDKInt(position.LeveragedLpAmount)).Dec().TruncateInt()

	ammPoolTVL, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	collateralDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, position.Collateral.Denom)

	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
	// Amount required to pay back
	repayAmount := debt.GetBigDecTotalLiablities().Mul(closingRatio).Dec().TruncateInt()
	repayValue := collateralDenomPrice.Mul(osmomath.BigDecFromSDKInt(repayAmount))

	ammPoolTotalShareAmountDec := osmomath.BigDecFromSDKInt(ammPool.TotalShares.Amount)
	lpSharesForRepay := repayValue.Mul(ammPoolTotalShareAmountDec).Quo(ammPoolTVL).Ceil().Dec().TruncateInt() // round up

	// Bot failed to close position here, Position is unhealthy, override closing ratio to 1 and position health is < 1
	// Also important because round up above. If position health is 1, then lpSharesForRepay == totalLpAmountToClose, rounding up before might have made it 1 higher
	if lpSharesForRepay.GT(totalLpAmountToClose) {
		closingRatio = osmomath.OneBigDec()
		totalLpAmountToClose = position.LeveragedLpAmount
		repayAmount = debt.GetTotalLiablities()
		repayValue = collateralDenomPrice.Mul(osmomath.BigDecFromSDKInt(repayAmount))
		lpSharesForRepay = position.LeveragedLpAmount
	}

	// we calculate weight breaking fee (-ve of weightBalanceBonus if there is one) that could have occurred
	_, weightBalanceBonus, slippage, swapFee, takerFee, err := k.amm.ExitPoolEst(ctx, position.AmmPoolId, lpSharesForRepay, position.Collateral.Denom)
	if err != nil {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	exitFeeOnClosingPosition := osmomath.ZeroBigDec()
	percentageExitLeverageFee := osmomath.ZeroBigDec()
	weightBreakingFee := osmomath.ZeroBigDec()
	weightBreakingFeeValue := osmomath.ZeroBigDec()
	// weight breaking fee exists when weightBalanceBonus is negative
	// There is anything left to pay if totalLpAmountToClose > lpSharesForRepay, also it prevents denominator going 0
	if weightBalanceBonus.IsNegative() {
		weightBreakingFee = weightBalanceBonus.Abs()
		weightBreakingFeeValue = repayValue.Mul(weightBreakingFee)

		// it equals weightBreakingFeeValue / totalLpAmountToCloseValue
		exitFeeOnClosingPosition = weightBreakingFeeValue.Mul(ammPoolTotalShareAmountDec).Quo(osmomath.BigDecFromSDKInt(totalLpAmountToClose).Mul(ammPoolTVL))
	}

	slippageValue := slippage.Mul(repayValue)
	swapFeeValue := swapFee.Mul(repayValue)
	takerFeeValue := takerFee.Mul(repayValue)
	exitSlippageFeeOnClosingPosition := (slippageValue).Mul(ammPoolTotalShareAmountDec).Quo(osmomath.BigDecFromSDKInt(totalLpAmountToClose).Mul(ammPoolTVL))
	exitSwapFeeOnClosingPosition := (swapFeeValue).Mul(ammPoolTotalShareAmountDec).Quo(osmomath.BigDecFromSDKInt(totalLpAmountToClose).Mul(ammPoolTVL))
	exitTakerFeeOnClosingPosition := (takerFeeValue).Mul(ammPoolTotalShareAmountDec).Quo(osmomath.BigDecFromSDKInt(totalLpAmountToClose).Mul(ammPoolTVL))

	if totalLpAmountToClose.GT(lpSharesForRepay) {
		// weightBreakingFeeValue / ((totalLpAmountToClose - lpSharesForRepay) x LP Price)
		denominator := osmomath.BigDecFromSDKInt(totalLpAmountToClose.Sub(lpSharesForRepay)).Mul(ammPoolTVL).Quo(ammPoolTotalShareAmountDec)
		percentageExitLeverageFee = (weightBreakingFeeValue.Add(slippageValue).Add(swapFeeValue).Add(takerFeeValue)).Quo(denominator)
	}

	if percentageExitLeverageFee.GT(osmomath.OneBigDec()) {
		percentageExitLeverageFee = osmomath.OneBigDec()
	}

	// We subtract here because CheckAmmBalance in hooks will validate if there is enough amount to pay liabilities, since we're exiting this one, it needs reduced value
	k.stableKeeper.SubtractPoolLiabilities(ctx, position.AmmPoolId, sdk.NewCoin(position.Collateral.Denom, repayAmount))

	// Subtract amount that is being reduced from lp pool as it gets checked in CheckAmmBalance
	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(lpSharesForRepay)
	pool.UpdateAssetLeveragedAmount(ctx, position.Collateral.Denom, lpSharesForRepay, false)
	// pool is set here
	k.UpdatePoolHealth(ctx, pool)

	_, _, _, _, _, err = k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, lpSharesForRepay, sdk.Coins{}, position.Collateral.Denom, isLiquidation, false)
	if err != nil {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	// Updating ammPool
	ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)

	// Important to add back exact same amount we reduced as Repay function subtracts the whatever amount that is actually paid back
	k.stableKeeper.AddPoolLiabilities(ctx, position.AmmPoolId, sdk.NewCoin(position.Collateral.Denom, repayAmount))

	// Check if position has enough coins to repay else repay partial
	positionCollateralBalance := k.bankKeeper.GetBalance(ctx, position.GetPositionAddress(), position.Collateral.Denom)
	// Position is unhealthy, pay back whatever position can
	if positionCollateralBalance.Amount.LT(repayAmount) {
		repayAmount = positionCollateralBalance.Amount
	}

	if repayAmount.IsPositive() {
		err = k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, repayAmount), position.BorrowPoolId, position.AmmPoolId)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
	}

	// Position is healthy, position gets rewards in two tokens
	var coinsLeftAfterRepay sdk.Coins
	sharesLeft := math.ZeroInt()
	if totalLpAmountToClose.GT(lpSharesForRepay) {
		sharesLeft = totalLpAmountToClose.Sub(lpSharesForRepay)
		coinsLeftAfterRepay, _, _, _, _, err = k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, sharesLeft, sdk.Coins{}, "", isLiquidation, false)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		// Updating ammPool
		ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)
	}

	// Set collateral to same % as reduction in LP position
	collateralReduced := osmomath.BigDecFromSDKInt(position.Collateral.Amount).Mul(closingRatio).Dec().TruncateInt()
	position.Collateral.Amount = position.Collateral.Amount.Sub(collateralReduced)
	// Update leveragedLpAmount
	position.LeveragedLpAmount = position.LeveragedLpAmount.Sub(totalLpAmountToClose)
	posHealth, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	position.PositionHealth = posHealth.Dec()

	// Update Liabilities
	debt = k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.BorrowPoolId)
	position.Liabilities = debt.GetTotalLiablities()

	// Update the pool health.
	if sharesLeft.IsPositive() {
		pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(sharesLeft)
		pool.UpdateAssetLeveragedAmount(ctx, position.Collateral.Denom, sharesLeft, false)
	}

	var coinsForAmm sdk.Coins
	if percentageExitLeverageFee.IsPositive() && len(coinsLeftAfterRepay) > 0 {
		coinsForAmm = ammkeeper.PortionCoins(coinsLeftAfterRepay, percentageExitLeverageFee)
		weightBreakingFeePortion := k.amm.GetParams(ctx).WeightBreakingFeePortion

		coinsToAmmRebalancer := ammkeeper.PortionCoins(coinsForAmm, osmomath.BigDecFromDec(weightBreakingFeePortion))
		coinsToAmmPool := coinsForAmm.Sub(coinsToAmmRebalancer...)

		// Very important to fetch this again, Updating ammPool
		ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)
		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(ammPool.Address), coinsToAmmPool)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		err = k.amm.AddToPoolBalanceAndUpdateLiquidity(ctx, &ammPool, math.ZeroInt(), coinsToAmmPool)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(ammPool.RebalanceTreasury), coinsToAmmRebalancer)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

	}

	// anything left over in position balance goes to user
	// We do not do finalUserTokens = coinsLeftAfterRepay.Sub(coinsForAmm) because there might be some rewards by leverageLP position owner
	finalUserTokens := k.bankKeeper.GetAllBalances(ctx, position.GetPositionAddress())
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	if len(finalUserTokens) > 0 {
		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), positionOwner, finalUserTokens)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
	}

	k.SetPosition(ctx, position)
	// Final pool update, it's also set here
	k.UpdatePoolHealth(ctx, pool)

	if position.LeveragedLpAmount.IsZero() {
		// As we have already exited the pool, we need to delete the position
		err = k.masterchefKeeper.ClaimRewards(ctx, position.GetPositionAddress(), []uint64{position.AmmPoolId}, positionOwner)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		err = k.DestroyPosition(ctx, positionOwner, position.Id)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		// Delete Debt even if it's not paid fully, can happen only if bot fails to close and position health goes below 1
		if position.Liabilities.IsPositive() {
			err = k.stableKeeper.CloseOnUnableToRepay(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
			if err != nil {
				return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
			}
		}
	}

	// Important to run hooks to updated accounted pool balances
	if k.hooks != nil {
		// Updating ammPool
		ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)
		err = k.hooks.AfterLeverageLpPositionClose(ctx, position.GetOwnerAddress(), ammPool)
		if err != nil {
			return osmomath.ZeroBigDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, osmomath.ZeroBigDec(), stopLossReached, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
	}

	return closingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, finalUserTokens, exitFeeOnClosingPosition, stopLossReached, weightBreakingFee, exitSlippageFeeOnClosingPosition, exitSwapFeeOnClosingPosition, exitTakerFeeOnClosingPosition, slippageValue, swapFeeValue, takerFeeValue, weightBreakingFeeValue, nil
}
