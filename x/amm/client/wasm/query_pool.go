package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (oq *Querier) queryPool(ctx sdk.Context, pool *types.QueryGetPoolRequest) ([]byte, error) {
	res, err := oq.keeper.Pool(ctx, pool)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pool")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pool response")
	}
	return responseBytes, nil
}
