package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clockkeeper "github.com/elys-network/elys/x/clock/keeper"
	"github.com/elys-network/elys/x/clock/types"
)

func (oq *Querier) queryClockContracts(ctx sdk.Context, query *types.QueryClockContracts) ([]byte, error) {
	querier := clockkeeper.NewQuerier(*oq.keeper)
	res, err := querier.ClockContracts(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get clock contracts")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize clock contracts response")
	}
	return responseBytes, nil
}
