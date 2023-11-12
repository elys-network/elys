package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
)

func (oq *Querier) queryAccountedPoolAll(ctx sdk.Context, req *types.QueryAllAccountedPoolRequest) ([]byte, error) {
	res, err := oq.keeper.AccountedPoolAll(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get accounted pool all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize accounted pool all response")
	}
	return responseBytes, nil
}
