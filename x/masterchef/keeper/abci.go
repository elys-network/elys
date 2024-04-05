package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

// EndBlocker of amm module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// distribute LP rewards
	k.ProcessLPRewardDistribution(ctx)
	// distribute external rewards
	k.ProcessExternalRewardsDistribution(ctx)
	// TODO: calculate APR (for external rewards)
}

func (k Keeper) ProcessExternalRewardsDistribution(ctx sdk.Context) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	lpIncentive := params.LpIncentives

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())

	externalIncentives := k.GetAllExternalIncentives(ctx)
	for _, externalIncentive := range externalIncentives {
		pool, found := k.GetPool(ctx, externalIncentive.PoolId)
		if !found {
			continue
		}

		if externalIncentive.FromBlock < curBlockHeight.Uint64() && curBlockHeight.Uint64() <= externalIncentive.ToBlock {
			k.UpdateAccPerShare(ctx, externalIncentive.PoolId, externalIncentive.RewardDenom, externalIncentive.AmountPerBlock)

			hasRewardDenom := false
			poolRewardDenoms := pool.ExternalRewardDenoms
			for _, poolRewardDenom := range poolRewardDenoms {
				if poolRewardDenom == externalIncentive.RewardDenom {
					hasRewardDenom = true
				}
			}
			if !hasRewardDenom {
				pool.ExternalRewardDenoms = append(pool.ExternalRewardDenoms, externalIncentive.RewardDenom)
				k.SetPool(ctx, pool)
			}

			ammPool, found := k.amm.GetPool(ctx, pool.PoolId)
			if found {
				tvl, err := ammPool.TVL(ctx, k.oracleKeeper)
				if err == nil {
					yearlyIncentiveRewardsTotal := externalIncentive.AmountPerBlock.
						Mul(lpIncentive.TotalBlocksPerYear).
						Quo(pool.NumBlocks)

					entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
					if found {
						baseCurrency := entry.Denom
						pool.ExternalIncentiveApr = sdk.NewDecFromInt(yearlyIncentiveRewardsTotal).
							Mul(k.GetTokenPrice(ctx, externalIncentive.RewardDenom, baseCurrency)).
							Quo(tvl)
						k.SetPool(ctx, pool)
					}
				}
			}

		}

		if curBlockHeight.Uint64() == externalIncentive.ToBlock {
			k.RemoveExternalIncentive(ctx, externalIncentive.Id)
		}
	}
}

func (k Keeper) ProcessLPRewardDistribution(ctx sdk.Context) {
	canDistribute := k.CanDistributeLPRewards(ctx)
	if canDistribute {
		err := k.UpdateLPRewardsUnclaimed(ctx)
		if err != nil {
			ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
		}
	}
}

func (k Keeper) CanDistributeLPRewards(ctx sdk.Context) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false
	}

	// If we don't have enough params
	if params.LpIncentives == nil {
		return false
	}

	// Incentive params initialize
	lpIncentive := params.LpIncentives

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())
	if lpIncentive.DistributionStartBlock.GT(curBlockHeight) {
		return false
	}

	// Increase current epoch of incentive param
	lpIncentive.CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks.Add(sdk.OneInt())
	if lpIncentive.CurrentEpochInBlocks.GTE(lpIncentive.TotalBlocksPerYear) || curBlockHeight.GT(lpIncentive.TotalBlocksPerYear.Add(lpIncentive.DistributionStartBlock)) {
		params.LpIncentives = nil
		k.SetParams(ctx, params)
		return false
	}

	params.LpIncentives.CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks
	k.SetParams(ctx, params)

	return true
}

