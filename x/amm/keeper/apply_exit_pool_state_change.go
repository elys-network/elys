package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) ApplyExitPoolStateChange(
	ctx sdk.Context, pool types.Pool,
	exiter sdk.AccAddress, numShares sdkmath.Int,
	exitCoins sdk.Coins, isLiquidation bool,
	weightBalanceBonus sdkmath.LegacyDec, takerFees sdkmath.LegacyDec,
	swapFee sdkmath.LegacyDec, slippageCoins sdk.Coins,
) error {
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

	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())

	swapFeeInCoins := sdk.Coins{}
	// As swapfee will always be positive only if there is single asset in exitCoins
	if swapFee.IsPositive() {
		swapFeeAmount := (exitCoins[0].Amount.ToLegacyDec().Mul(swapFee)).Quo(sdkmath.LegacyOneDec().Sub(swapFee))
		swapFeeInCoins = sdk.Coins{sdk.NewCoin(exitCoins[0].Denom, swapFeeAmount.RoundInt())}
	}

	// send swap fee to rebalance treasury
	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err := k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
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

	// Taker fees
	takerFeesInCoins := sdk.Coins{}
	if takerFees.IsPositive() {
		takerFeesInCoins = PortionCoins(exitCoins, takerFees)
	}

	// send taker fee to protocol treasury
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

	// convert the fees into USD
	swapFeeValueInUSD := k.CalculateCoinsUSDValue(ctx, swapFeeInCoins).String()
	slippageAmountInUSD := k.CalculateCoinsUSDValue(ctx, slippageCoins).String()

	takerFeesAmountInUSD := k.CalculateCoinsUSDValue(ctx, takerFeesInCoins).String()

	// emit swap fees event
	types.EmitSwapFeesCollectedEvent(ctx, swapFeeValueInUSD, slippageAmountInUSD, "0", "0", takerFeesAmountInUSD)

	types.EmitRemoveLiquidityEvent(ctx, exiter, pool.GetPoolId(), exitCoins)
	if k.hooks != nil {
		err = k.hooks.AfterExitPool(ctx, exiter, pool, numShares, exitCoins)
		if err != nil {
			return err
		}
	}
	return nil
}
