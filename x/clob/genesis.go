package clob

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/keeper"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, v := range genState.SubAccounts {
		k.SetSubAccount(ctx, v)
	}
	for _, v := range genState.PerpetualMarkets {
		k.SetPerpetualMarket(ctx, v)
	}
	for _, v := range genState.Perpetuals {
		k.SetPerpetual(ctx, v)
	}
	for _, v := range genState.PerpetualOwners {
		k.SetPerpetualOwner(ctx, v)
	}
	for _, v := range genState.OrderBooks {
		k.SetPerpetualOrder(ctx, v)
	}
	for _, v := range genState.OrderOwners {
		k.SetOrderOwner(ctx, v)
	}
	for _, v := range genState.TwapPrices {
		k.SetTwapPricesStruct(ctx, v)
	}
	for _, v := range genState.PerpetualMarketCounters {
		k.SetPerpetualMarketCounter(ctx, v)
	}
	for _, v := range genState.FundingRates {
		k.SetFundingRate(ctx, v)
	}
	for _, v := range genState.PerpetualADLs {
		k.SetPerpetualADL(ctx, v)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Params = k.GetParams(ctx)
	genesis.PerpetualMarkets = k.GetAllPerpetualMarket(ctx)
	genesis.PerpetualMarketCounters = k.GetAllPerpetualMarketCounter(ctx)
	genesis.SubAccounts = k.GetAllSubAccount(ctx)
	genesis.Perpetuals = k.GetAllPerpetuals(ctx)
	genesis.PerpetualOwners = k.GetAllPerpetualOwners(ctx)
	genesis.OrderBooks = k.GetAllPerpetualOrders(ctx)
	genesis.OrderOwners = k.GetAllOrderOwners(ctx)
	genesis.FundingRates = k.GetAllFundingRate(ctx)
	genesis.TwapPrices = k.GetAllTwapPrices(ctx)
	genesis.PerpetualADLs = k.GetAllPerpetualADLs(ctx)
	return genesis
}
