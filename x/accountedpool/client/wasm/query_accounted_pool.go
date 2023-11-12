package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
)

func (oq *Querier) queryAccountedPool(ctx sdk.Context, req *types.QueryGetAccountedPoolRequest) ([]byte, error) {
	res, err := oq.keeper.AccountedPool(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get accounted pool")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize accounted pool response")
	}
	return responseBytes, nil
}
