package tokenomics

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the airdrop
	for _, elem := range genState.AirdropList {
		k.SetAirdrop(ctx, elem)
	}
	// Set if defined
	if genState.GenesisInflation != nil {
		k.SetGenesisInflation(ctx, *genState.GenesisInflation)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.AirdropList = k.GetAllAirdrop(ctx)
	// Get all genesisInflation
	genesisInflation, found := k.GetGenesisInflation(ctx)
	if found {
		genesis.GenesisInflation = &genesisInflation
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
