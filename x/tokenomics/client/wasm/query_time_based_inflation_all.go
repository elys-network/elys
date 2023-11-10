package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (oq *Querier) queryTimeBasedInflationAll(ctx sdk.Context, query *types.QueryAllTimeBasedInflationRequest) ([]byte, error) {
	res, err := oq.keeper.TimeBasedInflationAll(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get time based inflation all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize time based inflation all response")
	}
	return responseBytes, nil
}
