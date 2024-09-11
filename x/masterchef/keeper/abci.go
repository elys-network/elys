package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

// EndBlocker of amm module
func (k Keeper) EndBlocker(ctx sdk.Context) {

	k.DeleteFeeInfo(ctx)

	// distribute LP rewards
	k.ProcessLPRewardDistribution(ctx)
	// distribute external rewards
	k.ProcessExternalRewardsDistribution(ctx)
}

func (k Keeper) GetPoolTVL(ctx sdk.Context, poolId uint64) math.LegacyDec {
	if poolId == stabletypes.PoolId {
		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return math.LegacyZeroDec()
		}
		return k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
	}
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if found {
		tvl, err := ammPool.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return math.LegacyZeroDec()
		}
		return tvl
	}
	return math.LegacyZeroDec()
}

func (k Keeper) ProcessExternalRewardsDistribution(ctx sdk.Context) {
	baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
	curBlockHeight := ctx.BlockHeight()
	totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear

	externalIncentives := k.GetAllExternalIncentives(ctx)
	externalIncentiveAprs := make(map[uint64]math.LegacyDec)
	for _, externalIncentive := range externalIncentives {
		pool, found := k.GetPoolInfo(ctx, externalIncentive.PoolId)
		if !found {
			continue
		}

		if externalIncentive.FromBlock < curBlockHeight && curBlockHeight <= externalIncentive.ToBlock {
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
				k.SetPoolInfo(ctx, pool)
			}

			tvl := k.GetPoolTVL(ctx, pool.PoolId)
			if tvl.IsPositive() {
				yearlyIncentiveRewardsTotal := externalIncentive.AmountPerBlock.
					Mul(sdk.NewInt(totalBlocksPerYear))

				apr := sdk.NewDecFromInt(yearlyIncentiveRewardsTotal).
					Mul(k.amm.GetTokenPrice(ctx, externalIncentive.RewardDenom, baseCurrency)).
					Quo(tvl)
				externalIncentive.Apr = apr
				k.SetExternalIncentive(ctx, externalIncentive)
				poolExternalApr, ok := externalIncentiveAprs[pool.PoolId]
				if !ok {
					poolExternalApr = math.LegacyZeroDec()
				}

				poolExternalApr = poolExternalApr.Add(apr)
				externalIncentiveAprs[pool.PoolId] = poolExternalApr
				pool.ExternalIncentiveApr = poolExternalApr
				k.SetPoolInfo(ctx, pool)
			}
		}

		if curBlockHeight == externalIncentive.ToBlock {
			k.RemoveExternalIncentive(ctx, externalIncentive.Id)
		}
	}
}

func (k Keeper) ProcessLPRewardDistribution(ctx sdk.Context) {
	// Read tokenomics time based inflation params and update incentive module params.
	k.ProcessUpdateIncentiveParams(ctx)

	err := k.UpdateLPRewards(ctx)
	if err != nil {
		ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
	}
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) {
	// Non-linear inflation per year happens and this includes yearly inflation data
	listTimeBasedInflations := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflations) == 0 {
		return
	}

	params := k.GetParams(ctx)
	for _, inflation := range listTimeBasedInflations {
		// Finding only current inflation data - and skip rest
		if inflation.StartBlockHeight > uint64(ctx.BlockHeight()) || inflation.EndBlockHeight < uint64(ctx.BlockHeight()) {
			continue
		}

		totalBlocks := inflation.EndBlockHeight - inflation.StartBlockHeight + 1

		// If totalBlocks is zero, we skip this inflation to avoid division by zero
		if totalBlocks == 0 {
			continue
		}
		blocksDistributed := ctx.BlockHeight() - int64(inflation.StartBlockHeight)

		params.LpIncentives = &types.IncentiveInfo{
			// reward amount in eden for 1 year
			EdenAmountPerYear: sdk.NewInt(int64(inflation.Inflation.LmRewards)),
			// number of blocks distributed
			BlocksDistributed: blocksDistributed,
		}
		k.SetParams(ctx, params)
		return
	}

	params.LpIncentives = nil
	k.SetParams(ctx, params)
}

