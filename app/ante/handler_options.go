package ante

import (
	corestoretypes "cosmossdk.io/core/store"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibcconsumerkeeper "github.com/cosmos/interchain-security/v6/x/ccv/consumer/keeper"
	parameterkeeper "github.com/elys-network/elys/v6/x/parameter/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	Cdc                   codec.BinaryCodec
	IBCKeeper             *ibckeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	ConsumerKeeper        ibcconsumerkeeper.Keeper
	BankKeeper            bankkeeper.Keeper
	ParameterKeeper       parameterkeeper.Keeper
	WasmConfig            *wasmTypes.WasmConfig
	TXCounterStoreService corestoretypes.KVStoreService
}
