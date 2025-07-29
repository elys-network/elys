package tradeshield

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tradeshield/keeper"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pendingSpotOrder
	for _, elem := range genState.PendingSpotOrderList {
		k.SetPendingSpotOrder(ctx, elem)
	}

	// Set pendingSpotOrder count
	k.SetPendingSpotOrderCount(ctx, genState.PendingSpotOrderCount)
	// Set all the pendingPerpetualOrder
	for _, elem := range genState.PendingPerpetualOrderList {
		k.SetPendingPerpetualOrder(ctx, elem)
	}

	// Set pendingPerpetualOrder count
	k.SetPendingPerpetualOrderCount(ctx, genState.PendingPerpetualOrderCount)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, &genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PendingSpotOrderList = k.GetAllPendingSpotOrder(ctx)
	genesis.PendingSpotOrderCount = k.GetPendingSpotOrderCount(ctx)
	genesis.PendingPerpetualOrderList = k.GetAllPendingPerpetualOrder(ctx)
	genesis.PendingPerpetualOrderCount = k.GetPendingPerpetualOrderCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
