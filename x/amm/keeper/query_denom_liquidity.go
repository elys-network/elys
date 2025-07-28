package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DenomLiquidityAll(goCtx context.Context, req *types.QueryAllDenomLiquidityRequest) (*types.QueryAllDenomLiquidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var denomLiquiditys []types.DenomLiquidity
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	denomLiquidityStore := prefix.NewStore(store, types.KeyPrefix(types.DenomLiquidityKeyPrefix))

	pageRes, err := query.Paginate(denomLiquidityStore, req.Pagination, func(key []byte, value []byte) error {
		var denomLiquidity types.DenomLiquidity
		if err := k.cdc.Unmarshal(value, &denomLiquidity); err != nil {
			return err
		}

		denomLiquiditys = append(denomLiquiditys, denomLiquidity)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDenomLiquidityResponse{DenomLiquidity: denomLiquiditys, Pagination: pageRes}, nil
}

func (k Keeper) DenomLiquidity(goCtx context.Context, req *types.QueryGetDenomLiquidityRequest) (*types.QueryGetDenomLiquidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetDenomLiquidity(
		ctx,
		req.Denom,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDenomLiquidityResponse{DenomLiquidity: val}, nil
}
