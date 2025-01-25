package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) ExitPool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	shareInAmount math.Int,
	tokenOutMins sdk.Coins,
	tokenOutDenom string,
	isLiquidation bool,
) (exitCoins sdk.Coins, weightBalanceBonus elystypes.Dec34, err error) {
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return sdk.Coins{}, elystypes.ZeroDec34(), types.ErrInvalidPoolId
	}

	totalSharesAmount := pool.GetTotalShares()
	if shareInAmount.GTE(totalSharesAmount.Amount) {
		return sdk.Coins{}, elystypes.ZeroDec34(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit >= the number of shares contained in the pool.")
	} else if shareInAmount.LTE(math.ZeroInt()) {
		return sdk.Coins{}, elystypes.ZeroDec34(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit a negative amount of shares")
	}
	params := k.GetParams(ctx)
	exitCoins, weightBalanceBonus, err = pool.ExitPool(ctx, k.oracleKeeper, k.accountedPoolKeeper, shareInAmount, tokenOutDenom, params)
	if err != nil {
		return sdk.Coins{}, elystypes.ZeroDec34(), err
	}
	if !tokenOutMins.DenomsSubsetOf(exitCoins) || tokenOutMins.IsAnyGT(exitCoins) {
		return sdk.Coins{}, elystypes.ZeroDec34(), errorsmod.Wrapf(types.ErrLimitMinAmount,
			"Exit pool returned %s , minimum tokens out specified as %s",
			exitCoins, tokenOutMins)
	}

	err = k.ApplyExitPoolStateChange(ctx, pool, sender, shareInAmount, exitCoins, isLiquidation, weightBalanceBonus)
	if err != nil {
		return sdk.Coins{}, elystypes.ZeroDec34(), err
	}

	err = k.RecordTotalLiquidityDecrease(ctx, exitCoins)
	if err != nil {
		return sdk.Coins{}, elystypes.ZeroDec34(), err
	}

	return exitCoins, weightBalanceBonus, nil
}
