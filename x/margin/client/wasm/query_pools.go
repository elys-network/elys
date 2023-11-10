package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (oq *Querier) queryPools(ctx sdk.Context, query *types.QueryAllPoolRequest) ([]byte, error) {
	res, err := oq.keeper.Pools(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pools")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pools response")
	}
	return responseBytes, nil
}
