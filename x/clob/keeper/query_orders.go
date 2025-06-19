package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OrderBook(goCtx context.Context, req *types.OrderBookRequest) (*types.OrderBookResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var list []types.PerpetualOrder

	key := types.GetPerpetualOrderBookIteratorKey(req.MarketId, false)
	if req.IsBuy {
		key = types.GetPerpetualOrderBookIteratorKey(req.MarketId, true)
	}

	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var val types.PerpetualOrder
		if err := k.cdc.Unmarshal(value, &val); err != nil {
			return err
		}

		list = append(list, val)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.OrderBookResponse{Orders: list, Pagination: pageRes}, nil
}
