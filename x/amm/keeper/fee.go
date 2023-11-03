package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
)

var WeightRecoveryPortion = sdk.NewDecWithPrec(10, 2) // 10%

func PortionCoins(coins sdk.Coins, portion sdk.Dec) sdk.Coins {
	portionCoins := sdk.Coins{}
	for _, coin := range coins {
		portionAmount := sdk.NewDecFromInt(coin.Amount).Mul(portion).RoundInt()
		portionCoins = portionCoins.Add(sdk.NewCoin(
			coin.Denom, portionAmount,
		))
	}
	return portionCoins
}

func (k Keeper) OnCollectFee(ctx sdk.Context, pool types.Pool, fee sdk.Coins) error {
	poolRevenueAddress := types.NewPoolRevenueAddress(pool.PoolId)
	weightRecoveryFee := PortionCoins(fee, WeightRecoveryPortion)
	revenueAmount := fee.Sub(weightRecoveryFee...)
	err := k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(pool.RebalanceTreasury), poolRevenueAddress, revenueAmount)
	if err != nil {
		return nil
	}
	return k.SwapFeesToRevenueToken(ctx, pool, revenueAmount)
}

// No fee management required when doing swap from fees to revenue token
func (k Keeper) SwapFeesToRevenueToken(ctx sdk.Context, pool types.Pool, fee sdk.Coins) error {
	poolRevenueAddress := types.NewPoolRevenueAddress(pool.PoolId)
	for _, tokenIn := range fee {
		// skip for fee denom
		if tokenIn.Denom == pool.PoolParams.FeeDenom {
			continue
		}
		// Executes the swap in the pool and stores the output. Updates pool assets but
		// does not actually transfer any tokens to or from the pool.
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		tokenOutCoin, _, _, err := pool.SwapOutAmtGivenIn(ctx, k.oracleKeeper, &snapshot, sdk.Coins{tokenIn}, pool.PoolParams.FeeDenom, sdk.ZeroDec(), k.accountedPoolKeeper)
		if err != nil {
			return err
		}

		tokenOutAmount := tokenOutCoin.Amount

		if !tokenOutAmount.IsPositive() {
			return sdkerrors.Wrapf(types.ErrInvalidMathApprox, "token amount must be positive")
		}

		// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
		// Also emits a swap event and updates related liquidity metrics.
		err, _ = k.UpdatePoolForSwap(ctx, pool, poolRevenueAddress, tokenIn, tokenOutCoin, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
		if err != nil {
			return err
		}
	}
	return nil
}
