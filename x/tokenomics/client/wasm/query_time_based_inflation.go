package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (oq *Querier) queryTimeBasedInflation(ctx sdk.Context, query *types.QueryGetTimeBasedInflationRequest) ([]byte, error) {
	res, err := oq.keeper.TimeBasedInflation(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get time based inflation")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize time based inflation response")
	}
	return responseBytes, nil
}
