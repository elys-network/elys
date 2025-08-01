package keeper

import (
	"context"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) MatchAndExecuteOrders(goCtx context.Context, msg *types.MsgMatchAndExecuteOrders) (*types.MsgMatchAndExecuteOrdersResponse, error) {

	return &types.MsgMatchAndExecuteOrdersResponse{}, nil
}
