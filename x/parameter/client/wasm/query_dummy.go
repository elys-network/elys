package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryDummy(ctx sdk.Context, query *types.QueryDummy) ([]byte, error) {
	res, err := oq.keeper.Dummy(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get dummy")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize dummy response")
	}
	return responseBytes, nil
}
