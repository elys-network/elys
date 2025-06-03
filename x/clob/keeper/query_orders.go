package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) OrderBook(goCtx context.Context, req *types.OrderBookRequest) (*types.OrderBookResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var list []types.PerpetualOrder

	iterator := k.GetSellOrderIterator(ctx, req.MarketId)
	if req.IsBuy {
		iterator = k.GetBuyOrderIterator(ctx, req.MarketId)
	}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return &types.OrderBookResponse{Orders: list}, nil
}
