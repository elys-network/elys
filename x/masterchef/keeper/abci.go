package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

// EndBlocker of amm module
func (k Keeper) EndBlocker(ctx sdk.Context) error {

	k.DeleteFeeInfo(ctx)

	// distribute LP rewards
	err := k.ProcessLPRewardDistribution(ctx)
	if err != nil {
		return err
	}
	// distribute external rewards
	k.ProcessExternalRewardsDistribution(ctx)

	// convert balances in taker address to elys and burn them
	k.ProcessTakerFee(ctx)
	return nil
}

func (k Keeper) GetPoolTVL(ctx sdk.Context, poolId uint64) math.LegacyDec {
	if poolId >= stabletypes.UsdcPoolId {
		return k.stableKeeper.TVL(ctx, k.oracleKeeper, poolId)
	}
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if found {
		tvl, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
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
	totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)

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
					Mul(math.NewInt(totalBlocksPerYear))

				apr := yearlyIncentiveRewardsTotal.ToLegacyDec().
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

func (k Keeper) ProcessLPRewardDistribution(ctx sdk.Context) error {
	// Read tokenomics time based inflation params and update incentive module params.
	k.ProcessUpdateIncentiveParams(ctx)

	err := k.UpdateLPRewards(ctx)
	if err != nil {
		ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
		return err
	}
	return nil
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) {
	// Non-linear inflation per year happens and this includes yearly inflation data
	listTimeBasedInflation := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflation) == 0 {
		return
	}

	params := k.GetParams(ctx)
	for _, inflation := range listTimeBasedInflation {
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
			EdenAmountPerYear: math.NewInt(int64(inflation.Inflation.LmRewards)),
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
	gasFeesForLpsDec, err := k.CollectGasFees(ctx, baseCurrency)
	if err != nil {
		return err
	}
	_, _, rewardsPerPool, err := k.CollectDEXRevenue(ctx)
	if err != nil {
		return err
	}

	// USDC amount in math.LegacyDec type
	gasFeeUsdcAmountForLps := gasFeesForLpsDec.AmountOf(baseCurrency)

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have multiplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL, totalProxyTvlEdenEnable := k.CalculateProxyTVL(ctx, baseCurrency)

	// Ensure totalBlocksPerYear is not zero to avoid division by zero
	totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
	if totalBlocksPerYear == 0 {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate eden amount per block
	edenAmountPerYear := math.ZeroInt()
	if lpIncentive != nil && lpIncentive.EdenAmountPerYear.IsPositive() {
		edenAmountPerYear = lpIncentive.EdenAmountPerYear
	}
	lpsEdenAmount := edenAmountPerYear.Quo(math.NewInt(totalBlocksPerYear))

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

		poolShare := math.LegacyZeroDec()
		poolShareEdenEnable := math.LegacyZeroDec()
		if totalProxyTVL.IsPositive() {
			poolShare = proxyTVL.Quo(totalProxyTVL)
		}

		if totalProxyTvlEdenEnable.IsPositive() {
			poolShareEdenEnable = proxyTVL.Quo(totalProxyTvlEdenEnable)
		}

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := math.LegacyZeroDec()

		// Maximum eden APR - 30% by default
		poolMaxEdenAmount := params.MaxEdenRewardAprLps.
			Mul(proxyTVL).
			QuoInt64(totalBlocksPerYear).
			Quo(edenDenomPrice)

		// Use min amount (eden allocation from tokenomics and max apr based eden amount)
		if pool.EnableEdenRewards {
			newEdenAllocatedForPool = poolShareEdenEnable.MulInt(lpsEdenAmount)
			newEdenAllocatedForPool = math.LegacyMinDec(newEdenAllocatedForPool, poolMaxEdenAmount)
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
			dexRewardsAllocatedForPool = math.LegacyNewDec(0)
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
			if firstAccum.Timestamp == 0 || int64(firstAccum.Timestamp+dataLifetime) >= ctx.BlockTime().Unix() {
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
			pool.EdenApr = math.LegacyZeroDec()
		}

		k.SetPoolInfo(ctx, pool)
	}

	// Update APR for amm pools
	k.UpdateAmmStablePoolAPR(ctx, totalBlocksPerYear, totalProxyTVL, edenDenomPrice)

	return nil
}

// Move gas fees collected to dex revenue wallet
// Convert it into USDC
func (k Keeper) ConvertGasFeesToUsdc(ctx sdk.Context, baseCurrency string, address sdk.AccAddress) (sdk.Coins, error) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feesCollected := k.bankKeeper.GetAllBalances(ctx, address)

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
		pool, found := k.amm.GetBestPoolWithDenoms(ctx, []string{tokenIn.Denom, baseCurrency}, false)
		if !found {
			// If there is a denom for which pool doesn't exist, log it, otherwise
			// if pool exist, throw error later
			ctx.Logger().Info("Pool not found for denom: " + tokenIn.Denom)
			continue
		}

		tokenOutAmount, err := k.amm.InternalSwapExactAmountIn(ctx, address, address, pool, tokenIn, baseCurrency, math.ZeroInt(), math.LegacyZeroDec())
		if err != nil {
			// Continue as we can swap it when this amount is higher
			if err == ammtypes.ErrTokenOutAmountZero {
				ctx.Logger().Info("Token out amount is zero(skipping conversion) for denom: " + tokenIn.Denom)
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						types.TypeEvtSkipSwap,
						sdk.NewAttribute("Token denom", tokenIn.Denom),
						sdk.NewAttribute("Token amount", "0"),
					),
				})
				continue
			}
			return sdk.Coins{}, err
		}

		// Swapped USDC coin
		swappedCoins := sdk.NewCoins(sdk.NewCoin(baseCurrency, tokenOutAmount))

		// Sum total swapped
		totalSwappedCoins = totalSwappedCoins.Add(swappedCoins...)
	}
	if !totalSwappedCoins.IsZero() {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtUsdcFee,
				sdk.NewAttribute("amount", totalSwappedCoins.String()),
			),
		})
	}

	return totalSwappedCoins, nil
}

