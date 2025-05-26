package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PerpetualOrder(goCtx context.Context, req *types.PerpetualOrderRequest) (*types.PerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	perpetualOrder, err := k.GetPerpetualOrder(ctx, req.MarketId, req.OrderType, req.Price, req.BlockHeight)

	if err != nil {
		return nil, err
	}

	return &types.PerpetualOrderResponse{Order: perpetualOrder}, nil
}

func (k Keeper) AllPerpetualOrder(goCtx context.Context, req *types.AllPerpetualOrderRequest) (*types.AllPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var orders []types.PerpetualOrder

	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOrderPrefix)

	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var order types.PerpetualOrder
		if err := k.cdc.Unmarshal(value, &order); err != nil {
			return err
		}

		orders = append(orders, order)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.AllPerpetualOrderResponse{Orders: orders, Pagination: pageRes}, nil
}
