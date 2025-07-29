package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (oq *Querier) querySlippageTrackAll(ctx sdk.Context, slippageTrackAll *types.QuerySlippageTrackAllRequest) ([]byte, error) {
	res, err := oq.keeper.SlippageTrackAll(ctx, slippageTrackAll)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get slippage track all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize slippage track all response")
	}
	return responseBytes, nil
}
