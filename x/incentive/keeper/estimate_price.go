package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Estimate the price : eg, 1 Eden -> x usdc
func (k Keeper) EstimatePrice(ctx sdk.Context, tokenInDenom, baseCurrency string) math.LegacyDec {
	// Find a pool that can convert tokenIn to usdc
	pool, found := k.amm.GetBestPoolWithDenoms(ctx, []string{tokenInDenom, baseCurrency})
	if !found {
		return sdk.ZeroDec()
	}

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, pool)

	rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, baseCurrency, k.accountedPoolKeeper)
	if err != nil {
		return sdk.ZeroDec()
	}

	return rate
}

func (k Keeper) GetEdenPrice(ctx sdk.Context, baseCurrency string) math.LegacyDec {
	// Calc Eden price in usdc
	// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
	edenPrice := k.EstimatePrice(ctx, ptypes.Elys, baseCurrency)
	if edenPrice.IsZero() {
		edenPrice = sdk.OneDec()
	}
	return edenPrice
}
