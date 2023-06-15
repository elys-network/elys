package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	atypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Move gas fees collected to incentive module
func (k Keeper) CollectGasFeesToIncentiveModule(ctx sdk.Context) sdk.Coins {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	// Get pool Ids that can convert Elys
	poolIds := k.amm.GetAllPoolIdsWithDenom(ctx, ptypes.Elys)
	poolId := uint64(9999999)

	// Choose a pool that can convert Elys to USDC
	// TODO:
	// Later on: add a logic to choose best pool
	for _, pId := range poolIds {
		// Get a pool with poolId
		pool, found := k.amm.GetPool(ctx, pId)
		if !found {
			continue
		}

		// Loop pool assets to find out USDC pair
		for _, asset := range pool.PoolAssets {
			// if USDC available,
			if asset.Token.Denom == ptypes.USDC {
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
		return sdk.Coins{}
	}

	// Prepare route, considering every pool will have USDC pair
	// So use only 1 route. Elys to USDC
	routes := make([]atypes.SwapAmountInRoute, 0)
	route := atypes.SwapAmountInRoute{
		PoolId:        poolId,
		TokenOutDenom: ptypes.USDC,
	}

	// Routes
	routes = append(routes, route)

	// Amount in (Elys amount in sdk.Int)
	amtIn := feesCollectedInt.AmountOf(ptypes.Elys)
	// Elys token amount
	tokenIn := sdk.NewCoin(ptypes.Elys, amtIn)

	// Set zero to min out amount in order to have result all the time.
	tokenOutMinAmount := sdk.ZeroInt()

	// Convert Elys to USDC
	tokenOutAmount, err := k.amm.RouteExactAmountIn(ctx, feeCollector.GetAddress(), routes, tokenIn, tokenOutMinAmount)
	if err != nil {
		return sdk.Coins{}
	}

	// Swapped USDC coin
	swappedCoin := sdk.NewCoin(ptypes.USDC, tokenOutAmount)
	swappedCoins := sdk.NewCoins(swappedCoin)

	// Transfer converted USDC fees to the distribution module account
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, swappedCoins)
	if err != nil {
		panic(err)
	}

	return swappedCoins
}

// Pull DEX revenus collected to incentive module
// TODO:
// + transfer collected fees from different wallets(liquidity pool, margin module etc) to the distribution module account
// Assume this is already in USDC.
func (k Keeper) CollectDEXRevenusToIncentiveModule(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{}
}
