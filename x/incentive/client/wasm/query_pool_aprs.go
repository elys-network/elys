package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (oq *Querier) queryPoolAprs(ctx sdk.Context, query *types.QueryPoolAprsRequest) ([]byte, error) {
	data := oq.keeper.CalculatePoolAprs(ctx, query.PoolIds)

	resp := types.QueryPoolAprsResponse{
		Data: data,
	}
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pool aprs response")
	}
	return responseBytes, nil
}
