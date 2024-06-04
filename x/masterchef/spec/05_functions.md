<!--
order: 5
-->

# Functions

## EndBlocker

The `EndBlocker` function is called at the end of each block to perform necessary updates and maintenance for the `masterchef` module. It processes LP rewards and external incentives distribution.

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
    k.ProcessLPRewardDistribution(ctx)
    k.ProcessExternalRewardsDistribution(ctx)
}
```

### ProcessLPRewardDistribution

The `ProcessLPRewardDistribution` function distributes rewards to liquidity providers. It updates the incentive parameters and calculates the rewards based on the collected fees and staking conditions.

```go
func (k Keeper) ProcessLPRewardDistribution(ctx sdk.Context) {
    k.ProcessUpdateIncentiveParams(ctx)
    err := k.UpdateLPRewards(ctx)
    if err != nil {
        ctx.Logger().Error("Failed to update LP rewards", "error", err)
    }
}
```

### ProcessExternalRewardsDistribution

The `ProcessExternalRewardsDistribution` function distributes external incentives to the specified pools within the defined block range.

```go
func (k Keeper) ProcessExternalRewardsDistribution(ctx sdk.Context) {
    baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
    curBlockHeight := sdk.NewInt(ctx.BlockHeight())
    totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear

    externalIncentives := k.GetAllExternalIncentives(ctx)
    externalIncentiveAprs := make(map[uint64]math.LegacyDec)
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
                    Mul(sdk.NewInt(totalBlocksPerYear)).
                    Quo(pool.NumBlocks)

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
                k.SetPool(ctx, pool)
            }
        }

        if curBlockHeight.Uint64() == externalIncentive.ToBlock {
            k.RemoveExternalIncentive(ctx, externalIncentive.Id)
        }
    }
}
```

### UpdateLPRewards

The `UpdateLPRewards` function updates the rewards for liquidity providers by calculating the total rewards based on the collected fees and staking conditions.

```go
func (k Keeper) UpdateLPRewards(ctx sdk.Context) error {
    params := k.GetParams(ctx)
    lpIncentive := params.LpIncentives

    baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
    if !found {
        return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
    }

    gasFeesForLpsDec := k.CollectGasFees(ctx, baseCurrency)
    _, dexRevenueForLps, rewardsPerPool := k.CollectDEXRevenue(ctx)

    dexUsdcAmountForLps := dexRevenueForLps.AmountOf(baseCurrency)
    gasFeeUsdcAmountForLps := gasFeesForLpsDec.AmountOf(baseCurrency)

    totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)
    totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear
    if totalBlocksPerYear == 0 {
        return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
    }

    edenAmountPerYear := sdk.ZeroInt()
    if lpIncentive != nil && lpIncentive.EdenAmountPerYear.IsPositive() {
        edenAmountPerYear = lpIncentive.EdenAmountPerYear
    }
    lpsEdenAmount = edenAmountPerYear.Quo(sdk.NewInt(totalBlocksPerYear))

    edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)
    if edenDenomPrice.IsZero() {
        return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid eden price")
    }

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

        newEdenAllocatedForPool := poolShare.MulInt(lpsEdenAmount)

        poolMaxEdenAmount := params.MaxEdenRewardAprLps.
            Mul(tvl).
            QuoInt64(totalBlocksPerYear).
            Quo(edenDenomPrice)

        newEdenAllocatedForPool = sdk.MinDec(newEdenAllocatedForPool, poolMaxEdenAmount)
        if newEdenAllocatedForPool.IsPositive() {
            err = k.cmk.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ptypes.Eden, newEdenAllocatedForPool.TruncateInt())})
            if err != nil {
                return err
            }
        }

        gasRewardsAllocatedForPool := poolShare.Mul(gasFeeUsdcAmountForLps)

        dexRewardsAllocatedForPool, ok := rewardsPerPool[pool.PoolId]
        if !ok {
            dexRewardsAllocatedForPool = sdk.NewDec(0)
        }

        k.UpdateAccPerShare(ctx, pool.PoolId, ptypes.Eden, newEdenAllocatedForPool.TruncateInt())
        k.UpdateAccPerShare(ctx, pool.PoolId, k.GetBaseCurrencyDenom(ctx), gasRewardsAllocatedForPool.Add(dexRewardsAllocatedForPool).TruncateInt())

        pool.EdenRewardAmountGiven = newEdenAllocatedForPool.RoundInt()
        pool.DexRewardAmountGiven = gasRewardsAllocatedForPool.Add(dexRewardsAllocatedForPool)
        k.SetPool(ctx,

 pool)
    }

    params.DexRewardsLps.NumBlocks = sdk.OneInt()
    params.DexRewardsLps.Amount = dexUsdcAmountForLps.Add(gasFeeUsdcAmountForLps)
    k.SetParams(ctx, params)

    k.UpdateAmmPoolAPR(ctx, totalBlocksPerYear, totalProxyTVL, edenDenomPrice)

    return nil
}
```

### UpdateAccPerShare

The `UpdateAccPerShare` function updates the accumulated reward per share for a specific pool.

```go
func (k Keeper) UpdateAccPerShare(ctx sdk.Context, poolId uint64, rewardDenom string, amount sdk.Int) {
    poolRewardInfo, found := k.GetPoolRewardInfo(ctx, poolId, rewardDenom)
    if !found {
        poolRewardInfo = types.PoolRewardInfo{
            PoolId:                poolId,
            RewardDenom:           rewardDenom,
            PoolAccRewardPerShare: sdk.NewDec(0),
            LastUpdatedBlock:      0,
        }
    }

    totalCommit := k.GetPoolTotalCommit(ctx, poolId)
    if totalCommit.IsZero() {
        return
    }
    poolRewardInfo.PoolAccRewardPerShare = poolRewardInfo.PoolAccRewardPerShare.Add(
        math.LegacyNewDecFromInt(amount.Mul(ammtypes.OneShare)).
            Quo(math.LegacyNewDecFromInt(totalCommit)),
    )
    poolRewardInfo.LastUpdatedBlock = uint64(ctx.BlockHeight())
    k.SetPoolRewardInfo(ctx, poolRewardInfo)
}
```

### CollectGasFees

The `CollectGasFees` function collects gas fees and allocates them to LPs and the protocol.

```go
func (k Keeper) CollectGasFees(ctx sdk.Context, baseCurrency string) sdk.DecCoins {
    params := k.GetParams(ctx)

    fees := k.ConvertGasFeesToUsdc(ctx, baseCurrency)
    gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(fees...)

    gasFeesForLpsDec := gasFeeCollectedDec.MulDecTruncate(params.RewardPortionForLps)
    gasFeesForStakersDec := gasFeeCollectedDec.MulDecTruncate(params.RewardPortionForStakers)
    gasFeesForProtocolDec := gasFeeCollectedDec.Sub(gasFeesForLpsDec).Sub(gasFeesForStakersDec)

    lpsGasFeeCoins, _ := gasFeesForLpsDec.TruncateDecimal()
    protocolGasFeeCoins, _ := gasFeesForProtocolDec.TruncateDecimal()

    if lpsGasFeeCoins.IsAllPositive() {
        err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, types.ModuleName, lpsGasFeeCoins)
        if err != nil {
            panic(err)
        }
    }

    if protocolGasFeeCoins.IsAllPositive() {
        protocolRevenueAddress := sdk.MustAccAddressFromBech32(params.ProtocolRevenueAddress)
        err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, protocolRevenueAddress, protocolGasFeeCoins)
        if err != nil {
            panic(err)
        }
    }
    return gasFeesForLpsDec
}
```

### CollectDEXRevenue

The `CollectDEXRevenue` function collects DEX revenue and distributes it to LPs and the protocol.

```go
func (k Keeper) CollectDEXRevenue(ctx sdk.Context) (sdk.Coins, sdk.DecCoins, map[uint64]sdk.Dec) {
    amountTotalCollected := sdk.Coins{}
    amountLPsCollected := sdk.DecCoins{}
    rewardsPerPool := make(map[uint64]sdk.Dec)

    k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
        poolId := p.GetPoolId()
        revenueAddress := ammtypes.NewPoolRevenueAddress(poolId)

        revenue := k.bankKeeper.GetAllBalances(ctx, revenueAddress)
        if revenue.IsAllPositive() {
            err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, revenueAddress, types.ModuleName, revenue)
            if err != nil {
                panic(err)
            }
        }

        params := k.GetParams(ctx)
        rewardPortionForLps := params.RewardPortionForLps
        rewardPortionForStakers := params.RewardPortionForStakers

        revenueDec := sdk.NewDecCoinsFromCoins(revenue...)

        revenuePortionForLPs := revenueDec.MulDecTruncate(rewardPortionForLps)
        revenuePortionForStakers := revenueDec.MulDecTruncate(rewardPortionForStakers)
        revenuePortionForProtocol := revenueDec.Sub(revenuePortionForLPs).Sub(revenuePortionForStakers)
        stakerRevenueCoins, _ := revenuePortionForStakers.TruncateDecimal()
        protocolRevenueCoins, _ := revenuePortionForProtocol.TruncateDecimal()

        if stakerRevenueCoins.IsAllPositive() {
            err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, stakerRevenueCoins)
            if err != nil {
                panic(err)
            }
        }

        if protocolRevenueCoins.IsAllPositive() {
            protocolRevenueAddress := sdk.MustAccAddressFromBech32(params.ProtocolRevenueAddress)
            err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolRevenueAddress, protocolRevenueCoins)
            if err != nil {
                panic(err)
            }
        }

        rewardsPerPool[poolId] = revenuePortionForLPs.AmountOf(ptypes.BaseCurrency)

        amountTotalCollected = amountTotalCollected.Add(revenue...)
        amountLPsCollected = amountLPsCollected.Add(revenuePortionForLPs...)

        return false
    })

    return amountTotalCollected, amountLPsCollected, rewardsPerPool
}
```

### CalculateProxyTVL

The `CalculateProxyTVL` function calculates the proxy total value locked (TVL) for the pools.

```go
func (k Keeper) CalculateProxyTVL(ctx sdk.Context, baseCurrency string) sdk.Dec {
    stableStakePoolId := uint64(stabletypes.PoolId)
    _, found := k.GetPool(ctx, stableStakePoolId)
    if !found {
        k.InitStableStakePoolParams(ctx, stableStakePoolId)
    }

    multipliedShareSum := sdk.ZeroDec()
    for _, pool := range k.GetAllPools(ctx) {
        tvl := k.GetPoolTVL(ctx, pool.PoolId)
        proxyTVL := tvl.Mul(pool.Multiplier)

        multipliedShareSum = multipliedShareSum.Add(proxyTVL)
    }

    return multipliedShareSum
}
```

### UpdateAmmPoolAPR

The `UpdateAmmPoolAPR` function updates the APR for AMM pools.

```go
func (k Keeper) UpdateAmmPoolAPR(ctx sdk.Context, totalBlocksPerYear int64, totalProxyTVL sdk.Dec, edenDenomPrice sdk.Dec) {
    baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
    usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)

    k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
        tvl, err := p.TVL(ctx, k.oracleKeeper)
        if err != nil {
            return false
        }

        poolId := p.GetPoolId()
        poolInfo, found := k.GetPool(ctx, poolId)
        if !found {
            k.InitPoolParams(ctx, poolId)
            poolInfo, _ = k.GetPool(ctx, poolId)
        }

        poolInfo.NumBlocks = sdk.OneInt()

        if tvl.IsZero() {
            return false
        }

        if poolInfo.NumBlocks.IsZero() {
            return false
        }

        yearlyDexRewardsTotal := poolInfo.DexRewardAmountGiven.
            MulInt64(totalBlocksPerYear).
            QuoInt(poolInfo.NumBlocks)
        poolInfo.DexApr = yearlyDexRewardsTotal.Mul(usdcDenomPrice).Quo(tvl)

        yearlyEdenRewardsTotal := poolInfo.EdenRewardAmountGiven.
            Mul(sdk.NewInt(totalBlocksPerYear)).
            Quo(poolInfo.NumBlocks)

        poolInfo.EdenApr = sdk.NewDecFromInt(yearlyEdenRewardsTotal).
            Mul(edenDenomPrice).
            Quo(tvl)

        k.SetPool(ctx, poolInfo)
        return false
    })
}
```