// Update unclaimed token amount
// Called back through epoch hook
func (k Keeper) UpdateLPRewardsUnclaimed(ctx sdk.Context) error {
	// Fetch incentive params
	params := k.GetParams(ctx)
	lpIncentive := params.LpIncentives

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx, baseCurrency)

	// Collect DEX revenue while tracking 65% of it for LPs reward calculation
	// Assume these are collected in USDC
	_, dexRevenueForLpsPerDistribution := k.CollectDEXRevenue(ctx)

	// Calculate each portion of Gas fees collected - stakers, LPs
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(k.tci.TotalFeesCollected...)
	rewardPortionForLps := k.GetParams(ctx).RewardPortionForLps
	gasFeesForLpsPerDistribution := gasFeeCollectedDec.MulDecTruncate(rewardPortionForLps)

	// USDC amount in sdk.Dec type
	dexRevenueLPsAmtPerDistribution := dexRevenueForLpsPerDistribution.AmountOf(baseCurrency)
	gasFeesLPsAmtPerDistribution := gasFeesForLpsPerDistribution.AmountOf(baseCurrency)

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

	// Ensure lpIncentive.TotalBlocksPerYear or lpIncentive.EpochNumBlocks are not zero to avoid division by zero
	if lpIncentive.TotalBlocksPerYear.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate eden amount per epoch
	epochLpsEdenAmount := lpIncentive.EdenAmountPerYear.
		Quo(lpIncentive.TotalBlocksPerYear)

	// Maximum eden based per distribution epoch on maximum APR - 30% by default
	// Allocated for staking per day = (0.3/365)* (total weighted proxy TVL)
	edenDenomPrice := k.GetEdenDenomPrice(ctx, baseCurrency)

	// Ensure edenDenomPrice is not zero to avoid division by zero
	if edenDenomPrice.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid eden price")
	}

	// TODO: apply update on incentive module
	epochLpsMaxEdenAmount := params.MaxEdenRewardAprLps.
		Mul(totalProxyTVL).
		QuoInt(lpIncentive.TotalBlocksPerYear).
		Quo(edenDenomPrice)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	epochLpsEdenAmount = sdk.MinInt(epochLpsEdenAmount, epochLpsMaxEdenAmount.TruncateInt())

	stableStakePoolId := uint64(stabletypes.PoolId)
	// Distribute Eden / USDC Rewards
	for _, pool := range k.GetAllPools(ctx) {
		var proxyTVL sdk.Dec
		if pool.PoolId == stableStakePoolId {
			tvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
			proxyTVL = tvl.Mul(pool.Multiplier)
		} else {
			ammPool, found := k.amm.GetPool(ctx, pool.PoolId)
			if !found {
				continue
			}

			// ------------ New Eden calculation -------------------
			// -----------------------------------------------------
			// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
			// Pool share = 80
			// edenAmountPerEpochLp = 100
			tvl, err := ammPool.TVL(ctx, k.oracleKeeper)
			if err != nil {
				continue
			}
			// Calculate Proxy TVL share considering multiplier
			proxyTVL = tvl.Mul(pool.Multiplier)
		}

		poolShare := sdk.ZeroDec()
		if totalProxyTVL.IsPositive() {
			poolShare = proxyTVL.Quo(totalProxyTVL)
		}

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShare.MulInt(epochLpsEdenAmount)

		// Get gas fee rewards per pool
		gasRewardsAllocatedForPool := poolShare.Mul(gasFeesLPsAmtPerDistribution)

		// ------------------- DEX rewards calculation -------------------
		// ---------------------------------------------------------------
		// Get dex rewards per pool
		// Get track key
		trackKey := types.GetPoolRevenueTrackKey(pool.PoolId)
		// Get tracked amount for Lps per pool
		dexRewardsAllocatedForPool, ok := k.tci.PoolRevenueTrack[trackKey]
		if !ok {
			dexRewardsAllocatedForPool = sdk.NewDec(0)
		}

		// Distribute Eden
		k.UpdateAccPerShare(ctx, pool.PoolId, ptypes.Eden, newEdenAllocatedForPool.TruncateInt())
		// Distribute Gas fees + Dex rewards (USDC)
		k.UpdateAccPerShare(ctx, pool.PoolId, k.GetBaseCurrencyDenom(ctx), gasRewardsAllocatedForPool.Add(dexRewardsAllocatedForPool).TruncateInt())

		// Update Pool Info
		pool.EdenRewardAmountGiven = newEdenAllocatedForPool.RoundInt()
		pool.DexRewardAmountGiven = gasRewardsAllocatedForPool.Add(dexRewardsAllocatedForPool)
		k.SetPool(ctx, pool)
	}

	// Increase block number
	params.DexRewardsLps.NumBlocks = sdk.OneInt()
	// Incrase total dex rewards given
	params.DexRewardsLps.Amount = dexRevenueLPsAmtPerDistribution.Add(gasFeesLPsAmtPerDistribution)
	k.SetParams(ctx, params)

	// Update APR for amm pools
	k.UpdateAmmPoolAPR(ctx, lpIncentive.TotalBlocksPerYear, totalProxyTVL, edenDenomPrice)

	return nil
}

