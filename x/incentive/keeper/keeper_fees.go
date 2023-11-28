package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// FindPool function gets a pool that can convert in_denom token to out_denom token
// TODO:
// Later on: add a logic to choose best pool
func (k Keeper) FindPool(ctx sdk.Context, in_denom string, out_denom string) (ammtypes.Pool, bool) {
	// Get pool ids that can convert tokenIn
	poolIds := k.amm.GetAllPoolIdsWithDenom(ctx, in_denom)
	poolId := uint64(9999999)

	for _, pId := range poolIds {
		// Get a pool with poolId
		pool, found := k.amm.GetPool(ctx, pId)
		if !found {
			continue
		}

		// Loop pool assets to find out USDC pair
		for _, asset := range pool.PoolAssets {
			// if USDC available,
			if asset.Token.Denom == out_denom {
				poolId = pool.PoolId
				break
			}
		}

		// If already found a pool matched, exit loop
		if poolId != uint64(9999999) {
			break
		}
	}

	// If the pool is not availble,
	if poolId == uint64(9999999) {
		return ammtypes.Pool{}, false
	}

	// Return a pool found
	return k.amm.GetPool(ctx, poolId)
}

// Move gas fees collected to dex revenue wallet
// Convert it into USDC
func (k Keeper) CollectGasFeesToIncentiveModule(ctx sdk.Context, baseCurrency string) sdk.Coins {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	// Total Swapped coin
	totalSwappedCoins := sdk.Coins{}

	for _, tokenIn := range feesCollectedInt {
		// skip for fee denom - usdc
		if tokenIn.Denom == baseCurrency {
			continue
		}

		// Find a pool that can convert tokenIn to usdc
		pool, found := k.FindPool(ctx, tokenIn.Denom, baseCurrency)
		if !found {
			continue
		}

		// Executes the swap in the pool and stores the output. Updates pool assets but
		// does not actually transfer any tokens to or from the pool.
		snapshot := k.amm.GetPoolSnapshotOrSet(ctx, pool)
		tokenOutCoin, _, _, err := k.amm.SwapOutAmtGivenIn(ctx, pool.PoolId, k.oracleKeeper, &snapshot, sdk.Coins{tokenIn}, baseCurrency, sdk.ZeroDec())
		if err != nil {
			continue
		}

		tokenOutAmount := tokenOutCoin.Amount

		if !tokenOutAmount.IsPositive() {
			continue
		}

		// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
		// Also emits a swap event and updates related liquidity metrics.
		err, _ = k.amm.UpdatePoolForSwap(ctx, pool, feeCollector.GetAddress(), feeCollector.GetAddress(), tokenIn, tokenOutCoin, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
		if err != nil {
			continue
		}

		// Swapped USDC coin
		swappedCoin := sdk.NewCoin(baseCurrency, tokenOutAmount)
		swappedCoins := sdk.NewCoins(swappedCoin)

		// Transfer converted USDC fees to the Dex revenue module account
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, k.dexRevCollectorName, swappedCoins)
		if err != nil {
			panic(err)
		}

		// Sum total swapped
		totalSwappedCoins = totalSwappedCoins.Add(swappedCoins...)
	}

	return totalSwappedCoins
}

// Collect all DEX revenues to DEX revenue wallet,
// while tracking the 65% of it for LPs reward distribution
// transfer collected fees from different wallets(liquidity pool, margin module etc) to the distribution module account
// Assume this is already in USDC.
// TODO:
// + Collect revenue from margin, lend module
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

		// Revenue amount
		revenue := k.bankKeeper.GetAllBalances(ctx, revenueAddress)

		denoms := revenue.Denoms()
		if len(denoms) < 1 {
			return false
		}

		// Transfer revenue to a single wallet of DEX revenue wallet.
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, revenueAddress, k.dexRevCollectorName, revenue)
		if err != nil {
			panic(err)
		}

		// LPs Portion param
		rewardPortionForLps := k.GetDEXRewardPortionForLPs(ctx)

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
