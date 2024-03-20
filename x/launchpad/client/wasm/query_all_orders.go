package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

func (oq *Querier) queryAllOrders(ctx sdk.Context, query *types.QueryAllOrdersRequest) ([]byte, error) {
	orders := oq.keeper.GetAllOrders(ctx)
	res := types.QueryAllOrdersResponse{
		Purchases: orders,
	}
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize all orders")
	}
	return responseBytes, nil
}
