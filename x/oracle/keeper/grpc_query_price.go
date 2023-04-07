package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/oracle/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PriceAll(c context.Context, req *types.QueryAllPriceRequest) (*types.QueryAllPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var prices []types.Price
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	priceStore := prefix.NewStore(store, types.KeyPrefix(types.PriceKeyPrefix))

	pageRes, err := query.Paginate(priceStore, req.Pagination, func(key []byte, value []byte) error {
		var price types.Price
		if err := k.cdc.Unmarshal(value, &price); err != nil {
			return err
		}

		prices = append(prices, price)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPriceResponse{Price: prices, Pagination: pageRes}, nil
}

func (k Keeper) Price(c context.Context, req *types.QueryGetPriceRequest) (*types.QueryGetPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	fmt.Println("Price1")
	// if both source and timestamp are defined, use specific value
	if req.Source != "" && req.Timestamp != 0 {
		val, found := k.GetPrice(ctx, req.Asset, req.Source, req.Timestamp)
		if !found {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return &types.QueryGetPriceResponse{Price: val}, nil
	}

	fmt.Println("Price2")
	// if source is specified use latest price from source
	if req.Source != "" {
		val, found := k.GetLatestPriceFromAssetAndSource(ctx, req.Asset, req.Source)
		if !found {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return &types.QueryGetPriceResponse{Price: val}, nil
	}

	fmt.Println("Price3")
	// try out band source
	val, found := k.GetLatestPriceFromAssetAndSource(ctx, req.Asset, types.BAND)
	if found {
		return &types.QueryGetPriceResponse{Price: val}, nil
	}

	fmt.Println("Price4")
	// try out binance source
	val, found = k.GetLatestPriceFromAssetAndSource(ctx, req.Asset, types.BINANCE)
	if found {
		return &types.QueryGetPriceResponse{Price: val}, nil
	}

	fmt.Println("Price5")
	// try out osmosis source
	val, found = k.GetLatestPriceFromAssetAndSource(ctx, req.Asset, types.OSMOSIS)
	if found {
		return &types.QueryGetPriceResponse{Price: val}, nil
	}

	fmt.Println("Price6")
	// find from any source if band source does not exist
	val, found = k.GetLatestPriceFromAnySource(ctx, req.Asset)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return &types.QueryGetPriceResponse{Price: val}, nil
}
