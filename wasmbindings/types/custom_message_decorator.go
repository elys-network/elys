package types

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
)

func CustomMessageDecorator(
	moduleMessengers []ModuleMessenger,
	amm *ammkeeper.Keeper,
	auth *authkeeper.AccountKeeper,
	bank *bankkeeper.BaseKeeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:          old,
			moduleMessengers: moduleMessengers,
			amm:              amm,
			auth:             auth,
			bank:             bank,
		}
	}
}
