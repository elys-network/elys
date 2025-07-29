package membershiptier

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tier/keeper"
	"github.com/elys-network/elys/v7/x/tier/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the portfolio
	for _, elem := range genState.PortfolioList {
		k.SetPortfolio(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortfolioList = k.GetAllPortfolio(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
