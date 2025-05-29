package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/accountedpool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountedPoolAll(goCtx context.Context, req *types.QueryAllAccountedPoolRequest) (*types.QueryAllAccountedPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var accountedPools []types.AccountedPool
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	accountedPoolStore := prefix.NewStore(store, types.KeyPrefix(types.AccountedPoolKeyPrefix))

	pageRes, err := query.Paginate(accountedPoolStore, req.Pagination, func(key []byte, value []byte) error {
		var accountedPool types.AccountedPool
		if err := k.cdc.Unmarshal(value, &accountedPool); err != nil {
			return err
		}

		accountedPools = append(accountedPools, accountedPool)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAccountedPoolResponse{AccountedPool: accountedPools, Pagination: pageRes}, nil
}

func (k Keeper) AccountedPool(goCtx context.Context, req *types.QueryGetAccountedPoolRequest) (*types.QueryGetAccountedPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetAccountedPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAccountedPoolResponse{AccountedPool: val}, nil
}
