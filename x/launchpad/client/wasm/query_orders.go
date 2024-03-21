package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

func (oq *Querier) queryOrders(ctx sdk.Context, query *types.QueryOrdersRequest) ([]byte, error) {
	res, err := oq.keeper.Orders(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user orders")
	}
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user orders")
	}
	return responseBytes, nil
}
