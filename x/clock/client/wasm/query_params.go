package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clockkeeper "github.com/elys-network/elys/x/clock/keeper"
	"github.com/elys-network/elys/x/clock/types"
)

func (oq *Querier) queryParams(ctx sdk.Context, query *types.QueryParamsRequest) ([]byte, error) {
	querier := clockkeeper.NewQuerier(*oq.keeper)
	res, err := querier.Params(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get params")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize params response")
	}
	return responseBytes, nil
}
