package estaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Create validators for Eden and EdenB
	if genState.Params.EdenCommitVal == "" {
		edenValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.Eden))
		k.Keeper.Hooks().AfterValidatorCreated(ctx, edenValAddr)
		genState.Params.EdenCommitVal = edenValAddr.String()
	}

	if genState.Params.EdenbCommitVal == "" {
		edenBValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.EdenB))
		k.Keeper.Hooks().AfterValidatorCreated(ctx, edenBValAddr)
		genState.Params.EdenbCommitVal = edenBValAddr.String()
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	// TODO: remove rewards management in incentive module (might be good to completely remove incentive module)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
