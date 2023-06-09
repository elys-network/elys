package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
)

// Calculate new Eden token amounts based on LpElys committed and MElys committed
func (k Keeper) CalculateRewardsForLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, commitments ctypes.Commitments, edenAmountPerEpochLp sdk.Int) (sdk.Int, sdk.Int) {
	// Method 2 - Using Proxy TVL
	totalNewEdenAllocated := sdk.ZeroInt()
	totalDexRewardsAllocated := sdk.ZeroDec()

	// Iterate to calculate total Eden from LpElys, MElys committed
	k.Lpk.IterateLiquidityPools(ctx, func(l LiquidityPool) bool {
		// ------------ New Eden calculation -------------------
		// -----------------------------------------------------
		// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
		// Pool share = 80
		// edenAmountPerEpochLp = 100

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := sdk.NewDecFromInt(l.TVL).MulInt64(l.multiplier)
		poolShare := proxyTVL.Quo(totalProxyTVL)

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShare.MulInt(edenAmountPerEpochLp)

		// this lp token committed
		commmittedLpToken := commitments.GetCommittedAmountForDenom(l.lpToken)
		// this lp token total committed
		totalCommittedLpToken, ok := k.tci.TotalLpTokensCommitted[l.lpToken]
		if !ok {
			return false
		}

		// If total committed LP token amount is zero
		if totalCommittedLpToken.LTE(sdk.ZeroInt()) {
			return false
		}

		// Calculalte lp token share of the pool
		lpShare := sdk.NewDecFromInt(commmittedLpToken).QuoInt(totalCommittedLpToken)

		// Calculate new Eden allocated per LP
		newEdenAllocated := lpShare.Mul(newEdenAllocatedForPool).TruncateInt()

		// Sum the total amount
		totalNewEdenAllocated = totalNewEdenAllocated.Add(newEdenAllocated)
		// -------------------------------------------------------

		// ------------------- DEX rewards calculation -------------------
		// ---------------------------------------------------------------
		// Get dex rewards per pool
		dexRewardsAllocatedForPool := l.poolRewards
		// Calculate dex rewards per lp
		dexRewardsForLP := lpShare.Mul(dexRewardsAllocatedForPool)
		// Sum total rewards per commitment
		totalDexRewardsAllocated = totalDexRewardsAllocated.Add(dexRewardsForLP)

		//----------------------------------------------------------------
		return false
	})

	// return
	return totalNewEdenAllocated, totalDexRewardsAllocated.TruncateInt()
}
