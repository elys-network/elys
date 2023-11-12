package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (oq *Querier) queryPosition(ctx sdk.Context, query *types.PositionRequest) ([]byte, error) {
	res, err := oq.keeper.Position(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get position")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize position response")
	}
	return responseBytes, nil
}
