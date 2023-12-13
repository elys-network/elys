package tvl

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/keeper"
	"github.com/elys-network/elys/x/accountedpool/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the accountedPool
	for _, elem := range genState.AccountedPoolList {
		k.SetAccountedPool(ctx, elem)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.AccountedPoolList = k.GetAllAccountedPool(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
