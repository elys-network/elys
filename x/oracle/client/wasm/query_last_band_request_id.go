package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryLastBandRequestId(ctx sdk.Context, req *oracletypes.QueryLastBandRequestIdRequest) ([]byte, error) {
	res, err := oq.keeper.LastBandRequestId(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query last band request")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize last band request response")
	}
	return responseBytes, nil
}
