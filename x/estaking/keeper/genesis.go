package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/v7/x/estaking/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	var shouldRunEdenValHook = false
	var shouldRunEdenBValHook = false
	edenValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.Eden))
	edenBValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.EdenB))

	// Create validators for Eden and EdenB
	if genState.Params.EdenCommitVal == "" {
		genState.Params.EdenCommitVal = edenValAddr.String()
		shouldRunEdenValHook = true
	}

	if genState.Params.EdenbCommitVal == "" {
		genState.Params.EdenbCommitVal = edenBValAddr.String()
		shouldRunEdenBValHook = true
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	for _, snap := range genState.StakingSnapshots {
		k.SetElysStaked(ctx, snap)
	}

	if k.Hooks() != nil {
		if shouldRunEdenValHook {
			err := k.Hooks().AfterValidatorCreated(ctx, edenValAddr)
			if err != nil {
				panic(err)
			}
		}
		if shouldRunEdenBValHook {
			err := k.Hooks().AfterValidatorCreated(ctx, edenBValAddr)
			if err != nil {
				panic(err)
			}
		}
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.StakingSnapshots = k.GetAllElysStaked(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
