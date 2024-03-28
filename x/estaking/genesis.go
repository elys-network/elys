package estaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/estaking/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	// TODO: create validators for Eden and EdenB - inactive validators
	// TODO: configure params if no validators set
	// TODO: there could be invariant checker to check not_bonded_pool balance
	// and total not bonded validators sum
	// TODO: prevent direct interaction to Eden and EdenB validators
	// through staking module - antehandler & wasmbinding restriction

	// TODO: remove rewards management in incentive module (might be good to completely remove incentive module)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
