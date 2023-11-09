package wasm

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (oq *Querier) querySwapEstimation(ctx sdk.Context, query *ammtypes.QuerySwapEstimationRequest) ([]byte, error) {
	return nil, wasmvmtypes.UnsupportedRequest{Kind: "QuerySwapEstimation, not implemented yet"}
}
