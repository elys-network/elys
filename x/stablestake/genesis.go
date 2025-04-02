package stablestake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the debts
	for _, elem := range genState.DebtList {
		k.SetDebt(ctx, elem)
	}

	// Set all the interests
	for _, elem := range genState.InterestList {
		k.SetInterestForPool(ctx, elem)
	}

	// Set all pools
	for _, elem := range genState.Pools {
		k.SetPool(ctx, elem)
	}

	// Set all amm pools
	for _, elem := range genState.AmmPools {
		k.SetAmmPool(ctx, elem)
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.DebtList = k.GetAllDebts(ctx)
	genesis.InterestList = k.GetAllInterest(ctx)

	genesis.Pools = k.GetAllPools(ctx)
	genesis.AmmPools = k.GetAllAmmPools(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