func (k Keeper) UpdateLPRewards(ctx sdk.Context) error {
	// Fetch incentive params
	params := k.GetParams(ctx)
	lpIncentive := params.LpIncentives

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// Collect Gas fees + swap fees
	gasFeesForLpsDec := k.CollectGasFees(ctx, baseCurrency)
	_, _, rewardsPerPool := k.CollectDEXRevenue(ctx)

	// USDC amount in sdk.Dec type
	gasFeeUsdcAmountForLps := gasFeesForLpsDec.AmountOf(baseCurrency)

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have multiplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL, totalProxyTvlEdenEnable := k.CalculateProxyTVL(ctx, baseCurrency)

	// Ensure totalBlocksPerYear is not zero to avoid division by zero
	totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear
	if totalBlocksPerYear == 0 {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate eden amount per block
	edenAmountPerYear := sdk.ZeroInt()
	if lpIncentive != nil && lpIncentive.EdenAmountPerYear.IsPositive() {
		edenAmountPerYear = lpIncentive.EdenAmountPerYear
	}
	lpsEdenAmount := edenAmountPerYear.Quo(sdk.NewInt(totalBlocksPerYear))

	// Ensure edenDenomPrice is not zero to avoid division by zero
	edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
	if edenDenomPrice.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid eden price")
	}

	// Distribute Eden / USDC Rewards
	for _, pool := range k.GetAllPoolInfos(ctx) {
		var err error
		tvl := k.GetPoolTVL(ctx, pool.PoolId)
		proxyTVL := tvl.Mul(pool.Multiplier)
		if proxyTVL.IsZero() {
			continue
		}

		poolShare := sdk.ZeroDec()
		poolShareEdenEnable := sdk.ZeroDec()
		if totalProxyTVL.IsPositive() {
			poolShare = proxyTVL.Quo(totalProxyTVL)
		}

		if totalProxyTvlEdenEnable.IsPositive() {
			poolShareEdenEnable = proxyTVL.Quo(totalProxyTvlEdenEnable)
		}

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShareEdenEnable.MulInt(lpsEdenAmount)

		// Maximum eden APR - 30% by default
		poolMaxEdenAmount := params.MaxEdenRewardAprLps.
			Mul(proxyTVL).
			QuoInt64(totalBlocksPerYear).
			Quo(edenDenomPrice)

		// Use min amount (eden allocation from tokenomics and max apr based eden amount)
		if pool.EnableEdenRewards {
			newEdenAllocatedForPool = sdk.MinDec(newEdenAllocatedForPool, poolMaxEdenAmount)
			if newEdenAllocatedForPool.IsPositive() {
				err = k.commitmentKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ptypes.Eden, newEdenAllocatedForPool.TruncateInt())})
				if err != nil {
					return err
				}
			}
		}

		// Get gas fee rewards per pool
		gasRewardsAllocatedForPool := poolShare.Mul(gasFeeUsdcAmountForLps)

		// ------------------- DEX rewards calculation -------------------
		// ---------------------------------------------------------------
		// Get dex rewards per pool
		// Get tracked amount for Lps per pool
		dexRewardsAllocatedForPool, ok := rewardsPerPool[pool.PoolId]
		if !ok {
			dexRewardsAllocatedForPool = sdk.NewDec(0)
		}

		k.AddEdenInfo(ctx, newEdenAllocatedForPool)

		// Distribute Eden
		if pool.EnableEdenRewards {
			k.UpdateAccPerShare(ctx, pool.PoolId, ptypes.Eden, newEdenAllocatedForPool.TruncateInt())
		}
		// Distribute Gas fees + Dex rewards (USDC)
		k.UpdateAccPerShare(ctx, pool.PoolId, k.GetBaseCurrencyDenom(ctx), gasRewardsAllocatedForPool.Add(dexRewardsAllocatedForPool).TruncateInt())

		// Track pool rewards accumulation
		edenReward := newEdenAllocatedForPool
		if !pool.EnableEdenRewards {
			edenReward = sdk.ZeroDec()
		}

		k.AddPoolRewardsAccum(
			ctx,
			pool.PoolId,
			uint64(ctx.BlockTime().Unix()),
			ctx.BlockHeight(),
			dexRewardsAllocatedForPool,
			gasRewardsAllocatedForPool,
			edenReward,
		)
		params := k.parameterKeeper.GetParams(ctx)
		dataLifetime := params.RewardsDataLifetime
		for {
			firstAccum := k.FirstPoolRewardsAccum(ctx, pool.PoolId)
			if firstAccum.Timestamp == 0 || int64(firstAccum.Timestamp)+dataLifetime >= ctx.BlockTime().Unix() {
				break
			}
			k.DeletePoolRewardsAccum(ctx, firstAccum)
		}

		if pool.EnableEdenRewards {
			pool.EdenApr = newEdenAllocatedForPool.
				MulInt64(totalBlocksPerYear).
				Mul(edenDenomPrice).
				Quo(tvl)
		} else {
			pool.EdenApr = sdk.ZeroDec()
		}

		k.SetPoolInfo(ctx, pool)
	}

	// Update APR for amm pools
	k.UpdateAmmPoolAPR(ctx, totalBlocksPerYear, totalProxyTVL, edenDenomPrice)

	return nil
}

