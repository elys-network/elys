package vaults

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v7/x/vaults/keeper"
	"github.com/elys-network/elys/v7/x/vaults/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the vaults
	for _, elem := range genState.VaultList {
		if err := k.SetVault(ctx, elem); err != nil {
			panic(err)
		}
	}
	// Set all the pool info
	for _, elem := range genState.PoolInfoList {
		k.SetPoolInfo(ctx, elem)
	}
	// Set all the pool reward info
	for _, elem := range genState.PoolRewardInfoList {
		k.SetPoolRewardInfo(ctx, elem)
	}
	// Set all the user reward info
	for _, elem := range genState.UserRewardInfoList {
		k.SetUserRewardInfo(ctx, elem)
	}
	// Set all the pool rewards accum
	for _, elem := range genState.PoolRewardsAccumList {
		k.SetPoolRewardsAccum(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.VaultList = k.GetAllVaults(ctx)
	genesis.PoolInfoList = k.GetAllPoolInfos(ctx)
	genesis.PoolRewardInfoList = k.GetAllPoolRewardInfos(ctx)
	genesis.UserRewardInfoList = k.GetAllUserRewardInfos(ctx)
	genesis.PoolRewardsAccumList = k.GetAllPoolRewardsAccum(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
