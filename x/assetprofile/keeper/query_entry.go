package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/assetprofile/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) EntryAll(goCtx context.Context, req *types.QueryAllEntryRequest) (*types.QueryAllEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var entries []types.Entry
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	entryStore := prefix.NewStore(store, types.KeyPrefix(types.EntryKeyPrefix))

	pageRes, err := query.Paginate(entryStore, req.Pagination, func(key []byte, value []byte) error {
		var entry types.Entry
		if err := k.cdc.Unmarshal(value, &entry); err != nil {
			return err
		}

		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEntryResponse{Entry: entries, Pagination: pageRes}, nil
}

func (k Keeper) Entry(goCtx context.Context, req *types.QueryEntryRequest) (*types.QueryEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetEntry(ctx, req.BaseDenom)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryEntryResponse{Entry: val}, nil
}

func (k Keeper) EntryByDenom(goCtx context.Context, req *types.QueryEntryByDenomRequest) (*types.QueryEntryByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetEntryByDenom(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryEntryByDenomResponse{Entry: val}, nil
}
