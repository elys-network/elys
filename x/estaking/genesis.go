package estaking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	var shouldRunEdenValHook = false
	var shouldRunEdenBValHook = false
	edenValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.Eden))
	edenBValAddr := sdk.ValAddress(authtypes.NewModuleAddress(ptypes.EdenB))
	fmt.Println("k.Keeper", k.Keeper)
	fmt.Println("k", k)

	// 	commKeeper         types.CommitmentKeeper
	// distrKeeper        types.DistrKeeper
	// tokenomicsKeeper   types.TokenomicsKeeper
	// assetProfileKeeper types.AssetProfileKeeper
	// authority          string

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
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
