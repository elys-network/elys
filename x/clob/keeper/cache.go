package keeper

import (
	"fmt"
	"sync"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// MarketCache provides caching for frequently accessed market data
type MarketCache struct {
	mu          sync.RWMutex
	markets     map[uint64]types.PerpetualMarket
	bestBids    map[uint64]math.LegacyDec
	bestAsks    map[uint64]math.LegacyDec
	blockHeight int64
}

// NewMarketCache creates a new market cache
func NewMarketCache() *MarketCache {
	return &MarketCache{
		markets:     make(map[uint64]types.PerpetualMarket),
		bestBids:    make(map[uint64]math.LegacyDec),
		bestAsks:    make(map[uint64]math.LegacyDec),
		blockHeight: 0,
	}
}

// GetMarket retrieves a market from cache or loads it
func (k *Keeper) GetCachedMarket(ctx sdk.Context, marketId uint64) (types.PerpetualMarket, error) {
	// Check if we have a cache and it's current
	if k.marketCache != nil {
		k.marketCache.mu.RLock()
		if k.marketCache.blockHeight == ctx.BlockHeight() {
			market, exists := k.marketCache.markets[marketId]
			k.marketCache.mu.RUnlock()
			if exists {
				return market, nil
			}
		} else {
			k.marketCache.mu.RUnlock()
			// Clear cache if block height changed
			k.ClearMarketCache()
		}
	}

	// Load from store
	market, err := k.GetPerpetualMarket(ctx, marketId)
	if err != nil {
		return types.PerpetualMarket{}, err
	}

	// Update cache
	if k.marketCache != nil {
		k.marketCache.mu.Lock()
		k.marketCache.markets[marketId] = market
		k.marketCache.blockHeight = ctx.BlockHeight()
		k.marketCache.mu.Unlock()
	}

	return market, nil
}

// GetCachedBestBid retrieves the best bid price from cache
func (k *Keeper) GetCachedBestBid(ctx sdk.Context, marketId uint64) math.LegacyDec {
	if k.marketCache != nil {
		k.marketCache.mu.RLock()
		if k.marketCache.blockHeight == ctx.BlockHeight() {
			price, exists := k.marketCache.bestBids[marketId]
			k.marketCache.mu.RUnlock()
			if exists {
				return price
			}
		} else {
			k.marketCache.mu.RUnlock()
		}
	}

	// Load from store
	price := k.GetHighestBuyPrice(ctx, marketId)

	// Update cache
	if k.marketCache != nil {
		k.marketCache.mu.Lock()
		k.marketCache.bestBids[marketId] = price
		k.marketCache.blockHeight = ctx.BlockHeight()
		k.marketCache.mu.Unlock()
	}

	return price
}

// GetCachedBestAsk retrieves the best ask price from cache
func (k *Keeper) GetCachedBestAsk(ctx sdk.Context, marketId uint64) math.LegacyDec {
	if k.marketCache != nil {
		k.marketCache.mu.RLock()
		if k.marketCache.blockHeight == ctx.BlockHeight() {
			price, exists := k.marketCache.bestAsks[marketId]
			k.marketCache.mu.RUnlock()
			if exists {
				return price
			}
		} else {
			k.marketCache.mu.RUnlock()
		}
	}

	// Load from store
	price := k.GetLowestSellPrice(ctx, marketId)

	// Update cache
	if k.marketCache != nil {
		k.marketCache.mu.Lock()
		k.marketCache.bestAsks[marketId] = price
		k.marketCache.blockHeight = ctx.BlockHeight()
		k.marketCache.mu.Unlock()
	}

	return price
}

// ClearMarketCache clears all cached market data
func (k *Keeper) ClearMarketCache() {
	if k.marketCache != nil {
		k.marketCache.mu.Lock()
		k.marketCache.markets = make(map[uint64]types.PerpetualMarket)
		k.marketCache.bestBids = make(map[uint64]math.LegacyDec)
		k.marketCache.bestAsks = make(map[uint64]math.LegacyDec)
		k.marketCache.mu.Unlock()
	}
}

// BatchGetMarkets retrieves multiple markets efficiently
func (k *Keeper) BatchGetMarkets(ctx sdk.Context, marketIds []uint64) (map[uint64]types.PerpetualMarket, error) {
	markets := make(map[uint64]types.PerpetualMarket)

	for _, id := range marketIds {
		market, err := k.GetCachedMarket(ctx, id)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to get market %d: %v", id, err))
			continue
		}
		markets[id] = market
	}

	return markets, nil
}

// PriceCache provides caching for oracle prices
type PriceCache struct {
	mu          sync.RWMutex
	prices      map[string]math.LegacyDec
	blockHeight int64
}

// NewPriceCache creates a new price cache
func NewPriceCache() *PriceCache {
	return &PriceCache{
		prices:      make(map[string]math.LegacyDec),
		blockHeight: 0,
	}
}

// GetCachedAssetPrice retrieves an asset price from cache or oracle
func (k *Keeper) GetCachedAssetPrice(ctx sdk.Context, denom string) (math.LegacyDec, error) {
	// Check if we have a cache and it's current
	if k.priceCache != nil {
		k.priceCache.mu.RLock()
		if k.priceCache.blockHeight == ctx.BlockHeight() {
			price, exists := k.priceCache.prices[denom]
			k.priceCache.mu.RUnlock()
			if exists {
				return price, nil
			}
		} else {
			k.priceCache.mu.RUnlock()
			// Clear cache if block height changed
			k.ClearPriceCache()
		}
	}

	// Load from oracle
	price, err := k.GetAssetPriceFromDenom(ctx, denom)
	if err != nil {
		return math.LegacyDec{}, err
	}

	// Update cache
	if k.priceCache != nil {
		k.priceCache.mu.Lock()
		k.priceCache.prices[denom] = price
		k.priceCache.blockHeight = ctx.BlockHeight()
		k.priceCache.mu.Unlock()
	}

	return price, nil
}

// BatchGetPrices retrieves multiple prices efficiently
func (k *Keeper) BatchGetPrices(ctx sdk.Context, denoms []string) (map[string]math.LegacyDec, error) {
	prices := make(map[string]math.LegacyDec)

	for _, denom := range denoms {
		price, err := k.GetCachedAssetPrice(ctx, denom)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to get price for %s: %v", denom, err))
			continue
		}
		prices[denom] = price
	}

	return prices, nil
}

// ClearPriceCache clears all cached price data
func (k *Keeper) ClearPriceCache() {
	if k.priceCache != nil {
		k.priceCache.mu.Lock()
		k.priceCache.prices = make(map[string]math.LegacyDec)
		k.priceCache.mu.Unlock()
	}
}
