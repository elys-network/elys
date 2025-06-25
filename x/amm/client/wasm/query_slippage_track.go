package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
)

func (oq *Querier) querySlippageTrack(ctx sdk.Context, slippageTrack *types.QuerySlippageTrackRequest) ([]byte, error) {
	res, err := oq.keeper.SlippageTrack(ctx, slippageTrack)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get slippage track")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize slippage track response")
	}
	return responseBytes, nil
}
