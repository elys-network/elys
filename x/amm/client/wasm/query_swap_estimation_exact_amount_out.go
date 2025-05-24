package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
)

func (oq *Querier) QuerySwapEstimationExactAmountOut(ctx sdk.Context, query *ammtypes.QuerySwapEstimationExactAmountOutRequest) ([]byte, error) {
	res, err := oq.keeper.SwapEstimationExactAmountOut(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get swap estimation exact amount out")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize swap estimation exact amount out response")
	}
	return responseBytes, nil
}
