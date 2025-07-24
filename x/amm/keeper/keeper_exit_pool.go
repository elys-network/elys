package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) ExitPool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	shareInAmount math.Int,
	tokenOutMins sdk.Coins,
	tokenOutDenom string,
	isLiquidation, applyWeightBreakingFee bool,
) (exitCoins sdk.Coins, weightBalanceBonus osmomath.BigDec, swapFee osmomath.BigDec, slippage osmomath.BigDec, takerFeesFinal osmomath.BigDec, err error) {
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrInvalidPoolId
	}

	totalSharesAmount := pool.GetTotalShares()
	if shareInAmount.GTE(totalSharesAmount.Amount) {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit >= the number of shares contained in the pool.")
	} else if shareInAmount.LTE(math.ZeroInt()) {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit a negative amount of shares")
	}
	params := k.GetParams(ctx)
	takersFees := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()
	snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)
	exitCoins, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, swapInfos, err := pool.ExitPool(ctx, k.oracleKeeper, k.accountedPoolKeeper, snapshot, shareInAmount, tokenOutDenom, params, takersFees, applyWeightBreakingFee)
	if err != nil {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	// Check treasury and update weightBalance
	if weightBalanceBonus.IsPositive() && exitCoins.Len() == 1 {
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, exitCoins[0].Denom).Amount

		bonusTokenAmount := osmomath.BigDecFromSDKInt(exitCoins[0].Amount).Mul(weightBalanceBonus).Dec().TruncateInt()

		if treasuryTokenAmount.LT(bonusTokenAmount) {
			weightBalanceBonus = osmomath.BigDecFromSDKInt(treasuryTokenAmount).Quo(osmomath.BigDecFromSDKInt(exitCoins[0].Amount))
		}
	}

	if pool.PoolParams.UseOracle {
		if slippageCoins.IsAllPositive() {
			k.TrackWeightBreakingSlippage(ctx, pool.PoolId, sdk.NewCoin(slippageCoins[0].Denom, slippageCoins[0].Amount))
		}
	}

	if !tokenOutMins.DenomsSubsetOf(exitCoins) || tokenOutMins.IsAnyGT(exitCoins) {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), errorsmod.Wrapf(types.ErrLimitMinAmount,
			"Exit pool returned %s , minimum tokens out specified as %s",
			exitCoins, tokenOutMins)
	}

	err = k.ApplyExitPoolStateChange(ctx, pool, sender, shareInAmount, exitCoins, isLiquidation, weightBalanceBonus, takerFeesFinal, swapFee, slippageCoins, swapInfos)
	if err != nil {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	err = k.RecordTotalLiquidityDecrease(ctx, exitCoins)
	if err != nil {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	return exitCoins, weightBalanceBonus, swapFee, slippage, takerFeesFinal, nil
}
