package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (oq *Querier) queryPositions(ctx sdk.Context, query *types.PositionsRequest) ([]byte, error) {
	res, err := oq.keeper.GetPositions(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get positions")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize positions response")
	}
	return responseBytes, nil
}
