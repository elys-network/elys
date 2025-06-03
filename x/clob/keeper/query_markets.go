package keeper

import (
	"context"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetMarketInfo(ctx sdk.Context, market types.PerpetualMarket) types.MarketResponse {
	currentTwapPrice := k.GetCurrentTwapPrice(ctx, market.Id)
	lastAverageTradePrice := k.GetLastAverageTradePrice(ctx, market.Id)
	highestBuyPrice := k.GetHighestBuyPrice(ctx, market.Id)
	lowestSellPrice := k.GetLowestSellPrice(ctx, market.Id)
	midPrice, err := k.GetMidPrice(ctx, market.Id)
	if err != nil {
		midPrice = math.LegacyZeroDec()
	}

	return types.MarketResponse{
		Market:                market,
		CurrentTwapPrice:      currentTwapPrice,
		LastAverageTradePrice: lastAverageTradePrice,
		HighestBuyPrice:       highestBuyPrice,
		LowestSellPrice:       lowestSellPrice,
		MidPrice:              midPrice,
	}
}

func (k Keeper) AllMarkets(goCtx context.Context, req *types.AllMarketsRequest) (*types.AllMarketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var responses []types.MarketResponse

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.PerpetualMarketPrefix)

	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var market types.PerpetualMarket
		if err := k.cdc.Unmarshal(value, &market); err != nil {
			return err
		}

		marketInfo := k.GetMarketInfo(ctx, market)
		responses = append(responses, marketInfo)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.AllMarketsResponse{Markets: responses, Pagination: pageRes}, nil
}

func (k Keeper) Market(goCtx context.Context, req *types.MarketRequest) (*types.MarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := k.GetPerpetualMarket(ctx, req.MarketId)
	if err != nil {
		return nil, err
	}

	marketInfo := k.GetMarketInfo(ctx, market)

	return &marketInfo, nil
}
