package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
)

func (oq *Querier) querySwapEstimationByDenom(ctx sdk.Context, query *ammtypes.QuerySwapEstimationByDenomRequest) ([]byte, error) {
	res, err := oq.keeper.SwapEstimationByDenom(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get swap estimation by denom")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize swap estimation by denom response")
	}
	return responseBytes, nil
}
