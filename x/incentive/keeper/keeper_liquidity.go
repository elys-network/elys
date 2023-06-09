package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------Placeholder implementation---------
// --------------------------------------------
// Liquidity pool
type LiquidityPool struct {
	index       string
	TVL         sdk.Int
	multiplier  int64
	lpToken     string
	poolRewards sdk.Dec
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

// Add dummy liquidity pool
func (k LiquidityKeeper) AddLiquidityPool(index string, TVL sdk.Int, multiplier int64, lpToken string, poolRewards sdk.Dec) {
	lp := &LiquidityPool{
		index:       index,
		TVL:         TVL,
		multiplier:  multiplier,
		lpToken:     lpToken,
		poolRewards: poolRewards,
	}

	k.liquidityPool[index] = lp
}

// Add dummy liquidity pool
func (k LiquidityKeeper) DepositTokenToLP(index string, LpIndex string, amount sdk.Int, ownerAddress string, denom string) {
	lp, ok := k.liquidityPool[LpIndex]
	if !ok {
		return
	}

	lpDeposit := LpDeposit{
		index:        index,
		LpIndex:      LpIndex,
		amount:       amount,
		ownerAddress: ownerAddress,
		denom:        denom,
	}

	lp.TVL = lp.TVL.Add(amount)
	k.lpDeposits = append(k.lpDeposits, lpDeposit)
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

// Draft function
// TODO:
// + After we distribute rewards from each pool, we should reset the amount of rewards accumulated
// + in each pool in order to avoid double paying.
// + Regarding reward wallet topic, I still need to discuss with team members. So draft this part.
func (k LiquidityKeeper) UpdateRewardsAccmulated(ctx sdk.Context) {
}
