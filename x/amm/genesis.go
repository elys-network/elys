package amm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/keeper"
	"github.com/elys-network/elys/v5/x/amm/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pool
	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}
	// Set all the denomLiquidity
	for _, elem := range genState.DenomLiquidityList {
		k.SetDenomLiquidity(ctx, elem)
	}
	for _, track := range genState.SlippageTracks {
		k.SetSlippageTrack(ctx, track)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PoolList = k.GetAllPool(ctx)
	genesis.DenomLiquidityList = k.GetAllDenomLiquidity(ctx)
	genesis.SlippageTracks = k.AllSlippageTracks(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
