package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

func (oq *Querier) queryExitPoolEstimation(ctx sdk.Context, query *ammtypes.QueryExitPoolEstimationRequest) ([]byte, error) {
	res, err := oq.keeper.ExitPoolEstimation(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get exit pool estimation")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize exit pool estimation response")
	}
	return responseBytes, nil
}
