package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) ExitPool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	shareInAmount math.Int,
	tokenOutMins sdk.Coins,
	tokenOutDenom string,
	isLiquidation, applyWeightBreakingFee bool,
) (exitCoins sdk.Coins, weightBalanceBonus math.LegacyDec, swapFee math.LegacyDec, slippage math.LegacyDec, err error) {
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrInvalidPoolId
	}

	totalSharesAmount := pool.GetTotalShares()
	if shareInAmount.GTE(totalSharesAmount.Amount) {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit >= the number of shares contained in the pool.")
	} else if shareInAmount.LTE(math.ZeroInt()) {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit a negative amount of shares")
	}
	params := k.GetParams(ctx)
	exitCoins, weightBalanceBonus, slippage, swapFee, err = pool.ExitPool(ctx, k.oracleKeeper, k.accountedPoolKeeper, shareInAmount, tokenOutDenom, params, applyWeightBreakingFee)
	if err != nil {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}
	if !tokenOutMins.DenomsSubsetOf(exitCoins) || tokenOutMins.IsAnyGT(exitCoins) {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrLimitMinAmount,
			"Exit pool returned %s , minimum tokens out specified as %s",
			exitCoins, tokenOutMins)
	}

	err = k.ApplyExitPoolStateChange(ctx, pool, sender, shareInAmount, exitCoins, isLiquidation, weightBalanceBonus, swapFee)
	if err != nil {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	err = k.RecordTotalLiquidityDecrease(ctx, exitCoins)
	if err != nil {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	return exitCoins, weightBalanceBonus, slippage, swapFee, nil
}
