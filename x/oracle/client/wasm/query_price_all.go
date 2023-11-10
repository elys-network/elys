package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryPriceAll(ctx sdk.Context, req *oracletypes.QueryAllPriceRequest) ([]byte, error) {
	// Calling the PriceAll function and handling its response
	res, err := oq.keeper.PriceAll(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get all prices")
	}

	// Serializing the response to a JSON byte array
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize all prices response")
	}

	return responseBytes, nil
}
