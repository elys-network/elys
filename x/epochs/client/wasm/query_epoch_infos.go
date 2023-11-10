package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/epochs/types"
)

func (oq *Querier) queryEpochInfos(ctx sdk.Context, query *types.QueryEpochsInfoRequest) ([]byte, error) {
	res, err := oq.keeper.EpochInfos(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get epoch infos")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize epoch infos response")
	}
	return responseBytes, nil
}
