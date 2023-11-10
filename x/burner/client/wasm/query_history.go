package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/burner/types"
)

func (oq *Querier) queryHistory(ctx sdk.Context, query *types.QueryGetHistoryRequest) ([]byte, error) {
	res, err := oq.keeper.History(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get history")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize history response")
	}
	return responseBytes, nil
}
