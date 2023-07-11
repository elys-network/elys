package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
)

// Calculate new Eden token amounts based on LpElys committed and MElys committed
func (k Keeper) CalculateRewardsForLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, commitments ctypes.Commitments, edenAmountPerEpochLp sdk.Int, gasFeesForLPs sdk.Dec) (sdk.Int, sdk.Int) {
	// Method 2 - Using Proxy TVL
	totalNewEdenAllocated := sdk.ZeroInt()
	totalDexRewardsAllocated := sdk.ZeroDec()

	// Iterate to calculate total Eden from LpElys, MElys committed
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		// ------------ New Eden calculation -------------------
		// -----------------------------------------------------
		// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
		// Pool share = 80
		// edenAmountPerEpochLp = 100
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool Id
		poolId := p.GetPoolId()

		// Get pool share denom - lp token
		lpToken := ammtypes.GetPoolShareDenom(poolId)

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, poolId)
		if !found {
			return false
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := tvl.Mul(poolInfo.Multiplier)
		poolShare := proxyTVL.Quo(totalProxyTVL)

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShare.MulInt(edenAmountPerEpochLp)

		// this lp token committed
		commmittedLpToken := commitments.GetCommittedAmountForDenom(lpToken)
		// this lp token total committed
		totalCommittedLpToken, ok := k.tci.TotalLpTokensCommitted[lpToken]
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
		// Get track key
		trackKey := types.GetPoolRevenueTrackKey(poolId)
		// Get tracked amount for Lps per pool
		dexRewardsAllocatedForPool, ok := k.tci.PoolRevenueTrack[trackKey]
		if !ok {
			dexRewardsAllocatedForPool = sdk.NewDec(0)
		}

		// Calculate dex rewards per lp
		dexRewardsForLP := lpShare.Mul(dexRewardsAllocatedForPool)
		// Sum total rewards per commitment
		totalDexRewardsAllocated = totalDexRewardsAllocated.Add(dexRewardsForLP)

		//----------------------------------------------------------------

		// ------------------- Gas rewards calculation -------------------
		// ---------------------------------------------------------------
		// Get gas fee rewards per pool
		gasRewardsAllocatedForPool := poolShare.Mul(gasFeesForLPs)
		// Calculate gas fee rewards per lp
		gasRewardsForLP := lpShare.Mul(gasRewardsAllocatedForPool)
		// Sum total rewards per commitment
		totalDexRewardsAllocated = totalDexRewardsAllocated.Add(gasRewardsForLP)

		//----------------------------------------------------------------
		return false
	})

	// return
	return totalNewEdenAllocated, totalDexRewardsAllocated.TruncateInt()
}