// Update total commitment info
// TODO: should be removed or updated for complexity in calculation
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context, baseCurrency string) {
	// Initialize with amount zero
	k.tci.TotalFeesCollected = sdk.Coins{}
	// Initialize Lp tokens amount
	k.tci.TotalLpTokensCommitted = make(map[string]math.Int)
	// Reinitialize Pool revenue tracker
	k.tci.PoolRevenueTrack = make(map[string]sdk.Dec)

	// Collect gas fees collected
	fees := k.CollectGasFeesToIncentiveModule(ctx, baseCurrency)

	// Calculate total fees - Gas fees collected
	k.tci.TotalFeesCollected = k.tci.TotalFeesCollected.Add(fees...)

	// Iterate to calculate total Eden, Eden boost and Lp tokens committed
	k.cmk.IterateCommitments(ctx, func(commitments ctypes.Commitments) bool {
		// Iterate to calculate total Lp tokens committed
		k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
			lpToken := ammtypes.GetPoolShareDenom(p.GetPoolId())

			committedLpToken := commitments.GetCommittedAmountForDenom(lpToken)
			amt, ok := k.tci.TotalLpTokensCommitted[lpToken]
			if !ok {
				k.tci.TotalLpTokensCommitted[lpToken] = committedLpToken
			} else {
				k.tci.TotalLpTokensCommitted[lpToken] = amt.Add(committedLpToken)
			}
			return false
		})

		// handle stable stake pool lp token
		lpStableStakeDenom := stabletypes.GetShareDenom()
		committedLpToken := commitments.GetCommittedAmountForDenom(lpStableStakeDenom)
		amt, ok := k.tci.TotalLpTokensCommitted[lpStableStakeDenom]
		if !ok {
			k.tci.TotalLpTokensCommitted[lpStableStakeDenom] = committedLpToken
		} else {
			k.tci.TotalLpTokensCommitted[lpStableStakeDenom] = amt.Add(committedLpToken)
		}
		return false
	})
}

// Move gas fees collected to dex revenue wallet
// Convert it into USDC
func (k Keeper) CollectGasFeesToIncentiveModule(ctx sdk.Context, baseCurrency string) sdk.Coins {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollected := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	// Total Swapped coin
	totalSwappedCoins := sdk.Coins{}

	for _, tokenIn := range feesCollected {
		// if it is base currency - usdc, we don't need convert. We just need to collect it to fee wallet.
		if tokenIn.Denom == baseCurrency {
			// Transfer converted USDC fees to the Dex revenue module account
			err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, k.dexRevCollectorName, sdk.Coins{tokenIn})
			if err != nil {
				panic(err)
			}

			// Sum total swapped
			totalSwappedCoins = totalSwappedCoins.Add(tokenIn)
			continue
		}

		// Find a pool that can convert tokenIn to usdc
		pool, found := k.amm.GetBestPoolWithDenoms(ctx, []string{tokenIn.Denom, baseCurrency})
		if !found {
			continue
		}

		// Executes the swap in the pool and stores the output. Updates pool assets but
		// does not actually transfer any tokens to or from the pool.
		snapshot := k.amm.GetPoolSnapshotOrSet(ctx, pool)
		tokenOutCoin, _, _, _, err := k.amm.SwapOutAmtGivenIn(ctx, pool.PoolId, k.oracleKeeper, &snapshot, sdk.Coins{tokenIn}, baseCurrency, sdk.ZeroDec())
		if err != nil {
			continue
		}

		tokenOutAmount := tokenOutCoin.Amount
		if !tokenOutAmount.IsPositive() {
			continue
		}

		// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
		// Also emits a swap event and updates related liquidity metrics.
		cacheCtx, write := ctx.CacheContext()
		_, err = k.amm.UpdatePoolForSwap(cacheCtx, pool, feeCollector.GetAddress(), feeCollector.GetAddress(), tokenIn, tokenOutCoin, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
		if err != nil {
			continue
		}
		write()

		// Swapped USDC coin
		swappedCoins := sdk.NewCoins(sdk.NewCoin(baseCurrency, tokenOutAmount))

		// Transfer converted USDC fees to the Dex revenue module account
		if swappedCoins.IsAllPositive() {
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, k.dexRevCollectorName, swappedCoins)
			if err != nil {
				panic(err)
			}
		}

		// Sum total swapped
		totalSwappedCoins = totalSwappedCoins.Add(swappedCoins...)
	}

	return totalSwappedCoins
}

