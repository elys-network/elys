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

func (k Keeper) OwnerPerpetuals(goCtx context.Context, req *types.OwnerPerpetualsRequest) (*types.OwnerPerpetualsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	var prefixStore prefix.Store
	if req.SubAccountId == 0 {
		prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetPerpetualOwnerAddressKey(owner))
	} else {
		key := types.GetPerpetualOwnerAddressKey(owner)
		key = append(key, sdk.Uint64ToBigEndian(req.SubAccountId)...)
		key = append(key, []byte("/")...)
		prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)

	}

	var list []types.Perpetual

	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var perpetualOwner types.PerpetualOwner
		if err := k.cdc.Unmarshal(value, &perpetualOwner); err != nil {
			return err
		}

		p, err := k.GetPerpetual(ctx, perpetualOwner.MarketId, perpetualOwner.PerpetualId)
		if err != nil {
			return err
		}

		list = append(list, p)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.OwnerPerpetualsResponse{List: list, Pagination: pageRes}, nil
}
