package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/burner/types"
)

func (oq *Querier) queryHistoryAll(ctx sdk.Context, query *types.QueryAllHistoryRequest) ([]byte, error) {
	res, err := oq.keeper.HistoryAll(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get history all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize history all response")
	}
	return responseBytes, nil
}
