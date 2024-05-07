package keeper

import (
	"errors"

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
	canDistribute := k.CanDistributeLPRewards(ctx)
	if !canDistribute {
		return
	}

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

			tvl := k.GetPoolTVL(ctx, pool.PoolId)
			if tvl.IsPositive() {
				yearlyIncentiveRewardsTotal := externalIncentive.AmountPerBlock.
					Mul(lpIncentive.TotalBlocksPerYear).
					Quo(pool.NumBlocks)

				baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
				if found {
					pool.ExternalIncentiveApr = sdk.NewDecFromInt(yearlyIncentiveRewardsTotal).
						Mul(k.amm.GetTokenPrice(ctx, externalIncentive.RewardDenom, baseCurrency)).
						Quo(tvl)
					k.SetPool(ctx, pool)
				}
			}
		}

		if curBlockHeight.Uint64() == externalIncentive.ToBlock {
			k.RemoveExternalIncentive(ctx, externalIncentive.Id)
		}
	}
}

func (k Keeper) ProcessLPRewardDistribution(ctx sdk.Context) {
	// Read tokenomics time based inflation params and update incentive module params.
	if !k.ProcessUpdateIncentiveParams(ctx) {
		ctx.Logger().Error("Invalid tokenomics params", "error", errors.New("invalid tokenomics params"))
		return
	}

	canDistribute := k.CanDistributeLPRewards(ctx)
	if canDistribute {
		err := k.UpdateLPRewards(ctx)
		if err != nil {
			ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
		}
	}
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) bool {
	// Non-linear inflation per year happens and this includes yearly inflation data
	listTimeBasedInflations := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflations) < 1 {
		return false
	}

	params := k.GetParams(ctx)

	for _, inflation := range listTimeBasedInflations {
		// Finding only current inflation data - and skip rest
		if inflation.StartBlockHeight > uint64(ctx.BlockHeight()) || inflation.EndBlockHeight < uint64(ctx.BlockHeight()) {
			continue
		}

		totalBlocksPerYear := sdk.NewInt(int64(inflation.EndBlockHeight - inflation.StartBlockHeight + 1))

		// If totalBlocksPerYear is zero, we skip this inflation to avoid division by zero
		if totalBlocksPerYear == sdk.ZeroInt() {
			continue
		}
		blocksDistributed := sdk.NewInt(ctx.BlockHeight() - int64(inflation.StartBlockHeight))

		incentiveInfo := types.IncentiveInfo{
			// reward amount in eden for 1 year
			EdenAmountPerYear: sdk.NewInt(int64(inflation.Inflation.LmRewards)),
			// starting block height of the distribution
			DistributionStartBlock: sdk.NewInt(int64(inflation.StartBlockHeight)),
			// distribution duration - block number per year
			TotalBlocksPerYear: totalBlocksPerYear,
			// number of blocks distributed
			BlocksDistributed: blocksDistributed,
		}

		if params.LpIncentives == nil {
			params.LpIncentives = &incentiveInfo
		} else {
			// If any of block number related parameter changed, we re-calculate the current epoch
			if params.LpIncentives.DistributionStartBlock != incentiveInfo.DistributionStartBlock ||
				params.LpIncentives.TotalBlocksPerYear != incentiveInfo.TotalBlocksPerYear {
				params.LpIncentives.BlocksDistributed = blocksDistributed
			}
			params.LpIncentives.EdenAmountPerYear = incentiveInfo.EdenAmountPerYear
			params.LpIncentives.DistributionStartBlock = incentiveInfo.DistributionStartBlock
			params.LpIncentives.TotalBlocksPerYear = incentiveInfo.TotalBlocksPerYear
		}
		break
	}

	k.SetParams(ctx, params)
	return true
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
	lpIncentive.BlocksDistributed = lpIncentive.BlocksDistributed.Add(sdk.OneInt())
	if lpIncentive.BlocksDistributed.GTE(lpIncentive.TotalBlocksPerYear) || curBlockHeight.GT(lpIncentive.TotalBlocksPerYear.Add(lpIncentive.DistributionStartBlock)) {
		params.LpIncentives = nil
		k.SetParams(ctx, params)
		return false
	}

	params.LpIncentives.BlocksDistributed = lpIncentive.BlocksDistributed
	k.SetParams(ctx, params)

	return true
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
	_, dexRevenueForLps := k.CollectDEXRevenue(ctx)

	// USDC amount in sdk.Dec type
	dexUsdcAmountForLps := dexRevenueForLps.AmountOf(baseCurrency)
	gasFeeUsdcAmountForLps := gasFeesForLpsDec.AmountOf(baseCurrency)

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

	// Ensure lpIncentive.TotalBlocksPerYear is not zero to avoid division by zero
	if lpIncentive.TotalBlocksPerYear.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate eden amount per block
	lpsEdenAmount := lpIncentive.EdenAmountPerYear.
		Quo(lpIncentive.TotalBlocksPerYear)

	// Maximum eden APR - 30% by default
	// Allocated for staking per day = (0.3/365)* (total weighted proxy TVL)
	edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

	// Ensure edenDenomPrice is not zero to avoid division by zero
	if edenDenomPrice.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid eden price")
	}

	// Distribute Eden / USDC Rewards
	for _, pool := range k.GetAllPools(ctx) {
		var err error
		tvl := k.GetPoolTVL(ctx, pool.PoolId)
		proxyTVL := tvl.Mul(pool.Multiplier)
		if proxyTVL.IsZero() {
			continue
		}

		poolShare := sdk.ZeroDec()
		if totalProxyTVL.IsPositive() {
			poolShare = proxyTVL.Quo(totalProxyTVL)
		}

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShare.MulInt(lpsEdenAmount)

		poolMaxEdenAmount := params.MaxEdenRewardAprLps.
			Mul(tvl).
			QuoInt(lpIncentive.TotalBlocksPerYear).
			Quo(edenDenomPrice)

		// Use min amount (eden allocation from tokenomics and max apr based eden amount)
		newEdenAllocatedForPool = sdk.MinDec(newEdenAllocatedForPool, poolMaxEdenAmount)
		err = k.cmk.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ptypes.Eden, newEdenAllocatedForPool.TruncateInt())})
		if err != nil {
			panic(err)
		}

		// Get gas fee rewards per pool
		gasRewardsAllocatedForPool := poolShare.Mul(gasFeeUsdcAmountForLps)

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

	// Set DexRewards info
	params.DexRewardsLps.NumBlocks = sdk.OneInt()
	params.DexRewardsLps.Amount = dexUsdcAmountForLps.Add(gasFeeUsdcAmountForLps)
	k.SetParams(ctx, params)

	// Update APR for amm pools
	k.UpdateAmmPoolAPR(ctx, lpIncentive.TotalBlocksPerYear, totalProxyTVL, edenDenomPrice)

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
		protocolRevenueAddress := sdk.MustAccAddressFromBech32(params.ProtocolRevenueAddress)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, protocolRevenueAddress, protocolGasFeeCoins)
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
// TODO:
// + Collect revenue from perpetual, lend module
func (k Keeper) CollectDEXRevenue(ctx sdk.Context) (sdk.Coins, sdk.DecCoins) {
	// Total colllected revenue amount
	amountTotalCollected := sdk.Coins{}
	amountLPsCollected := sdk.DecCoins{}

	k.tci.PoolRevenueTrack = make(map[string]sdk.Dec)

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
	_, found := k.GetPool(ctx, stableStakePoolId)
	// Ensure stablestakePoolParams exist
	if !found {
		k.InitStableStakePoolParams(ctx, stableStakePoolId)
	}
	for _, pool := range k.GetAllPools(ctx) {
		tvl := k.GetPoolTVL(ctx, pool.PoolId)
		proxyTVL := tvl.Mul(pool.Multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)
	}

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum
}

// InitPoolParams: creates a poolInfo at the time of pool creation.
func (k Keeper) InitPoolParams(ctx sdk.Context, poolId uint64) bool {
	_, found := k.GetPool(ctx, poolId)
	if !found {
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
		k.SetPool(ctx, poolInfo)
	}

	return true
}

// InitStableStakePoolMultiplier: create a stable stake pool information responding to the pool creation.
func (k Keeper) InitStableStakePoolParams(ctx sdk.Context, poolId uint64) bool {
	_, found := k.GetPool(ctx, poolId)
	if !found {
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
		k.SetPool(ctx, poolInfo)
	}

	return true
}

// Calculate pool share for stable stake pool
func (k Keeper) CalculatePoolShareForStableStakeLPs(ctx sdk.Context, totalProxyTVL sdk.Dec, baseCurrency string) sdk.Dec {
	// ------------ New Eden calculation -------------------
	// -----------------------------------------------------
	// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
	// Pool share = 80
	// edenAmountLp = 100
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
