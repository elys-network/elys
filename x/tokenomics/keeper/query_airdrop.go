package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AirdropAll(goCtx context.Context, req *types.QueryAllAirdropRequest) (*types.QueryAllAirdropResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var airdrops []types.Airdrop
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	airdropStore := prefix.NewStore(store, types.KeyPrefix(types.AirdropKeyPrefix))

	pageRes, err := query.Paginate(airdropStore, req.Pagination, func(key []byte, value []byte) error {
		var airdrop types.Airdrop
		if err := k.cdc.Unmarshal(value, &airdrop); err != nil {
			return err
		}

		airdrops = append(airdrops, airdrop)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAirdropResponse{Airdrop: airdrops, Pagination: pageRes}, nil
}

func (k Keeper) Airdrop(goCtx context.Context, req *types.QueryGetAirdropRequest) (*types.QueryGetAirdropResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetAirdrop(ctx, req.Intent)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAirdropResponse{Airdrop: val}, nil
}
