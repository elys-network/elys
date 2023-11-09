package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (oq *Querier) queryParams(ctx sdk.Context, params *types.QueryParamsRequest) ([]byte, error) {
	// Your logic here
	return json.Marshal(&types.QueryParamsResponse{})
}
