package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

// Calculate new Eden token amounts based on LpElys committed and MElys committed
func (k Keeper) CalculateRewardsForStableStakeLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, commitments ctypes.Commitments, edenAmountPerEpochLp sdk.Int, gasFeesForLPs sdk.Dec, baseCurrency string) (sdk.Int, sdk.Int) {
	// Method 2 - Using Proxy TVL
	totalDexRewardsAllocated := sdk.ZeroDec()

	// ------------ New Eden calculation -------------------
	// -----------------------------------------------------
	// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
	// Pool share = 80
	// edenAmountPerEpochLp = 100
	tvl := stabletypes.TVL(ctx, k.authKeeper, k.bankKeeper, baseCurrency)

	// Get pool Id
	poolId := uint64(stabletypes.PoolId)

	// Get pool share denom - stable stake lp token
	lpToken := stabletypes.GetShareDenom()

	// Get pool info from incentive param
	poolInfo, found := k.GetPoolInfo(ctx, poolId)
	if !found {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// Calculate Proxy TVL share considering multiplier
	proxyTVL := sdk.NewDecFromInt(tvl).Mul(poolInfo.Multiplier)
	poolShare := proxyTVL.Quo(totalProxyTVL)

	// Calculate new Eden for this pool
	newEdenAllocatedForPool := poolShare.MulInt(edenAmountPerEpochLp)

	// this lp token committed
	commmittedLpToken := commitments.GetCommittedAmountForDenom(lpToken)
	// this lp token total committed
	totalCommittedLpToken, ok := k.tci.TotalLpTokensCommitted[lpToken]
	if !ok {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// If total committed LP token amount is zero
	if totalCommittedLpToken.LTE(sdk.ZeroInt()) {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// Calculalte lp token share of the pool
	lpShare := sdk.NewDecFromInt(commmittedLpToken).QuoInt(totalCommittedLpToken)

	// Calculate new Eden allocated per LP
	newEdenAllocated := lpShare.Mul(newEdenAllocatedForPool).TruncateInt()

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

	// return
	return newEdenAllocated, totalDexRewardsAllocated.TruncateInt()
}
