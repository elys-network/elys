package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

func (oq *Querier) queryOutRouteByDenom(ctx sdk.Context, query *ammtypes.QueryOutRouteByDenomRequest) ([]byte, error) {
	res, err := oq.keeper.OutRouteByDenom(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get out route by denom")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize out route by denom response")
	}
	return responseBytes, nil
}
