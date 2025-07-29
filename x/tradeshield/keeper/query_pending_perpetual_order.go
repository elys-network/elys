package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PendingPerpetualOrderAll(goCtx context.Context, req *types.QueryAllPendingPerpetualOrderRequest) (*types.QueryAllPendingPerpetualOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pendingPerpetualOrdersExtraInfo []types.PerpetualOrderExtraInfo
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	pendingPerpetualOrderStore := prefix.NewStore(store, types.PendingPerpetualOrderKey)

	pageRes, err := query.Paginate(pendingPerpetualOrderStore, req.Pagination, func(key []byte, value []byte) error {
		var pendingPerpetualOrder types.PerpetualOrder
		if err := k.cdc.Unmarshal(value, &pendingPerpetualOrder); err != nil {
			return err
		}

		pendingPerpetualOrderExtraInfo, err := k.ConstructPerpetualOrderExtraInfo(ctx, pendingPerpetualOrder)
		if err != nil {
			return err
		}

		pendingPerpetualOrdersExtraInfo = append(pendingPerpetualOrdersExtraInfo, *pendingPerpetualOrderExtraInfo)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPendingPerpetualOrderResponse{PendingPerpetualOrder: pendingPerpetualOrdersExtraInfo, Pagination: pageRes}, nil
}

func (k Keeper) PendingPerpetualOrder(goCtx context.Context, req *types.QueryGetPendingPerpetualOrderRequest) (*types.QueryGetPendingPerpetualOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pendingPerpetualOrder, found := k.GetPendingPerpetualOrder(ctx, sdk.MustAccAddressFromBech32(req.OwnerAddress), req.PoolId, req.OrderId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	pendingPerpetualOrderExtraInfo, err := k.ConstructPerpetualOrderExtraInfo(ctx, pendingPerpetualOrder)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPendingPerpetualOrderResponse{PendingPerpetualOrder: *pendingPerpetualOrderExtraInfo}, nil
}
