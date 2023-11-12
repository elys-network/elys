package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryAnteHandlerParamAll(ctx sdk.Context, query *types.QueryAllAnteHandlerParamRequest) ([]byte, error) {
	res, err := oq.keeper.AnteHandlerParamAll(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get ante handler param all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize ante handler param all response")
	}
	return responseBytes, nil
}
