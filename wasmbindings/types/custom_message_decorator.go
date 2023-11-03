package types

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
)

func CustomMessageDecorator(moduleMessengers []ModuleMessenger, amm *ammkeeper.Keeper, margin *marginkeeper.Keeper, staking *stakingkeeper.Keeper, commitment *commitmentkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:          old,
			moduleMessengers: moduleMessengers,
			amm:              amm,
			margin:           margin,
			staking:          staking,
			commitment:       commitment,
		}
	}
}
