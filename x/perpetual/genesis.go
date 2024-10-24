package perpetual

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pool
	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}

	// Set all the pool
	for _, elem := range genState.MtpList {
		err := k.SetMTP(ctx, &elem)
		if err != nil {
			panic(err)
		}
	}

	// Set genesis MTP count
	k.SetMTPCount(ctx, (uint64)(len(genState.MtpList)))
	// Set genesis open MTP count
	k.SetOpenMTPCount(ctx, (uint64)(len(genState.MtpList)))

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
	genesis.MtpList = k.GetAllMTPs(ctx)

	whitelist := k.GetAllWhitelistedAddress(ctx)
	whitelistAddresses := make([]string, len(whitelist))
	for i, whitelistAddress := range whitelist {
		whitelistAddresses[i] = whitelistAddress.String()
	}
	genesis.AddressWhitelist = whitelistAddresses

	return genesis
}
