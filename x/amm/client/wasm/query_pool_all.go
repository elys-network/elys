package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
)

func (oq *Querier) queryPoolAll(ctx sdk.Context, poolAll *types.QueryAllPoolRequest) ([]byte, error) {
	res, err := oq.keeper.PoolAll(ctx, poolAll)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pool all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pool all response")
	}
	return responseBytes, nil
}