func (k Keeper) CollectGasFees(ctx sdk.Context, baseCurrency string) (sdk.DecCoins, error) {
	params := k.GetParams(ctx)
	estakingParams := k.estakingKeeper.GetParams(ctx)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	// Calculate each portion of Gas fees collected - stakers, LPs
	fees, err := k.ConvertGasFeesToUsdc(ctx, baseCurrency, feeCollector.GetAddress())
	if err != nil {
		return sdk.DecCoins{}, err
	}
	if fees.IsZero() {
		return sdk.DecCoins{}, nil
	}
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(fees...)

	gasFeesForLpsDec := gasFeeCollectedDec.MulDecTruncate(params.RewardPortionForLps)
	gasFeesForStakersDec := gasFeeCollectedDec.MulDecTruncate(params.RewardPortionForStakers)
	gasFeesForProtocolDec := gasFeeCollectedDec.Sub(gasFeesForLpsDec).Sub(gasFeesForStakersDec)

	k.AddFeeInfo(ctx, gasFeesForLpsDec.AmountOf(baseCurrency), gasFeesForStakersDec.AmountOf(baseCurrency), gasFeesForProtocolDec.AmountOf(baseCurrency), true)

	lpsGasFeeCoins, _ := gasFeesForLpsDec.TruncateDecimal()
	stakersGasFeeCoins, _ := gasFeesForStakersDec.TruncateDecimal() // Before ccv, this used to be remain in FeeCollectorName
	protocolGasFeeCoins, _ := gasFeesForProtocolDec.TruncateDecimal()

	if stakersGasFeeCoins.IsAllPositive() {
		// Earlier this amount remained in FeeCollectorName and distribution module handled it using FeeCollectorName.
		// But after ccv, distribution module only acts on ccvconsumertypes.ConsumerRedistributeName
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, ccvconsumertypes.ConsumerRedistributeName, stakersGasFeeCoins)
		if err != nil {
			return sdk.DecCoins{}, err
		}
	}

	// Send coins from fee collector name to masterchef
	if lpsGasFeeCoins.IsAllPositive() {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, types.ModuleName, lpsGasFeeCoins)
		if err != nil {
			return sdk.DecCoins{}, err
		}
	}

	// Send coins to protocol revenue address
	if protocolGasFeeCoins.IsAllPositive() {
		protocolRevenueAddress, err := sdk.AccAddressFromBech32(params.ProtocolRevenueAddress)
		if err != nil {
			// Handle the error by skipping the fee distribution
			ctx.Logger().Error("Invalid protocol revenue address", "error", err)
			return gasFeesForLpsDec, err
		}
		providerPortion := ammkeeper.PortionCoins(protocolGasFeeCoins, estakingParams.ProviderStakingRewardsPortion)
		consumerPortion := protocolGasFeeCoins.Sub(providerPortion...)

		// This will be sent to provider
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, ccvconsumertypes.ConsumerToSendToProviderName, providerPortion)
		if err != nil {
			return sdk.DecCoins{}, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, protocolRevenueAddress, consumerPortion)
		if err != nil {
			return sdk.DecCoins{}, err
		}
	}
	return gasFeesForLpsDec, nil
}

