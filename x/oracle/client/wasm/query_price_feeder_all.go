package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryPriceFeederAll(ctx sdk.Context, req *oracletypes.QueryAllPriceFeederRequest) ([]byte, error) {
	// Calling the PriceAll function and handling its response
	res, err := oq.keeper.PriceFeederAll(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get all price feeders")
	}

	// Serializing the response to a JSON byte array
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize all price feeders response")
	}

	return responseBytes, nil
}
