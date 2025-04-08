package commitment

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/commitment/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	// Set all the assetInfo
	for _, commitment := range genState.Commitments {
		k.SetCommitments(ctx, *commitment)
	}
	for _, val := range genState.AtomStakers {
		k.SetAtomStaker(ctx, *val)
	}
	for _, val := range genState.NftHolders {
		k.SetNFTHodler(ctx, *val)
	}
	for _, val := range genState.Governors {
		k.SetGovernor(ctx, *val)
	}
	for _, val := range genState.Cadets {
		k.SetCadet(ctx, *val)
	}
	for _, val := range genState.KolList {
		k.SetKol(ctx, *val)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.Commitments = k.GetAllCommitments(ctx)

	genesis.AtomStakers = k.GetAllAtomStakers(ctx)
	genesis.Cadets = k.GetAllCadets(ctx)
	genesis.Governors = k.GetAllGovernors(ctx)
	genesis.NftHolders = k.GetAllNFTHolders(ctx)
	genesis.KolList = k.GetAllKol(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
