package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/assetprofile/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) EntryAll(goCtx context.Context, req *types.QueryAllEntryRequest) (*types.QueryAllEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var entrys []types.Entry
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	entryStore := prefix.NewStore(store, types.KeyPrefix(types.EntryKeyPrefix))

	pageRes, err := query.Paginate(entryStore, req.Pagination, func(key []byte, value []byte) error {
		var entry types.Entry
		if err := k.cdc.Unmarshal(value, &entry); err != nil {
			return err
		}

		entrys = append(entrys, entry)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEntryResponse{Entry: entrys, Pagination: pageRes}, nil
}

func (k Keeper) Entry(goCtx context.Context, req *types.QueryGetEntryRequest) (*types.QueryGetEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetEntry(ctx, req.BaseDenom)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetEntryResponse{Entry: val}, nil
}
