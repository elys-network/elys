package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PendingSpotOrderAll(goCtx context.Context, req *types.QueryAllPendingSpotOrderRequest) (*types.QueryAllPendingSpotOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pendingSpotOrders []types.SpotOrder
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	pendingSpotOrderStore := prefix.NewStore(store, types.PendingSpotOrderKey)

	pageRes, err := query.Paginate(pendingSpotOrderStore, req.Pagination, func(key []byte, value []byte) error {
		var pendingSpotOrder types.SpotOrder
		if err := k.cdc.Unmarshal(value, &pendingSpotOrder); err != nil {
			return err
		}

		pendingSpotOrders = append(pendingSpotOrders, pendingSpotOrder)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPendingSpotOrderResponse{PendingSpotOrder: pendingSpotOrders, Pagination: pageRes}, nil
}

func (k Keeper) PendingSpotOrder(goCtx context.Context, req *types.QueryGetPendingSpotOrderRequest) (*types.QueryGetPendingSpotOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pendingSpotOrder, found := k.GetPendingSpotOrder(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetPendingSpotOrderResponse{PendingSpotOrder: pendingSpotOrder}, nil
}
