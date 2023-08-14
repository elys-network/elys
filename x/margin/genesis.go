package margin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pool
	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}

	// Set all the pool
	for _, elem := range genState.MtpList {
		k.SetMTP(ctx, &elem)
	}

	// Set genesis MTP count
	k.SetMTPCount(ctx, (uint64)(len(genState.MtpList)))
	// Set genesis open MTP count
	k.SetOpenMTPCount(ctx, (uint64)(len(genState.MtpList)))

	// Set all the whitelisted
	for _, elem := range genState.AddressWhitelist {
		k.WhitelistAddress(ctx, elem)
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, &genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PoolList = k.GetAllPools(ctx)
	genesis.MtpList = k.GetAllMTPs(ctx)
	genesis.AddressWhitelist = k.GetAllWhitelistedAddress(ctx)

	return genesis
}
