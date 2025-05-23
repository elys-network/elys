package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v5/x/oracle/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AssetInfoAll(c context.Context, req *types.QueryAllAssetInfoRequest) (*types.QueryAllAssetInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var assetInfos []types.AssetInfo
	ctx := sdk.UnwrapSDKContext(c)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	assetInfoStore := prefix.NewStore(store, types.KeyPrefix(types.AssetInfoKeyPrefix))

	pageRes, err := query.Paginate(assetInfoStore, req.Pagination, func(key []byte, value []byte) error {
		var assetInfo types.AssetInfo
		if err := k.cdc.Unmarshal(value, &assetInfo); err != nil {
			return err
		}

		assetInfos = append(assetInfos, assetInfo)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAssetInfoResponse{AssetInfo: assetInfos, Pagination: pageRes}, nil
}

func (k Keeper) AssetInfo(c context.Context, req *types.QueryGetAssetInfoRequest) (*types.QueryGetAssetInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetAssetInfo(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAssetInfoResponse{AssetInfo: val}, nil
}
