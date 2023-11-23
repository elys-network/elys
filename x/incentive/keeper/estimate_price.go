package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Estimate the price : eg, 1 Eden -> x usdc
func (k Keeper) EstimatePrice(ctx sdk.Context, tokenIn sdk.Coin, baseCurrency string) sdk.Int {
	// Find a pool that can convert tokenIn to usdc
	pool, found := k.FindPool(ctx, tokenIn.Denom, baseCurrency)
	if !found {
		return sdk.ZeroInt()
	}

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, pool)
	swapResult, err := k.amm.CalcOutAmtGivenIn(ctx, pool.PoolId, k.oracleKeeper, &snapshot, sdk.Coins{tokenIn}, baseCurrency, sdk.ZeroDec())

	if err != nil {
		return sdk.ZeroInt()
	}

	if swapResult.IsZero() {
		return sdk.ZeroInt()
	}
	return swapResult.Amount
}
