package client

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	externalsauth "github.com/elys-network/elys/v6/wasmbindings/externals/auth"
	"github.com/elys-network/elys/v6/wasmbindings/types"
	ammclientwasm "github.com/elys-network/elys/v6/x/amm/client/wasm"
	ammkeeper "github.com/elys-network/elys/v6/x/amm/keeper"
)

func RegisterCustomPlugins(
	amm *ammkeeper.Keeper,
	auth *authkeeper.AccountKeeper,
	bank *bankkeeper.BaseKeeper,
) []wasmkeeper.Option {
	ammQuerier := ammclientwasm.NewQuerier(amm)
	ammMessenger := ammclientwasm.NewMessenger(amm)

	authQuerier := externalsauth.NewQuerier(auth)
	authMessenger := externalsauth.NewMessenger(auth)

	moduleQueriers := []types.ModuleQuerier{
		ammQuerier,
		authQuerier,
	}

	wasmQueryPlugin := types.NewQueryPlugin(
		moduleQueriers,
		amm,
		auth,
		bank,
	)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: types.CustomQuerier(wasmQueryPlugin),
	})

	moduleMessengers := []types.ModuleMessenger{
		ammMessenger,
		authMessenger,
	}

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		types.CustomMessageDecorator(
			moduleMessengers,
			amm,
			auth,
			bank,
		),
	)
	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