// Move gas fees collected to dex revenue wallet
// Convert it into USDC
func (k Keeper) ConvertGasFeesToUsdc(ctx sdk.Context, baseCurrency string) sdk.Coins {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feesCollected := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	// Total Swapped coin
	totalSwappedCoins := sdk.Coins{}
	for _, tokenIn := range feesCollected {
		// if it is base currency - usdc, we don't need convert. We just need to collect it to fee wallet.
		if tokenIn.Denom == baseCurrency {
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

		// Sum total swapped
		totalSwappedCoins = totalSwappedCoins.Add(swappedCoins...)
	}

	return totalSwappedCoins
}

func (k Keeper) CollectGasFees(ctx sdk.Context, baseCurrency string) sdk.DecCoins {
	params := k.GetParams(ctx)

	// Calculate each portion of Gas fees collected - stakers, LPs
	fees := k.ConvertGasFeesToUsdc(ctx, baseCurrency)
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(fees...)

	gasFeesForLpsDec := gasFeeCollectedDec.MulDecTruncate(params.RewardPortionForLps)
	gasFeesForStakersDec := gasFeeCollectedDec.MulDecTruncate(params.RewardPortionForStakers)
	gasFeesForProtocolDec := gasFeeCollectedDec.Sub(gasFeesForLpsDec).Sub(gasFeesForStakersDec)

	k.AddFeeInfo(ctx, gasFeesForLpsDec.AmountOf(baseCurrency), gasFeesForStakersDec.AmountOf(baseCurrency), gasFeesForProtocolDec.AmountOf(baseCurrency), true)

	lpsGasFeeCoins, _ := gasFeesForLpsDec.TruncateDecimal()
	protocolGasFeeCoins, _ := gasFeesForProtocolDec.TruncateDecimal()

	// Send coins from fee collector name to masterchef
	if lpsGasFeeCoins.IsAllPositive() {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, types.ModuleName, lpsGasFeeCoins)
		if err != nil {
			panic(err)
		}
	}

	// Send coins to protocol revenue address
	if protocolGasFeeCoins.IsAllPositive() {
		protocolRevenueAddress, err := sdk.AccAddressFromBech32(params.ProtocolRevenueAddress)
		if err != nil {
			// Handle the error by skipping the fee distribution
			ctx.Logger().Error("Invalid protocol revenue address", "error", err)
			return gasFeesForLpsDec
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, protocolRevenueAddress, protocolGasFeeCoins)
		if err != nil {
			panic(err)
		}
	}
	return gasFeesForLpsDec
}

