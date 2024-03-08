package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (oq *Querier) queryApr(ctx sdk.Context, query *types.QueryAprRequest) ([]byte, error) {
	res, err := oq.keeper.CalculateApr(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get community pool")
	}

	resp := types.QueryAprResponse{
		Apr: res,
	}
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize apr response")
	}
	return responseBytes, nil
}
