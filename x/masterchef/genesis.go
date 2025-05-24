package masterchef

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/masterchef/keeper"
	"github.com/elys-network/elys/v5/x/masterchef/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if k.CheckBlockedAddress(genState.Params) {
		panic("protocol revenues address is blocked")
	}

	k.SetParams(ctx, genState.Params)

	k.SetExternalIncentiveIndex(ctx, genState.ExternalIncentiveIndex)
	for _, elem := range genState.ExternalIncentives {
		k.SetExternalIncentive(ctx, elem)
	}

	for _, elem := range genState.PoolInfos {
		k.SetPoolInfo(ctx, elem)
	}

	for _, elem := range genState.PoolRewardInfos {
		k.SetPoolRewardInfo(ctx, elem)
	}

	for _, elem := range genState.UserRewardInfos {
		k.SetUserRewardInfo(ctx, elem)
	}

	for _, elem := range genState.PoolRewardsAccum {
		k.SetPoolRewardsAccum(ctx, elem)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.ExternalIncentives = k.GetAllExternalIncentives(ctx)
	genesis.ExternalIncentiveIndex = k.GetExternalIncentiveIndex(ctx)
	genesis.PoolInfos = k.GetAllPoolInfos(ctx)
	genesis.PoolRewardInfos = k.GetAllPoolRewardInfos(ctx)
	genesis.UserRewardInfos = k.GetAllUserRewardInfos(ctx)
	genesis.PoolRewardsAccum = k.GetAllPoolRewardsAccum(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
