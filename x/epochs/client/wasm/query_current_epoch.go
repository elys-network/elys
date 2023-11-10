package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/epochs/types"
)

func (oq *Querier) queryCurrentEpoch(ctx sdk.Context, query *types.QueryCurrentEpochRequest) ([]byte, error) {
	res, err := oq.keeper.CurrentEpoch(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get current epoch")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize current epoch response")
	}
	return responseBytes, nil
}