// Collect all DEX revenues to DEX revenue wallet,
// while tracking the 65% of it for LPs reward distribution
// transfer collected fees from different wallets(liquidity pool, perpetual module etc) to the distribution module account
// Assume this is already in USDC.
// TODO:
// + Collect revenue from perpetual, lend module
func (k Keeper) CollectDEXRevenue(ctx sdk.Context) (sdk.Coins, sdk.DecCoins) {
	// Total colllected revenue amount
	amountTotalCollected := sdk.Coins{}
	amountLPsCollected := sdk.DecCoins{}

	// Iterate to calculate total Eden from LpElys, MElys committed
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		// Get pool Id
		poolId := p.GetPoolId()

		// Get dex rewards per pool
		revenueAddress := ammtypes.NewPoolRevenueAddress(poolId)

		// Transfer revenue to a single wallet of DEX revenue wallet.
		revenue := k.bankKeeper.GetAllBalances(ctx, revenueAddress)
		if revenue.IsAllPositive() {
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, revenueAddress, k.dexRevCollectorName, revenue)
			if err != nil {
				panic(err)
			}
		}

		// LPs Portion param
		rewardPortionForLps := k.GetParams(ctx).RewardPortionForLps

		// Calculate revenue portion for LPs
		revenueDec := sdk.NewDecCoinsFromCoins(revenue...)

		// LPs portion of pool revenue
		revenuePortionForLPs := revenueDec.MulDecTruncate(rewardPortionForLps)

		// Get track key
		trackKey := types.GetPoolRevenueTrackKey(poolId)

		// Store revenue portion for Lps temporarilly
		k.tci.PoolRevenueTrack[trackKey] = revenuePortionForLPs.AmountOf(ptypes.BaseCurrency)

		// Sum total collected amount
		amountTotalCollected = amountTotalCollected.Add(revenue...)

		// Sum total amount for LPs
		amountLPsCollected = amountLPsCollected.Add(revenuePortionForLPs...)

		return false
	})

	return amountTotalCollected, amountLPsCollected
}

// Calculate Proxy TVL
func (k Keeper) CalculateProxyTVL(ctx sdk.Context, baseCurrency string) sdk.Dec {
	multipliedShareSum := sdk.ZeroDec()
	stableStakePoolId := uint64(stabletypes.PoolId)
	for _, pool := range k.GetAllPools(ctx) {
		if pool.PoolId == stableStakePoolId {
			// Get pool info from incentive param
			poolInfo, found := k.GetPool(ctx, stableStakePoolId)
			if !found {
				k.InitStableStakePoolParams(ctx, stableStakePoolId)
				poolInfo, _ = k.GetPool(ctx, stableStakePoolId)
			}
			tvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
			proxyTVL := tvl.Mul(poolInfo.Multiplier)
			multipliedShareSum = multipliedShareSum.Add(proxyTVL)
			continue
		}

		ammPool, found := k.amm.GetPool(ctx, pool.PoolId)
		if !found {
			continue
		}

		tvl, err := ammPool.TVL(ctx, k.oracleKeeper)
		if err != nil {
			continue
		}

		// Get pool info from incentive param
		poolInfo, found := k.GetPool(ctx, ammPool.GetPoolId())
		if !found {
			continue
		}

		proxyTVL := tvl.Mul(poolInfo.Multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)
	}

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum
}

// InitPoolParams: creates a poolInfo at the time of pool creation.
func (k Keeper) InitPoolParams(ctx sdk.Context, poolId uint64) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	poolInfos := params.PoolInfos

	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return true
		}
	}

	// Initiate a new pool info
	poolInfo := types.PoolInfo{
		// reward amount
		PoolId: poolId,
		// reward wallet address
		RewardWallet: ammtypes.NewPoolRevenueAddress(poolId).String(),
		// multiplier for lp rewards
		Multiplier: sdk.NewDec(1),
		// Number of blocks since creation
		NumBlocks: sdk.NewInt(1),
		// Total dex rewards given since creation
		DexRewardAmountGiven: sdk.ZeroDec(),
		// Total eden rewards given since creation
		EdenRewardAmountGiven: sdk.ZeroInt(),
	}

	// Update pool information
	params.PoolInfos = append(params.PoolInfos, poolInfo)
	k.SetParams(ctx, params)

	return true
}

