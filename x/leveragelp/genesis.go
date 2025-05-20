package leveragelp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/leveragelp/keeper"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pool
	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}

	// Set all the pool
	for _, elem := range genState.PositionList {
		k.SetPosition(ctx, &elem)
	}

	// Set genesis Position count
	k.SetPositionCount(ctx, (uint64)(len(genState.PositionList)))
	// Set genesis open Position count
	k.SetOpenPositionCount(ctx, (uint64)(len(genState.PositionList)))

	// Set all the whitelisted
	for _, elem := range genState.AddressWhitelist {
		k.WhitelistAddress(ctx, sdk.MustAccAddressFromBech32(elem))
	}

	// this line is used by starport scaffolding # genesis/module/init
	err := k.SetParams(ctx, &genState.Params)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PoolList = k.GetAllPools(ctx)
	genesis.PositionList = k.GetAllPositions(ctx)
	whitelist := k.GetAllWhitelistedAddress(ctx)
	whitelistAddressStrings := make([]string, len(whitelist))
	for i, whitelistAddress := range whitelist {
		whitelistAddressStrings[i] = whitelistAddress.String()
	}
	genesis.AddressWhitelist = whitelistAddressStrings

	return genesis
}
