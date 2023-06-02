package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------Placeholder implementation---------
// --------------------------------------------
// Liquidity pool
type LiquidityPool struct {
	index      string
	TVL        sdk.Int
	multiplier int64
	lpToken    string
	rewards    sdk.Dec
}

// LP Deposits
type LpDeposit struct {
	index        string
	LpIndex      string
	amount       sdk.Int
	ownerAddress string
	denom        string
}

// Dummy Liquidity Keeper
type LiquidityKeeper struct {
	liquidityPool map[string]*LiquidityPool
	lpDeposits    []LpDeposit
}

// New Liquidity Keeper
func NewLiquidityKeeper() *LiquidityKeeper {
	return &LiquidityKeeper{
		liquidityPool: make(map[string]*LiquidityPool),
		lpDeposits:    make([]LpDeposit, 0),
	}
}

// IterateLiquidty iterates over all LiquidityPools and performs a
// callback.
func (k LiquidityKeeper) IterateLiquidityPools(
	ctx sdk.Context, handlerFn func(liquidityPool LiquidityPool) (stop bool),
) {
	for _, l := range k.liquidityPool {
		if handlerFn(*l) {
			break
		}
	}
}

// Caculate total TVL
func (k LiquidityKeeper) CalculateTVL() sdk.Int {
	tvl := sdk.ZeroInt()
	for _, l := range k.liquidityPool {
		tvl = tvl.Add(l.TVL)
	}

	return tvl
}

// Caculate pool share using mulitplier
func (k LiquidityKeeper) CalculateProxyTVL() sdk.Dec {
	multipliedShareSum := sdk.ZeroDec()
	for _, l := range k.liquidityPool {
		proxyTVL := sdk.NewDecFromInt(l.TVL).MulInt64(l.multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)
	}

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum
}

// Calculate total TVL across all pools
func (k LiquidityKeeper) CalculateLiquidateTVL(ctx sdk.Context, delegator string) sdk.Int {
	liquidatedAmt := sdk.ZeroInt()

	for _, l := range k.lpDeposits {
		if l.ownerAddress != delegator {
			continue
		}
		// Calculate total pool share by TVL and multiplier
		liquidatedAmt = liquidatedAmt.Add(l.amount)
	}

	return liquidatedAmt
}

// After rewarding, reset the rewards that were given to avoid double paid.
func (k LiquidityKeeper) UpdateRewardsAccmulated(ctx sdk.Context) {
}
