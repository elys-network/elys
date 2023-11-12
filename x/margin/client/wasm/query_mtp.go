package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (oq *Querier) queryMtp(ctx sdk.Context, query *types.MTPRequest) ([]byte, error) {
	res, err := oq.keeper.MTP(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get mtp")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize mtp response")
	}
	return responseBytes, nil
}
