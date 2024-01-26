package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func (oq *Querier) queryEntryByDenom(ctx sdk.Context, req *types.QueryGetEntryByDenomRequest) ([]byte, error) {
	res, err := oq.keeper.EntryByDenom(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get entry by denom")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize entry by denom response")
	}
	return responseBytes, nil
}
