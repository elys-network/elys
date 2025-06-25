package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/burner/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) HistoryAll(goCtx context.Context, req *types.QueryAllHistoryRequest) (*types.QueryAllHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var histories []types.History
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	historyStore := prefix.NewStore(store, types.HistoryKeyPrefix)

	pageRes, err := query.Paginate(historyStore, req.Pagination, func(key []byte, value []byte) error {
		var history types.History
		if err := k.cdc.Unmarshal(value, &history); err != nil {
			return err
		}

		histories = append(histories, history)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllHistoryResponse{History: histories, Pagination: pageRes}, nil
}

func (k Keeper) History(goCtx context.Context, req *types.QueryGetHistoryRequest) (*types.QueryGetHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetHistory(
		ctx,
		req.Block,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetHistoryResponse{History: val}, nil
}
