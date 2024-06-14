package app

import (
	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func wasmConfiguration(params types.Params) {

	labelMaxSize := 256                // to set the maximum label size on instantiation (default 128)
	wasmMaxSize := 1638400             //  to set the max size of compiled wasm to be accepted (default 819200)
	wasmMaxProposalWasmSize := 6291456 // to set the max size of gov proposal compiled wasm to be accepted (default 3145728)

	if !params.WasmMaxLabelSize.IsNil() {
		labelMaxSize = int(params.WasmMaxLabelSize.Int64())
	}

	if !params.WasmMaxSize.IsNil() {
		wasmMaxSize = int(params.WasmMaxSize.Int64())
	}

	if !params.WasmMaxProposalWasmSize.IsNil() {
		wasmMaxProposalWasmSize = int(params.WasmMaxProposalWasmSize.Int64())
	}

	// increase wasm size limit
	wasm.MaxLabelSize = labelMaxSize
	wasm.MaxWasmSize = wasmMaxSize
	wasm.MaxProposalWasmSize = wasmMaxProposalWasmSize
}
