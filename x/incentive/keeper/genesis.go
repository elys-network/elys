package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetFeePool(ctx, data.FeePool)
	k.SetParams(ctx, data.Params)
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	feePool := k.GetFeePool(ctx)
	params := k.GetParams(ctx)

	return types.NewGenesisState(params, feePool)
}
