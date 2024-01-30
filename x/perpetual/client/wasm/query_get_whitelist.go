package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (oq *Querier) queryGetWhitelist(ctx sdk.Context, query *types.WhitelistRequest) ([]byte, error) {
	res, err := oq.keeper.GetWhitelist(sdk.WrapSDKContext(ctx), query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get whitelist")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize whitelist response")
	}
	return responseBytes, nil
}
