package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func PortionCoins(coins sdk.Coins, portion osmomath.BigDec) sdk.Coins {
	portionCoins := sdk.Coins{}
	for _, coin := range coins {
		portionAmount := osmomath.BigDecFromSDKInt(coin.Amount).Mul(portion).Dec().RoundInt()
		portionCoins = portionCoins.Add(sdk.NewCoin(
			coin.Denom, portionAmount,
		))
	}
	return portionCoins
}

func (k Keeper) OnCollectFee(ctx sdk.Context, pool types.Pool, fee sdk.Coins) error {
	poolRevenueAddress := types.NewPoolRevenueAddress(pool.PoolId)
	revenueAmount := fee
	if pool.PoolParams.UseOracle {
		params := k.GetParams(ctx)
		weightRecoveryFee := PortionCoins(fee, params.GetBigDecWeightRecoveryFeePortion())
		revenueAmount = fee.Sub(weightRecoveryFee...)
	}

	err := k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(pool.RebalanceTreasury), poolRevenueAddress, revenueAmount)
	if err != nil {
		return err
	}

	// handling the case, pool does not enough liquidity to swap fees to revenue token when liquidity is being fully removed
	cacheCtx, write := ctx.CacheContext()
	err = k.SwapFeesToRevenueToken(cacheCtx, pool, revenueAmount)
	if err == nil {
		write()
	}
	return nil
}

// No fee management required when doing swap from fees to revenue token
func (k Keeper) SwapFeesToRevenueToken(ctx sdk.Context, pool types.Pool, fee sdk.Coins) error {
	poolRevenueAddress := types.NewPoolRevenueAddress(pool.PoolId)
	params := k.GetParams(ctx)
	takersFees := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()
	for _, tokenIn := range fee {
		// skip for fee denom
		if tokenIn.Denom == pool.PoolParams.FeeDenom {
			continue
		}
		// Executes the swap in the pool and stores the output. Updates pool assets but
		// does not actually transfer any tokens to or from the pool.
		snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)
		tokenOutCoin, _, slippageAmount, _, oracleOutAmount, _, err := pool.SwapOutAmtGivenIn(ctx, k.oracleKeeper, snapshot, sdk.Coins{tokenIn}, pool.PoolParams.FeeDenom, osmomath.ZeroBigDec(), k.accountedPoolKeeper, osmomath.OneBigDec(), params, takersFees)
		if err != nil {
			return err
		}

		tokenOutAmount := tokenOutCoin.Amount

		if !tokenOutAmount.IsPositive() {
			return errorsmod.Wrapf(types.ErrInvalidMathApprox, "token amount must be positive")
		}

		// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
		// Also emits a swap event and updates related liquidity metrics.
		_, err = k.UpdatePoolForSwap(ctx, pool, poolRevenueAddress, poolRevenueAddress, tokenIn, tokenOutCoin, osmomath.ZeroBigDec(), slippageAmount, sdkmath.ZeroInt(), oracleOutAmount.Dec().TruncateInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), false)
		if err != nil {
			return err
		}
	}
	return nil
}
