package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Estimate the price : eg, 1 Eden -> x usdc
func (k Keeper) EstimatePrice(ctx sdk.Context, tokenIn sdk.Coin, baseCurrency string) math.Int {
	// Find a pool that can convert tokenIn to usdc
	pool, found := k.amm.GetBestPoolWithDenoms(ctx, []string{tokenIn.Denom, baseCurrency})
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

func (k Keeper) GetEdenPrice(ctx sdk.Context, baseCurrency string) math.LegacyDec {
	// Calc Eden price in usdc
	// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
	// TODO: replace to use spot price
	// TODO: Remember to use the $ value of Eden price and not eden/usdc
	edenPrice := k.EstimatePrice(ctx, sdk.NewCoin(ptypes.Elys, sdk.NewInt(100000)), baseCurrency)
	edenPriceDec := sdk.NewDecFromInt(edenPrice).Quo(sdk.NewDec(100000))
	if edenPriceDec.IsZero() {
		edenPriceDec = sdk.OneDec()
	}
	return edenPriceDec
}
