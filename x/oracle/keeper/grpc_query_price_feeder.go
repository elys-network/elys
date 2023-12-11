package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/oracle/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PriceFeederAll(c context.Context, req *types.QueryAllPriceFeederRequest) (*types.QueryAllPriceFeederResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var priceFeeders []types.PriceFeeder
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	priceFeederStore := prefix.NewStore(store, types.KeyPrefix(types.PriceFeederKeyPrefix))

	pageRes, err := query.Paginate(priceFeederStore, req.Pagination, func(key []byte, value []byte) error {
		var priceFeeder types.PriceFeeder
		if err := k.cdc.Unmarshal(value, &priceFeeder); err != nil {
			return err
		}

		priceFeeders = append(priceFeeders, priceFeeder)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPriceFeederResponse{PriceFeeder: priceFeeders, Pagination: pageRes}, nil
}

func (k Keeper) PriceFeeder(c context.Context, req *types.QueryGetPriceFeederRequest) (*types.QueryGetPriceFeederResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetPriceFeeder(ctx, req.Feeder)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPriceFeederResponse{PriceFeeder: val}, nil
}
