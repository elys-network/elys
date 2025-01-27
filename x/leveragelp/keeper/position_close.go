package keeper

import (
	"cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CheckHealthStopLossThenRepayAndClose(ctx sdk.Context, position *types.Position, pool *types.Pool, closingRatio math.LegacyDec, isLiquidation bool) (math.LegacyDec, math.Int, sdk.Coins, math.Int, sdk.Coins, error) {
	// Ensure position.LeveragedLpAmount is not zero to avoid division by zero
	if position.LeveragedLpAmount.IsZero() {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, types.ErrAmountTooLow
	}
	//
	//if position.PositionHealth {
	//
	//}
	//
	//if position.StopLossPrice {
	//
	//}

	if closingRatio.IsNil() || closingRatio.IsNegative() || closingRatio.IsZero() || closingRatio.GT(math.LegacyOneDec()) {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, errors.New("invalid closing ratio for leverage lp")
	}

	// This will be the net amount that will be closed
	totalLpAmountToClose := closingRatio.MulInt(position.LeveragedLpAmount).TruncateInt()

	ammPool, found := k.amm.GetPool(ctx, position.AmmPoolId)
	if !found {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, types.ErrAmmPoolNotFound
	}
	ammPoolTVL, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
	}
	collateralDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)

	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress())
	// Amount required to pay back
	repayAmount := debt.GetTotalLiablities().ToLegacyDec().Mul(closingRatio).TruncateInt() // Important to round up
	repayValue := collateralDenomPrice.MulInt(repayAmount)

	lpSharesForRepay := repayValue.Mul(ammPool.TotalShares.Amount.ToLegacyDec()).Quo(ammPoolTVL).TruncateInt()

	// Bot failed to close position here, Position is unhealthy, override closing ratio to 1 and position health is < 1
	if lpSharesForRepay.GT(totalLpAmountToClose) {
		closingRatio = math.LegacyOneDec()
		totalLpAmountToClose = position.LeveragedLpAmount
		repayAmount = debt.GetTotalLiablities()
		repayValue = collateralDenomPrice.MulInt(repayAmount)
		lpSharesForRepay = position.LeveragedLpAmount
	}

	// we calculate weight breaking fee (-ve of weightBalanceBonus if there is one) that could have occurred
	_, weightBalanceBonus, err := k.amm.ExitPoolEst(ctx, position.AmmPoolId, lpSharesForRepay, position.Collateral.Denom)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
	}

	percentageExitLeverageFee := math.LegacyZeroDec()
	// weight breaking fee exists when weightBalanceBonus is negative
	// There is anything left to pay if totalLpAmountToClose > lpSharesForRepay, also it prevents denominator going 0
	if weightBalanceBonus.IsNegative() && totalLpAmountToClose.GT(lpSharesForRepay) {
		weightBreakingFee := weightBalanceBonus.Abs()
		weightBreakingFeeValue := repayValue.Mul(weightBreakingFee)

		// weightBreakingFeeValue / ((totalLpAmountToClose - lpSharesForRepay) x LP Price)
		percentageExitLeverageFee = weightBreakingFeeValue.Quo(totalLpAmountToClose.Sub(lpSharesForRepay).ToLegacyDec().Mul(ammPool.TotalShares.Amount.ToLegacyDec()).Quo(ammPoolTVL))

		if percentageExitLeverageFee.GT(math.LegacyOneDec()) {
			percentageExitLeverageFee = math.LegacyOneDec()
		}
	}
	// TODO reduce liabilities from global

	_, _, err = k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, lpSharesForRepay, sdk.Coins{}, position.Collateral.Denom, isLiquidation, false)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
	}

	// Check if position has enough coins to repay else repay partial
	positionCollateralBalance := k.bankKeeper.GetBalance(ctx, position.GetPositionAddress(), position.Collateral.Denom)
	// Position is unhealthy, pay back whatever position can
	if positionCollateralBalance.Amount.LT(repayAmount) {
		repayAmount = positionCollateralBalance.Amount
	}

	if repayAmount.IsPositive() {
		err = k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, repayAmount))
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}
	}

	// Position is healthy, position gets rewards in two tokens
	var coinsLeftAfterRepay sdk.Coins
	if totalLpAmountToClose.GT(lpSharesForRepay) {
		rewardShares := totalLpAmountToClose.Sub(lpSharesForRepay)
		coinsLeftAfterRepay, _, err = k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, rewardShares, sdk.Coins{}, "", isLiquidation, false)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}
	}

	// Set collateral to same % as reduction in LP position
	collateralReduced := position.Collateral.Amount.ToLegacyDec().Mul(closingRatio).TruncateInt()
	position.Collateral.Amount = position.Collateral.Amount.Sub(collateralReduced)
	// Update leveragedLpAmount
	position.LeveragedLpAmount = position.LeveragedLpAmount.Sub(totalLpAmountToClose)
	position.PositionHealth, err = k.GetPositionHealth(ctx, *position)
	if err != nil {
		return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
	}

	// Update Liabilities
	debt = k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress())
	position.Liabilities = debt.GetTotalLiablities()

	// Update the pool health.
	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(totalLpAmountToClose)

	var coinsForAmm sdk.Coins
	if percentageExitLeverageFee.IsPositive() && len(coinsLeftAfterRepay) > 0 {
		coinsForAmm = ammkeeper.PortionCoins(coinsLeftAfterRepay, percentageExitLeverageFee)
		weightBreakingFeePortion := k.amm.GetParams(ctx).WeightBreakingFeePortion

		coinsToAmmRebalancer := ammkeeper.PortionCoins(coinsLeftAfterRepay, weightBreakingFeePortion)
		coinsToAmmPool := coinsForAmm.Sub(coinsToAmmRebalancer...)

		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(ammPool.Address), coinsToAmmPool)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}
		err = k.amm.AddToPoolBalanceAndUpdateLiquidity(ctx, &ammPool, math.ZeroInt(), coinsToAmmPool)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}

		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(ammPool.RebalanceTreasury), coinsToAmmRebalancer)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}

	}

	// anything left over in position balance goes to user
	// We do not do finalUserRewards = coinsLeftAfterRepay.Sub(coinsForAmm) because there might be some rewards by leverageLP position owner
	finalUserRewards := k.bankKeeper.GetAllBalances(ctx, position.GetPositionAddress())
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	if len(finalUserRewards) > 0 {
		err = k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), positionOwner, finalUserRewards)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}
	}

	k.UpdatePoolHealth(ctx, pool)
	k.SetPosition(ctx, position)

	if position.LeveragedLpAmount.IsZero() {
		// As we have already exited the pool, we need to delete the position
		err = k.masterchefKeeper.ClaimRewards(ctx, position.GetPositionAddress(), []uint64{position.AmmPoolId}, positionOwner)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}
		err = k.DestroyPosition(ctx, positionOwner, position.Id)
		if err != nil {
			return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
		}

		// Delete Debt even if it's not paid fully, can happen only if bot fails to close and position health goes below 1
		if position.Liabilities.IsPositive() {
			err = k.stableKeeper.CloseOnUnableToRepay(ctx, position.GetPositionAddress())
			if err != nil {
				return math.LegacyZeroDec(), math.ZeroInt(), sdk.Coins{}, math.ZeroInt(), sdk.Coins{}, err
			}
		}
	}

	return closingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, finalUserRewards, nil
}