// Collect all DEX revenues to DEX revenue wallet,
// while tracking the 60% of it for LPs reward distribution
// transfer collected fees from different wallets(liquidity pool, perpetual module etc) to the distribution module account
// Assume this is already in USDC.
func (k Keeper) CollectDEXRevenue(ctx sdk.Context) (sdk.Coins, sdk.DecCoins, map[uint64]sdk.Dec) {
	// Total colllected revenue amount
	amountTotalCollected := sdk.Coins{}
	amountLPsCollected := sdk.DecCoins{}
	rewardsPerPool := make(map[uint64]sdk.Dec)

	// Iterate to calculate total Eden from LpElys, MElys committed
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		// Get pool Id
		poolId := p.GetPoolId()

		// Get dex rewards per pool
		revenueAddress := ammtypes.NewPoolRevenueAddress(poolId)

		// Transfer revenue to a single wallet of DEX revenue wallet.
		revenue := k.bankKeeper.GetAllBalances(ctx, revenueAddress)
		if revenue.IsAllPositive() {
			err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, revenueAddress, types.ModuleName, revenue)
			if err != nil {
				panic(err)
			}
		}

		// LPs Portion param
		params := k.GetParams(ctx)
		rewardPortionForLps := params.RewardPortionForLps
		rewardPortionForStakers := params.RewardPortionForStakers

		// Calculate revenue portion for LPs
		revenueDec := sdk.NewDecCoinsFromCoins(revenue...)

		// LPs portion of pool revenue
		revenuePortionForLPs := revenueDec.MulDecTruncate(rewardPortionForLps)
		revenuePortionForStakers := revenueDec.MulDecTruncate(rewardPortionForStakers)
		revenuePortionForProtocol := revenueDec.Sub(revenuePortionForLPs).Sub(revenuePortionForStakers)
		stakerRevenueCoins, _ := revenuePortionForStakers.TruncateDecimal()
		protocolRevenueCoins, _ := revenuePortionForProtocol.TruncateDecimal()

		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if found {
			k.AddFeeInfo(ctx, revenuePortionForLPs.AmountOf(baseCurrency), revenuePortionForStakers.AmountOf(baseCurrency), revenuePortionForProtocol.AmountOf(baseCurrency), false)
		}

		// Send coins to fee collector name
		if stakerRevenueCoins.IsAllPositive() {
			err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, stakerRevenueCoins)
			if err != nil {
				panic(err)
			}
		}

		// Send coins to protocol revenue address
		if protocolRevenueCoins.IsAllPositive() {
			protocolRevenueAddress := sdk.MustAccAddressFromBech32(params.ProtocolRevenueAddress)
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolRevenueAddress, protocolRevenueCoins)
			if err != nil {
				panic(err)
			}
		}

		// Store revenue portion for Lps temporarily
		if found {
			rewardsPerPool[poolId] = revenuePortionForLPs.AmountOf(baseCurrency)
		}

		// Sum total collected amount
		amountTotalCollected = amountTotalCollected.Add(revenue...)

		// Sum total amount for LPs
		amountLPsCollected = amountLPsCollected.Add(revenuePortionForLPs...)

		return false
	})

	return amountTotalCollected, amountLPsCollected, rewardsPerPool
}

// Calculate Proxy TVL
func (k Keeper) CalculateProxyTVL(ctx sdk.Context, baseCurrency string) (sdk.Dec, sdk.Dec) {
	// Ensure stablestakePoolParams exist
	stableStakePoolId := uint64(stabletypes.PoolId)
	_, found := k.GetPoolInfo(ctx, stableStakePoolId)
	if !found {
		k.InitStableStakePoolParams(ctx, stableStakePoolId)
	}

	multipliedShareSum := sdk.ZeroDec()
	multipliedShareSumOnlyEden := sdk.ZeroDec()
	for _, pool := range k.GetAllPoolInfos(ctx) {
		tvl := k.GetPoolTVL(ctx, pool.PoolId)
		proxyTVL := tvl.Mul(pool.Multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)

		/// Calculate total pool share by TVL and multiplier only when eden rewards is enable
		if pool.EnableEdenRewards {
			multipliedShareSumOnlyEden = multipliedShareSumOnlyEden.Add(proxyTVL)
		}
	}

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum, multipliedShareSumOnlyEden
}