// InitStableStakePoolMultiplier: create a stable stake pool information responding to the pool creation.
func (k Keeper) InitStableStakePoolParams(ctx sdk.Context, poolId uint64) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	poolInfos := params.PoolInfos

	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return true
		}
	}

	// Initiate a new pool info
	poolInfo := types.PoolInfo{
		// reward amount
		PoolId: poolId,
		// reward wallet address
		RewardWallet: stabletypes.PoolAddress().String(),
		// multiplier for lp rewards
		Multiplier: sdk.NewDec(1),
		// Number of blocks since creation
		NumBlocks: sdk.NewInt(1),
		// Total dex rewards given since creation
		DexRewardAmountGiven: sdk.ZeroDec(),
		// Total eden rewards given since creation
		EdenRewardAmountGiven: sdk.ZeroInt(),
	}

	// Update pool information
	params.PoolInfos = append(params.PoolInfos, poolInfo)
	k.SetParams(ctx, params)

	return true
}

// Calculate new Eden token amounts based on LpElys committed and MElys committed
func (k Keeper) CalcRewardsForLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, commitments ctypes.Commitments, edenAmountForLpPerDistribution math.Int, gasFeesForLPsPerDistribution sdk.Dec) (math.Int, math.Int) {
	// Method 2 - Using Proxy TVL
	totalNewEdenAllocatedPerDistribution := sdk.ZeroInt()
	totalDexRewardsAllocatedPerDistribution := sdk.ZeroDec()

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
		poolInfo, found := k.GetPool(ctx, poolId)
		if !found {
			k.InitPoolParams(ctx, poolId)
			poolInfo, _ = k.GetPool(ctx, poolId)
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := tvl.Mul(poolInfo.Multiplier)
		poolShare := sdk.ZeroDec()
		if totalProxyTVL.IsPositive() {
			poolShare = proxyTVL.Quo(totalProxyTVL)
		}

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShare.MulInt(edenAmountForLpPerDistribution)

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
		totalNewEdenAllocatedPerDistribution = totalNewEdenAllocatedPerDistribution.Add(newEdenAllocated)
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
		totalDexRewardsAllocatedPerDistribution = totalDexRewardsAllocatedPerDistribution.Add(dexRewardsForLP)

		//----------------------------------------------------------------

		// ------------------- Gas rewards calculation -------------------
		// ---------------------------------------------------------------
		// Get gas fee rewards per pool
		gasRewardsAllocatedForPool := poolShare.Mul(gasFeesForLPsPerDistribution)
		// Calculate gas fee rewards per lp
		gasRewardsForLP := lpShare.Mul(gasRewardsAllocatedForPool)
		// Sum total rewards per commitment
		totalDexRewardsAllocatedPerDistribution = totalDexRewardsAllocatedPerDistribution.Add(gasRewardsForLP)

		//----------------------------------------------------------------

		poolInfo.EdenRewardAmountGiven = newEdenAllocatedForPool.RoundInt()
		poolInfo.DexRewardAmountGiven = gasRewardsAllocatedForPool.Add(dexRewardsAllocatedForPool)
		// Update Pool Info
		k.SetPool(ctx, poolInfo)

		return false
	})

	// return
	return totalNewEdenAllocatedPerDistribution, totalDexRewardsAllocatedPerDistribution.TruncateInt()
}

// Calculate pool share for stable stake pool
func (k Keeper) CalculatePoolShareForStableStakeLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, baseCurrency string) sdk.Dec {
	// ------------ New Eden calculation -------------------
	// -----------------------------------------------------
	// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
	// Pool share = 80
	// edenAmountPerEpochLp = 100
	tvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)

	// Get pool Id
	poolId := uint64(stabletypes.PoolId)

	// Get pool info from incentive param
	poolInfo, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.ZeroDec()
	}

	// Calculate Proxy TVL share considering multiplier
	proxyTVL := tvl.Mul(poolInfo.Multiplier)
	if totalProxyTVL.IsZero() {
		return sdk.ZeroDec()
	}
	poolShare := proxyTVL.Quo(totalProxyTVL)

	return poolShare
}

