package wasm

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
)

func (oq *Querier) querySwapEstimation(ctx sdk.Context, query *wasmbindingstypes.QuerySwapEstimationRequest) ([]byte, error) {
	return nil, wasmvmtypes.UnsupportedRequest{Kind: "QuerySwapEstimation, not implemented yet"}
}
