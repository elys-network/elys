package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

func (oq *Querier) queryParams(ctx sdk.Context, query *types.QueryParamsRequest) ([]byte, error) {
	res, err := oq.keeper.Params(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get launchpad params")
	}
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize launchpad params")
	}
	return responseBytes, nil
}
