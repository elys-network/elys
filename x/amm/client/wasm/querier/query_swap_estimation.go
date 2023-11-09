package wasm

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/client/wasm/types"
)

func (oq *Querier) querySwapEstimation(ctx sdk.Context, query *types.QuerySwapEstimationRequest) ([]byte, error) {
	return nil, wasmvmtypes.UnsupportedRequest{Kind: "QuerySwapEstimation, not implemented yet"}
}
