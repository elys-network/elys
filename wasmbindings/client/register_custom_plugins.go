package client

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/elys-network/elys/wasmbindings/types"
	ammclientwasm "github.com/elys-network/elys/x/amm/client/wasm"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	commitmentclientwasm "github.com/elys-network/elys/x/commitment/client/wasm"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	oracleclientwasm "github.com/elys-network/elys/x/oracle/client/wasm"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
)

func RegisterCustomPlugins(
	amm *ammkeeper.Keeper,
	oracle *oraclekeeper.Keeper,
	margin *marginkeeper.Keeper,
	bank *bankkeeper.BaseKeeper,
	staking *stakingkeeper.Keeper,
	commitment *commitmentkeeper.Keeper,
) []wasmkeeper.Option {
	oracleQuerier := oracleclientwasm.NewQuerier(oracle)
	oracleMessenger := oracleclientwasm.NewMessenger(oracle)

	ammQuerier := ammclientwasm.NewQuerier(amm, bank, commitment)
	ammMessenger := ammclientwasm.NewMessenger(amm)

	commitmentQuerier := commitmentclientwasm.NewQuerier(commitment)
	commitmentMessenger := commitmentclientwasm.NewMessenger(commitment, staking)

	moduleQueriers := []types.ModuleQuerier{
		oracleQuerier,
		ammQuerier,
		commitmentQuerier,
	}

	wasmQueryPlugin := types.NewQueryPlugin(moduleQueriers, amm, oracle, bank, staking, commitment, margin)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: types.CustomQuerier(wasmQueryPlugin),
	})

	moduleMessengers := []types.ModuleMessenger{
		oracleMessenger,
		ammMessenger,
		commitmentMessenger,
	}

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		types.CustomMessageDecorator(moduleMessengers, amm, margin, staking, commitment),
	)
	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
