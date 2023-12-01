package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (oq *Querier) queryMinCollateral(ctx sdk.Context, query *types.QueryMinCollateralRequest) ([]byte, error) {
	res, err := oq.keeper.MinCollateral(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get min collateral")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize min collateral response")
	}
	return responseBytes, nil
}
