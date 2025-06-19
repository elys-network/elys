package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

func (oq *Querier) querySwapEstimation(ctx sdk.Context, query *ammtypes.QuerySwapEstimationRequest) ([]byte, error) {
	res, err := oq.keeper.SwapEstimation(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get swap estimation")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize swap estimation response")
	}
	return responseBytes, nil
}
