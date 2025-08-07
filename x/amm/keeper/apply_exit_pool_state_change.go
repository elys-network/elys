package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) ApplyExitPoolStateChange(
	ctx sdk.Context, pool types.Pool,
	exiter sdk.AccAddress, numShares sdkmath.Int,
	exitCoins sdk.Coins, isLiquidation bool,
	weightBalanceBonus osmomath.BigDec, takerFees osmomath.BigDec,
	swapFee osmomath.BigDec, slippageCoins sdk.Coins,
	swapInfos []types.SwapInfo,
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
		swapFeeAmount := (osmomath.BigDecFromSDKInt(exitCoins[0].Amount).Mul(swapFee)).Quo(osmomath.OneBigDec().Sub(swapFee))
		swapFeeInCoins = sdk.Coins{sdk.NewCoin(exitCoins[0].Denom, swapFeeAmount.Dec().RoundInt())}
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
	swapFeeValueInUSD := k.CalculateCoinsUSDValue(ctx, swapFeeInCoins)
	slippageAmountInUSD := k.CalculateCoinsUSDValue(ctx, slippageCoins)

	takerFeesAmountInUSD := k.CalculateCoinsUSDValue(ctx, takerFeesInCoins)

	// emit swap fees event
	if !(swapFeeValueInUSD.IsZero() &&
		slippageAmountInUSD.IsZero() &&
		takerFeesAmountInUSD.IsZero()) {
		types.EmitSwapFeesCollectedEvent(
			ctx,
			swapFeeValueInUSD.String(),
			slippageAmountInUSD.String(),
			"0",
			"0",
			takerFeesAmountInUSD.String(),
		)
	}

	if exitCoins.Len() == 1 {
		// swapInfos contains the internal swaps without fees
		types.EmitSwapsInfoEvent(ctx, pool.PoolId, exiter.String(), swapInfos)
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
