package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CurrentTwapPrice(goCtx context.Context, req *types.CurrentTwapPriceRequest) (*types.CurrentTwapPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.CurrentTwapPriceResponse{Price: k.GetCurrentTwapPrice(ctx, req.MarketId)}, nil
}

func (k Keeper) AllTwapPrices(goCtx context.Context, req *types.AllTwapPricesRequest) (*types.AllTwapPricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.TwapPricesPrefix)

	var twapPrices []types.TwapPrice
	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var twapPrice types.TwapPrice
		if err := k.cdc.Unmarshal(value, &twapPrice); err != nil {
			return err
		}

		twapPrices = append(twapPrices, twapPrice)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.AllTwapPricesResponse{TwapPrices: twapPrices, Pagination: pageRes}, nil
}

func (k Keeper) LastAverageTradePrice(goCtx context.Context, req *types.LastAverageTradePriceRequest) (*types.LastAverageTradePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.LastAverageTradePriceResponse{Price: k.GetLastAverageTradePrice(ctx, req.MarketId)}, nil
}

func (k Keeper) HighestBuyPrice(goCtx context.Context, req *types.HighestBuyPriceRequest) (*types.HighestBuyPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.HighestBuyPriceResponse{Price: k.GetHighestBuyPrice(ctx, req.MarketId)}, nil
}

func (k Keeper) LowestSellPrice(goCtx context.Context, req *types.LowestSellPriceRequest) (*types.LowestSellPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.LowestSellPriceResponse{Price: k.GetLowestSellPrice(ctx, req.MarketId)}, nil
}

func (k Keeper) MidPrice(goCtx context.Context, req *types.MidPriceRequest) (*types.MidPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	midPrice, err := k.GetMidPrice(ctx, req.MarketId)
	if err != nil {
		return nil, err
	}

	return &types.MidPriceResponse{Price: midPrice}, nil
}
