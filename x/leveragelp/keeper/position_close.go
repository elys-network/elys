package keeper

import (
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CheckHealthStopLossThenRepayAndClose(ctx sdk.Context, position *types.Position, pool *types.Pool, closingRatio math.LegacyDec, isLiquidation bool) (math.LegacyDec, math.Int, sdk.Coins, math.Int, sdk.Coins, math.LegacyDec, bool, math.LegacyDec, math.LegacyDec, math.LegacyDec, error) {

	positionHealth, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), false, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}
	safetyFactor := k.GetSafetyFactor(ctx)

	if positionHealth.LTE(safetyFactor) {
		closingRatio = math.LegacyOneDec()
	}

	ammPool, found := k.amm.GetPool(ctx, position.AmmPoolId)
	if !found {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), false, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrAmmPoolNotFound
	}

	// Note this price of One Share (multiplied by 1^18)
	lpTokenPrice, err := ammPool.LpTokenPriceForShare(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), false, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	stopLossReached := position.CheckStopLossReached(lpTokenPrice)
	if stopLossReached {
		closingRatio = math.LegacyOneDec()
	}

	if closingRatio.IsNil() || closingRatio.IsNegative() || closingRatio.IsZero() || closingRatio.GT(math.LegacyOneDec()) {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), errors.New("invalid closing ratio for leverage lp")
	}

	// This will be the net amount that will be closed
	totalLpAmountToClose := closingRatio.MulInt(position.LeveragedLpAmount).TruncateInt()

	ammPoolTVL, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}
	collateralDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)

	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
	// Amount required to pay back
	repayAmount := debt.GetTotalLiablities().ToLegacyDec().Mul(closingRatio).TruncateInt()
	repayValue := collateralDenomPrice.MulInt(repayAmount)

	lpSharesForRepay := repayValue.Mul(ammPool.TotalShares.Amount.ToLegacyDec()).Quo(ammPoolTVL).Ceil().TruncateInt() // round up

	// Bot failed to close position here, Position is unhealthy, override closing ratio to 1 and position health is < 1
	// Also important because round up above. If position health is 1, then lpSharesForRepay == totalLpAmountToClose, rounding up before might have made it 1 higher
	if lpSharesForRepay.GT(totalLpAmountToClose) {
		closingRatio = math.LegacyOneDec()
		totalLpAmountToClose = position.LeveragedLpAmount
		repayAmount = debt.GetTotalLiablities()
		repayValue = collateralDenomPrice.MulInt(repayAmount)
		lpSharesForRepay = position.LeveragedLpAmount
	}

	// we calculate weight breaking fee (-ve of weightBalanceBonus if there is one) that could have occurred
	_, weightBalanceBonus, slippage, swapFee, err := k.amm.ExitPoolEst(ctx, position.AmmPoolId, lpSharesForRepay, position.Collateral.Denom)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	exitFeeOnClosingPosition := math.LegacyZeroDec()
	percentageExitLeverageFee := math.LegacyZeroDec()
	weightBreakingFee := math.LegacyZeroDec()
	weightBreakingFeeValue := math.LegacyZeroDec()
	// weight breaking fee exists when weightBalanceBonus is negative
	// There is anything left to pay if totalLpAmountToClose > lpSharesForRepay, also it prevents denominator going 0
	if weightBalanceBonus.IsNegative() {
		weightBreakingFee = weightBalanceBonus.Abs()
		weightBreakingFeeValue = repayValue.Mul(weightBreakingFee)

		// it equals weightBreakingFeeValue / totalLpAmountToCloseValue
		exitFeeOnClosingPosition = weightBreakingFeeValue.Mul(ammPool.TotalShares.Amount.ToLegacyDec()).Quo(totalLpAmountToClose.ToLegacyDec().Mul(ammPoolTVL))
	}

	// panic("Swap fee: " + swapFee.String())
	slippageValue := slippage.Mul(repayValue)
	swapFeeValue := swapFee.Mul(repayValue)
	exitSlippageFeeOnClosingPosition := (slippageValue).Mul(ammPool.TotalShares.Amount.ToLegacyDec()).Quo(totalLpAmountToClose.ToLegacyDec().Mul(ammPoolTVL))
	exitSwapFeeOnClosingPosition := (swapFeeValue).Mul(ammPool.TotalShares.Amount.ToLegacyDec()).Quo(totalLpAmountToClose.ToLegacyDec().Mul(ammPoolTVL))

	if totalLpAmountToClose.GT(lpSharesForRepay) {
		// weightBreakingFeeValue / ((totalLpAmountToClose - lpSharesForRepay) x LP Price)
		denominator := (totalLpAmountToClose.Sub(lpSharesForRepay).ToLegacyDec()).Mul(ammPoolTVL).Quo(ammPool.TotalShares.Amount.ToLegacyDec())
		percentageExitLeverageFee = (weightBreakingFeeValue.Add(slippageValue).Add(swapFeeValue)).Quo(denominator)
	}

	if percentageExitLeverageFee.GT(math.LegacyOneDec()) {
		percentageExitLeverageFee = math.LegacyOneDec()
	}

	// We subtract here because CheckAmmUsdcBalance in hooks will validate if there is enough amount to pay liabilities, since we're exiting this one, it needs reduced value
	k.stableKeeper.SubtractPoolLiabilities(ctx, position.AmmPoolId, sdk.NewCoin(position.Collateral.Denom, repayAmount))

	// Subtract amount that is being reduced from lp pool as it gets checked in CheckAmmUsdcBalance
	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(lpSharesForRepay)
	// pool is set here
	k.UpdatePoolHealth(ctx, pool)

	_, _, _, _, err = k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, lpSharesForRepay, sdk.Coins{}, position.Collateral.Denom, isLiquidation, false)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
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
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}

	// Position is healthy, position gets rewards in two tokens
	var coinsLeftAfterRepay sdk.Coins
	sharesLeft := math.ZeroInt()
	if totalLpAmountToClose.GT(lpSharesForRepay) {
		sharesLeft = totalLpAmountToClose.Sub(lpSharesForRepay)
		coinsLeftAfterRepay, _, _, _, err = k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, sharesLeft, sdk.Coins{}, "", isLiquidation, false)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
		// Updating ammPool
		ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)
	}

	// Set collateral to same % as reduction in LP position
	collateralReduced := position.Collateral.Amount.ToLegacyDec().Mul(closingRatio).TruncateInt()
	position.Collateral.Amount = position.Collateral.Amount.Sub(collateralReduced)
	// Update leveragedLpAmount
	position.LeveragedLpAmount = position.LeveragedLpAmount.Sub(totalLpAmountToClose)
	position.PositionHealth, err = k.GetPositionHealth(ctx, *position)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	// Update Liabilities
	debt = k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.BorrowPoolId)
	position.Liabilities = debt.GetTotalLiablities()

	// Update the pool health.
	if sharesLeft.IsPositive() {
		pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(sharesLeft)
	}

	var coinsForAmm sdk.Coins
	if percentageExitLeverageFee.IsPositive() && len(coinsLeftAfterRepay) > 0 {
		coinsForAmm = ammkeeper.PortionCoins(coinsLeftAfterRepay, percentageExitLeverageFee)
		weightBreakingFeePortion := k.amm.GetParams(ctx).WeightBreakingFeePortion

		coinsToAmmRebalancer := ammkeeper.PortionCoins(coinsForAmm, weightBreakingFeePortion)
		coinsToAmmPool := coinsForAmm.Sub(coinsToAmmRebalancer...)

		// Very important to fetch this again, Updating ammPool
		ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)
		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(ammPool.Address), coinsToAmmPool)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
		err = k.amm.AddToPoolBalanceAndUpdateLiquidity(ctx, &ammPool, math.ZeroInt(), coinsToAmmPool)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}

		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(ammPool.RebalanceTreasury), coinsToAmmRebalancer)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}

	}

	// anything left over in position balance goes to user
	// We do not do finalUserTokens = coinsLeftAfterRepay.Sub(coinsForAmm) because there might be some rewards by leverageLP position owner
	finalUserTokens := k.bankKeeper.GetAllBalances(ctx, position.GetPositionAddress())
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	if len(finalUserTokens) > 0 {
		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), positionOwner, finalUserTokens)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}

	k.SetPosition(ctx, position)
	// Final pool update, it's also set here
	k.UpdatePoolHealth(ctx, pool)

	if position.LeveragedLpAmount.IsZero() {
		// As we have already exited the pool, we need to delete the position
		err = k.masterchefKeeper.ClaimRewards(ctx, position.GetPositionAddress(), []uint64{position.AmmPoolId}, positionOwner)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
		err = k.DestroyPosition(ctx, positionOwner, position.Id)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}

		// Delete Debt even if it's not paid fully, can happen only if bot fails to close and position health goes below 1
		if position.Liabilities.IsPositive() {
			err = k.stableKeeper.CloseOnUnableToRepay(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
			if err != nil {
				return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
			}
		}
	}

	// Important to run hooks to updated accounted pool balances
	if k.hooks != nil {
		// Updating ammPool
		ammPool, _ = k.amm.GetPool(ctx, position.AmmPoolId)
		err = k.hooks.AfterLeverageLpPositionClose(ctx, position.GetOwnerAddress(), ammPool)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, math.LegacyZeroDec(), stopLossReached, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}

	return closingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, finalUserTokens, exitFeeOnClosingPosition, stopLossReached, weightBreakingFee, exitSlippageFeeOnClosingPosition, exitSwapFeeOnClosingPosition, nil
}