// Collect all DEX revenues to DEX revenue wallet,
// while tracking the 60% of it for LPs reward distribution
// transfer collected fees from different wallets(liquidity pool, perpetual module etc) to the distribution module account
// Assume this is already in USDC.
func (k Keeper) CollectDEXRevenue(ctx sdk.Context) (sdk.Coins, sdk.DecCoins, map[uint64]math.LegacyDec, error) {
	// Total collected revenue amount
	amountTotalCollected := sdk.Coins{}
	amountLPsCollected := sdk.DecCoins{}
	rewardsPerPool := make(map[uint64]math.LegacyDec)
	// LPs Portion param
	params := k.GetParams(ctx)
	estakingParams := k.estakingKeeper.GetParams(ctx)
	protocolRevenueAddress, err := sdk.AccAddressFromBech32(params.ProtocolRevenueAddress)
	if err != nil {
		return nil, nil, nil, err
	}

	// Iterate to calculate total Eden from LpElys, MElys committed
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		// Get pool Id
		poolId := p.GetPoolId()

		// Get dex rewards per pool
		revenueAddress := ammtypes.NewPoolRevenueAddress(poolId)

		// Transfer revenue to a single wallet of DEX revenue wallet.
		revenue := k.bankKeeper.GetAllBalances(ctx, revenueAddress)
		if revenue.IsAllPositive() {
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, revenueAddress, types.ModuleName, revenue)
			if err != nil {
				return true
			}
		}

		// Calculate revenue portion for LPs
		revenueDec := sdk.NewDecCoinsFromCoins(revenue...)

		// LPs portion of pool revenue
		revenuePortionForLPs := revenueDec.MulDecTruncate(params.RewardPortionForLps)
		revenuePortionForStakers := revenueDec.MulDecTruncate(params.RewardPortionForStakers)
		revenuePortionForProtocol := revenueDec.Sub(revenuePortionForLPs).Sub(revenuePortionForStakers)
		stakerRevenueCoins, _ := revenuePortionForStakers.TruncateDecimal()
		protocolRevenueCoins, _ := revenuePortionForProtocol.TruncateDecimal()

		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if found {
			k.AddFeeInfo(ctx, revenuePortionForLPs.AmountOf(baseCurrency), revenuePortionForStakers.AmountOf(baseCurrency), revenuePortionForProtocol.AmountOf(baseCurrency), false)
		}

		// Send coins to fee collector name
		if stakerRevenueCoins.IsAllPositive() {
			// The distribution module picks from ccvconsumertypes.ConsumerRedistributeName
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, ccvconsumertypes.ConsumerRedistributeName, stakerRevenueCoins)
			if err != nil {
				return true
			}
		}

		// Send coins to protocol revenue address
		if protocolRevenueCoins.IsAllPositive() {
			providerPortion := ammkeeper.PortionCoins(protocolRevenueCoins, estakingParams.ProviderStakingRewardsPortion)
			consumerPortion := stakerRevenueCoins.Sub(providerPortion...)

			// This will be sent to provider
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, ccvconsumertypes.ConsumerToSendToProviderName, providerPortion)
			if err != nil {
				return true
			}

			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolRevenueAddress, consumerPortion)
			if err != nil {
				return true
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
	if err != nil {
		return nil, nil, nil, err
	}

	return amountTotalCollected, amountLPsCollected, rewardsPerPool, nil
}

