package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the assetInfo
	for _, elem := range genState.AssetInfos {
		k.SetAssetInfo(ctx, elem)
	}
	// Set all the price
	for _, elem := range genState.Prices {
		k.SetPrice(ctx, elem)
	}
	// Set all the priceFeeder
	for _, elem := range genState.PriceFeeders {
		k.SetPriceFeeder(ctx, elem)
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	//genesis.PortId = k.GetPort(ctx)
	genesis.AssetInfos = k.GetAllAssetInfo(ctx)
	genesis.Prices = k.GetAllLegacyPrice(ctx)
	genesis.PriceFeeders = k.GetAllPriceFeeder(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
