package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
)

func (oq *Querier) queryJoinPoolEstimation(ctx sdk.Context, query *ammtypes.QueryJoinPoolEstimationRequest) ([]byte, error) {
	res, err := oq.keeper.JoinPoolEstimation(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to join pool estimation")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize join pool estimation response")
	}
	return responseBytes, nil
}
