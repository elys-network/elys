package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

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
func (k Keeper) CollectDEXRevenue(ctx sdk.Context) (sdk.Coins, sdk.DecCoins, sdk.DecCoins) {
	// Total colllected revenue amount
	amountTotalCollected := sdk.Coins{}
	amountLPsCollected := sdk.DecCoins{}
	amountStakersCollected := sdk.DecCoins{}

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
		rewardPortionForLps := k.GetDEXRewardPortionForLPs(ctx)
		// Stakers Portion param
		rewardPortionForStakers := k.GetDEXRewardPortionForStakers(ctx)

		// Calculate revenue portion for LPs
		revenueDec := sdk.NewDecCoinsFromCoins(revenue...)

		// LPs portion of pool revenue
		revenuePortionForLPs := revenueDec.MulDecTruncate(rewardPortionForLps)
		revenuePortionForStakers := revenueDec.MulDecTruncate(rewardPortionForStakers)

		// Get track key
		trackKey := types.GetPoolRevenueTrackKey(poolId)

		// Store revenue portion for Lps temporarilly
		k.tci.PoolRevenueTrack[trackKey] = revenuePortionForLPs.AmountOf(ptypes.BaseCurrency)

		// Sum total collected amount
		amountTotalCollected = amountTotalCollected.Add(revenue...)

		// Sum total amount for LPs
		amountLPsCollected = amountLPsCollected.Add(revenuePortionForLPs...)
		amountStakersCollected = amountStakersCollected.Add(revenuePortionForStakers...)

		return false
	})

	return amountTotalCollected, amountLPsCollected, amountStakersCollected
}
