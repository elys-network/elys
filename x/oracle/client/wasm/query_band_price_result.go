package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryBandPriceResult(ctx sdk.Context, req *oracletypes.QueryBandPriceRequest) ([]byte, error) {
	res, err := oq.keeper.BandPriceResult(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query band price result")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize band price result response")
	}
	return responseBytes, nil
}
