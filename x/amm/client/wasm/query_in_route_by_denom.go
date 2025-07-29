package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
)

func (oq *Querier) queryInRouteByDenom(ctx sdk.Context, query *ammtypes.QueryInRouteByDenomRequest) ([]byte, error) {
	res, err := oq.keeper.InRouteByDenom(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get in route by denom")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize in route by denom response")
	}
	return responseBytes, nil
}
