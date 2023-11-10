package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (oq *Querier) queryCommunityPool(ctx sdk.Context, query *types.QueryCommunityPoolRequest) ([]byte, error) {
	res, err := oq.keeper.CommunityPool(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get community pool")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize community pool response")
	}
	return responseBytes, nil
}