// Calculate new Eden token amounts based on LpElys committed and MElys committed
func (k Keeper) CalcRewardsForStableStakeLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, commitments ctypes.Commitments, edenAmountPerEpochLp math.Int, gasFeesForLPs sdk.Dec, baseCurrency string) (math.Int, math.Int) {
	// Method 2 - Using Proxy TVL
	totalDexRewardsAllocated := sdk.ZeroDec()

	// Calculate pool share for stable stake pool
	poolShare := k.CalculatePoolShareForStableStakeLPs(ctx, totalProxyTVL, baseCurrency)

	// Calculate new Eden for this pool
	newEdenAllocatedForPool := poolShare.MulInt(edenAmountPerEpochLp)

	// Get pool share denom - stable stake lp token
	lpToken := stabletypes.GetShareDenom()

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
	// Get pool Id
	poolId := uint64(stabletypes.PoolId)
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

func (k Keeper) UpdateCommitments(
	ctx sdk.Context,
	creator string,
	commitments *ctypes.Commitments,
	newUnclaimedEdenTokens math.Int,
	newUnclaimedEdenBTokens math.Int,
	dexRewards math.Int,
	baseCurrency string,
) {
	// Update unclaimed Eden balances in the Commitments structure
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
	// Update unclaimed Eden-Boost token balances in the Commitments structure
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBTokens))

	// All dex revenue are collected to incentive module in USDC
	// Gas fees (Elys) are also converted into USDC and collected into total dex revenue wallet of incentive module.
	// Update USDC balances in the Commitments structure.
	// These are the rewards from each pool, perpetual, gas fee.
	commitments.AddRewardsUnclaimed(sdk.NewCoin(baseCurrency, dexRewards))

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update sub bucket commitment record
func (k Keeper) UpdateCommitmentsSubBuckets(ctx sdk.Context, creator string, commitments *ctypes.Commitments, rewardsByElysStaking sdk.Coins, rewardsByEdenCommitted sdk.Coins, rewardsByEdenBCommitted sdk.Coins, rewardsByUSDCDeposit sdk.Coins) {
	// Add to Elys staking bucket
	commitments.AddSubBucketRewardsByElysUnclaimed(rewardsByElysStaking)
	// Add to Eden committed bucket
	commitments.AddSubBucketRewardsByEdenUnclaimed(rewardsByEdenCommitted)
	// Add to EdenB committed bucket
	commitments.AddSubBucketRewardsByEdenBUnclaimed(rewardsByEdenBCommitted)
	// Add to USDC deposit bucket
	commitments.AddSubBucketRewardsByUsdcUnclaimed(rewardsByUSDCDeposit)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update APR for AMM pool
func (k Keeper) UpdateAmmPoolAPR(ctx sdk.Context, totalBlocksPerYear sdk.Int, totalProxyTVL sdk.Dec, edenDenomPrice sdk.Dec) {
	// Iterate to calculate total Eden from LpElys, MElys committed
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool Id
		poolId := p.GetPoolId()

		// Get pool info from incentive param
		poolInfo, found := k.GetPool(ctx, poolId)
		if !found {
			k.InitPoolParams(ctx, poolId)
			poolInfo, _ = k.GetPool(ctx, poolId)
		}

		poolInfo.NumBlocks = sdk.OneInt()
		// Invalid block number
		if poolInfo.NumBlocks.IsZero() {
			return false
		}

		if tvl.IsZero() {
			return false
		}

		// Dex reward Apr per pool =  total accumulated usdc rewards for 7 day * 52/ tvl of pool
		yearlyDexRewardsTotal := poolInfo.DexRewardAmountGiven.
			MulInt(totalBlocksPerYear).
			QuoInt(poolInfo.NumBlocks)
		poolInfo.DexApr = yearlyDexRewardsTotal.
			Quo(tvl)

		// Eden reward Apr per pool = (total LM Eden reward allocated per day*((tvl of pool * multiplier)/total proxy TVL) ) * 365 / TVL of pool
		yearlyEdenRewardsTotal := poolInfo.EdenRewardAmountGiven.
			Mul(totalBlocksPerYear).
			Quo(poolInfo.NumBlocks)

		poolInfo.EdenApr = sdk.NewDecFromInt(yearlyEdenRewardsTotal).
			Mul(edenDenomPrice).
			Quo(tvl)

		// Update Pool Info
		k.SetPool(ctx, poolInfo)

		return false
	})
}
