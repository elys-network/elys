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
	isLiquidation bool,
) (exitCoins, exitCoinsAfterExitFee sdk.Coins, err error) {
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return sdk.Coins{}, sdk.Coins{}, types.ErrInvalidPoolId
	}

	totalSharesAmount := pool.GetTotalShares()
	if shareInAmount.GTE(totalSharesAmount.Amount) {
		return sdk.Coins{}, sdk.Coins{}, errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit >= the number of shares contained in the pool.")
	} else if shareInAmount.LTE(sdk.ZeroInt()) {
		return sdk.Coins{}, sdk.Coins{}, errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit a negative amount of shares")
	}
	exitCoins, err = pool.ExitPool(ctx, k.oracleKeeper, k.accountedPoolKeeper, shareInAmount, tokenOutDenom)
	if err != nil {
		return sdk.Coins{}, sdk.Coins{}, err
	}
	if !tokenOutMins.DenomsSubsetOf(exitCoins) || tokenOutMins.IsAnyGT(exitCoins) {
		return sdk.Coins{}, sdk.Coins{}, errorsmod.Wrapf(types.ErrLimitMinAmount,
			"Exit pool returned %s , minimum tokens out specified as %s",
			exitCoins, tokenOutMins)
	}

	exitCoinsAfterExitFee, err = k.ApplyExitPoolStateChange(ctx, pool, sender, shareInAmount, exitCoins, isLiquidation)
	if err != nil {
		return sdk.Coins{}, sdk.Coins{}, err
	}

	// Decrease liquidty amount
	k.RecordTotalLiquidityDecrease(ctx, exitCoins)

	return exitCoins, exitCoinsAfterExitFee, nil
}
