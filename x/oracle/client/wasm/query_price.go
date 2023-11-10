package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryPrice(ctx sdk.Context, req *oracletypes.QueryGetPriceRequest) ([]byte, error) {
	// Calling the PriceAll function and handling its response
	res, err := oq.keeper.Price(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get price")
	}

	// Serializing the response to a JSON byte array
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize price response")
	}

	return responseBytes, nil
}
