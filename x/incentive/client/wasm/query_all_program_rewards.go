package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (oq *Querier) queryAllProgramRewards(ctx sdk.Context, query *types.QueryAllProgramRewardsRequest) ([]byte, error) {
	resp, err := oq.keeper.AllProgramRewards(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get all program rewards")
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize all program rewards response")
	}
	return responseBytes, nil
}