// InitPoolParams: creates a poolInfo at the time of pool creation.
func (k Keeper) InitPoolParams(ctx sdk.Context, poolId uint64) bool {
	_, found := k.GetPoolInfo(ctx, poolId)
	if !found {
		poolInfo := types.PoolInfo{
			// reward amount
			PoolId: poolId,
			// reward wallet address
			RewardWallet: ammtypes.NewPoolRevenueAddress(poolId).String(),
			// multiplier for lp rewards
			Multiplier: sdk.NewDec(1),
			// Eden APR, updated at every distribution
			EdenApr: math.LegacyZeroDec(),
			// Dex APR, updated at every distribution
			DexApr: math.LegacyZeroDec(),
			// Gas APR, updated at every distribution
			GasApr: math.LegacyZeroDec(),
			// External Incentive APR, updated at every distribution
			ExternalIncentiveApr: math.LegacyZeroDec(),
			// external reward denoms on the pool
			ExternalRewardDenoms: []string{},
		}
		k.SetPoolInfo(ctx, poolInfo)
	}

	return true
}

// InitStableStakePoolMultiplier: create a stable stake pool information responding to the pool creation.
func (k Keeper) InitStableStakePoolParams(ctx sdk.Context, poolId uint64) bool {
	_, found := k.GetPoolInfo(ctx, poolId)
	if !found {
		poolInfo := types.PoolInfo{
			// reward amount
			PoolId: poolId,
			// reward wallet address
			RewardWallet: stabletypes.PoolAddress().String(),
			// multiplier for lp rewards
			Multiplier: sdk.NewDec(1),
			// Eden APR, updated at every distribution
			EdenApr: math.LegacyZeroDec(),
			// Dex APR, updated at every distribution
			DexApr: math.LegacyZeroDec(),
			// Gas APR, updated at every distribution
			GasApr: math.LegacyZeroDec(),
			// External Incentive APR, updated at every distribution
			ExternalIncentiveApr: math.LegacyZeroDec(),
			// external reward denoms on the pool
			ExternalRewardDenoms: []string{},
		}
		k.SetPoolInfo(ctx, poolInfo)
	}

	return true
}

// Update APR for AMM pool
func (k Keeper) UpdateAmmPoolAPR(ctx sdk.Context, totalBlocksPerYear int64, totalProxyTVL sdk.Dec, edenDenomPrice sdk.Dec) {
	baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)

	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool Id
		poolId := p.GetPoolId()

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, poolId)
		if !found {
			k.InitPoolParams(ctx, poolId)
			poolInfo, _ = k.GetPoolInfo(ctx, poolId)
		}

		if tvl.IsZero() {
			return false
		}

		firstAccum := k.FirstPoolRewardsAccum(ctx, poolId)
		lastAccum := k.LastPoolRewardsAccum(ctx, poolId)
		if lastAccum.Timestamp == 0 {
			return false
		}

		if firstAccum.Timestamp == lastAccum.Timestamp {
			poolInfo.DexApr = lastAccum.DexReward.
				MulInt64(totalBlocksPerYear).
				Mul(usdcDenomPrice).
				Quo(tvl)

			poolInfo.GasApr = lastAccum.GasReward.
				MulInt64(totalBlocksPerYear).
				Mul(usdcDenomPrice).
				Quo(tvl)
		} else {
			duration := lastAccum.Timestamp - firstAccum.Timestamp
			secondsInYear := int64(86400 * 360)

			poolInfo.DexApr = lastAccum.DexReward.Sub(firstAccum.DexReward).
				MulInt64(secondsInYear).
				QuoInt64(int64(duration)).
				Mul(usdcDenomPrice).
				Quo(tvl)

			poolInfo.GasApr = lastAccum.GasReward.Sub(firstAccum.GasReward).
				MulInt64(secondsInYear).
				QuoInt64(int64(duration)).
				Mul(usdcDenomPrice).
				Quo(tvl)
		}
		k.SetPoolInfo(ctx, poolInfo)
		return false
	})
}