// Calculate Proxy TVL
func (k Keeper) CalculateProxyTVL(ctx sdk.Context, baseCurrency string) (math.LegacyDec, math.LegacyDec) {
	// Ensure stablestakePoolParams exist
	stableStakePoolId := uint64(stabletypes.UsdcPoolId)
	_, found := k.GetPoolInfo(ctx, stableStakePoolId)
	if !found {
		k.InitStableStakePoolParams(ctx, stableStakePoolId)
	}

	multipliedShareSum := math.LegacyZeroDec()
	multipliedShareSumOnlyEden := math.LegacyZeroDec()
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
			Multiplier: math.LegacyNewDec(1),
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
			// enable eden reward on the pool
			EnableEdenRewards: false,
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
			Multiplier: math.LegacyNewDec(1),
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

// Update APR for AMM pools and stable stake pools
func (k Keeper) UpdateAmmStablePoolAPR(ctx sdk.Context, totalBlocksPerYear int64, totalProxyTVL math.LegacyDec, edenDenomPrice math.LegacyDec) {
	baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)

	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
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

	k.stableKeeper.IterateLiquidityPools(ctx, func(p stabletypes.Pool) bool {
		tvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, p.Id)

		// Get pool Id
		poolId := p.GetId()

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

func (k Keeper) ProcessTakerFee(ctx sdk.Context) {
	collectionAddressString := k.parameterKeeper.GetParams(ctx).TakerFeeCollectionAddress
	// Convert balances in taker address to elys
	collectionAddress, err := sdk.AccAddressFromBech32(collectionAddressString)
	if err != nil {
		ctx.Logger().Error("Invalid taker fee collection address", "error", err)
		return
	}

	balances := k.bankKeeper.GetAllBalances(ctx, collectionAddress)
	for _, balance := range balances {
		if balance.Denom == ptypes.Elys {
			continue
		}
		_, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
			Sender:    collectionAddressString,
			Recipient: collectionAddressString,
			Amount:    sdk.NewCoin(balance.Denom, balance.Amount),
			DenomIn:   balance.Denom,
			DenomOut:  ptypes.Elys,
			MinAmount: sdk.NewCoin(ptypes.Elys, sdkmath.ZeroInt()),
		})
		if err != nil {
			ctx.Logger().Error("Failed to swap taker fee", "error", err)
			continue
		}
	}

	elysBalance := k.bankKeeper.GetBalance(ctx, collectionAddress, ptypes.Elys)
	if elysBalance.IsPositive() {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, collectionAddress, types.ModuleName, sdk.NewCoins(elysBalance))
		if err != nil {
			ctx.Logger().Error("Failed to send taker fee to masterchef", "error", err)
		} else {
			// burn elys token
			err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(elysBalance))
			if err != nil {
				ctx.Logger().Error("Failed to burn taker fee", "error", err)
			}

			// event for burning taker fees
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.TypeEvtTakerFeeBurn,
					sdk.NewAttribute("amount", elysBalance.String()),
				),
			})
		}
	}
}
