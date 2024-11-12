package ante

import (
	corestoretypes "cosmossdk.io/core/store"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	Cdc                   codec.BinaryCodec
	IBCKeeper             *ibckeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
	WasmConfig            *wasmtypes.WasmConfig

	BankKeeper      bankkeeper.Keeper
	ParameterKeeper parameterkeeper.Keeper
}
