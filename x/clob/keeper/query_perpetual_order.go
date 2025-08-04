package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OwnerOrder(goCtx context.Context, request *types.OwnerOrdersRequest) (*types.OwnerOrdersResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	var orders []types.Order
	var prefixStore prefix.Store

	if request.SubAccountId == 0 {
		prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetOrderOwnerAddressKey(owner))
	} else {
		key := types.GetOrderOwnerAddressKey(owner)
		key = append(key, sdk.Uint64ToBigEndian(request.SubAccountId)...)
		key = append(key, []byte("/")...)
		prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)

	}

	pageRes, err := query.Paginate(prefixStore, request.Pagination, func(key []byte, value []byte) error {
		var orderOwner types.PerpetualOrderOwner
		if err := k.cdc.Unmarshal(value, &orderOwner); err != nil {
			return err
		}

		order, found := k.GetPerpetualOrder(ctx, orderOwner.OrderId)
		if !found {
			return types.ErrPerpetualOrderNotFound
		}

		orders = append(orders, order)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.OwnerOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}
