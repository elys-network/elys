package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryPriceFeeder(ctx sdk.Context, req *oracletypes.QueryGetPriceFeederRequest) ([]byte, error) {
	// Calling the PriceAll function and handling its response
	res, err := oq.keeper.PriceFeeder(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get price feeder")
	}

	// Serializing the response to a JSON byte array
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize price feeder response")
	}

	return responseBytes, nil
}
