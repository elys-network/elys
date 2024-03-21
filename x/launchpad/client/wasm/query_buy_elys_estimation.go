package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

func (oq *Querier) queryBuyElysEst(ctx sdk.Context, query *types.QueryBuyElysEstRequest) ([]byte, error) {
	res, err := oq.keeper.BuyElysEst(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user bonus")
	}
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user bonus")
	}
	return responseBytes, nil
}
