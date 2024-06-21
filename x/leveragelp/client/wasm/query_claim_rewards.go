package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (oq *Querier) queryRewards(ctx sdk.Context, query *types.QueryRewardsRequest) ([]byte, error) {
	res, err := oq.keeper.Rewards(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get leveragelp rewards")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize leveragelp rewards response")
	}
	return responseBytes, nil
}
