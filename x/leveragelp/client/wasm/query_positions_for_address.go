package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (oq *Querier) queryPositionsForAddress(ctx sdk.Context, query *types.PositionsForAddressRequest) ([]byte, error) {
	res, err := oq.keeper.QueryPositionsForAddress(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get positions for address")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize positions for address response")
	}
	return responseBytes, nil
}
