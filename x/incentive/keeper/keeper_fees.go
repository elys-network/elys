package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Move gas fees collected to incentive module
func (k Keeper) CollectGasFeesToIncentiveModule(ctx sdk.Context) sdk.Coins {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	return feesCollectedInt
}

// Pull DEX revenus collected to incentive module
// TODO:
// + transfer collected fees from different wallets(liquidity pool, margin module etc) to the distribution module account
func (k Keeper) CollectDEXRevenusToIncentiveModule(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{}
}
