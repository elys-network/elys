package app

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
	}
}

type QueryPlugin struct {
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin() *QueryPlugin {
	return &QueryPlugin{}
}

func RegisterCustomPlugins() []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin()

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	return []wasm.Option{
		queryPluginOpt,
	}
}

// CustomQuerier dispatches custom CosmWasm bindings queries.
func CustomQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery ElysQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "elys query")
		}

		switch {
		case contractQuery.PriceAll != nil:
			pagination := contractQuery.PriceAll.Pagination

			_ = pagination

			res := AllPriceResponse{
				Price: []oracletypes.Price{},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   0,
				},
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "elys price all query response")
			}

			return bz, nil

		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown elys query variant"}
		}
	}
}

type ElysQuery struct {
	PriceAll *PriceAll `json:"price_all,omitempty"`
}

type PriceAll struct {
	// oracletypes.QueryAllPriceRequest
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type AllPriceResponse struct {
	// oracletypes.QueryAllPriceResponse
	Price      []oracletypes.Price `protobuf:"bytes,1,rep,name=price,proto3" json:"price"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
