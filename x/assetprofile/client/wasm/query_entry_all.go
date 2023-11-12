package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func (oq *Querier) queryEntryAll(ctx sdk.Context, req *types.QueryAllEntryRequest) ([]byte, error) {
	res, err := oq.keeper.EntryAll(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get entry all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize entry all response")
	}
	return responseBytes, nil
}
