package ante

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// HandlerOptions extends the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	sdkante.HandlerOptions
	Cdc               codec.BinaryCodec
	StakingKeeper     *stakingkeeper.Keeper
	BankKeeper        bankkeeper.Keeper
	IBCKeeper         *ibckeeper.Keeper
	WasmConfig        *wasmtypes.WasmConfig
	ParameterKeeper   parameterkeeper.Keeper
	TXCounterStoreKey storetypes.StoreKey
}
