package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func (oq *Querier) queryEntry(ctx sdk.Context, req *types.QueryGetEntryRequest) ([]byte, error) {
	res, err := oq.keeper.Entry(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get entry")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize entry response")
	}
	return responseBytes, nil
}
