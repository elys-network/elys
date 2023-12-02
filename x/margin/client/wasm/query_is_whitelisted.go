package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (oq *Querier) queryIsWhitelisted(ctx sdk.Context, query *types.IsWhitelistedRequest) ([]byte, error) {
	res, err := oq.keeper.IsWhitelisted(sdk.WrapSDKContext(ctx), query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get is whitelisted")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize is whitelisted response")
	}
	return responseBytes, nil
}
