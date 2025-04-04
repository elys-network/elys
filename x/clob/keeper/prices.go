package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetCurrentTwapPrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	reverseIterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

	var lastTwapPrice types.TwapPrice
	if reverseIterator.Valid() {
		k.cdc.MustUnmarshal(reverseIterator.Value(), &lastTwapPrice)
	}

	_ = reverseIterator.Close()

	iteratorForward := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer iteratorForward.Close()

	var firstTwapPrice types.TwapPrice
	if iteratorForward.Valid() {
		k.cdc.MustUnmarshal(iteratorForward.Value(), &firstTwapPrice)
	}

	num := lastTwapPrice.CumulativePrice.Sub(firstTwapPrice.CumulativePrice)
	if num.IsZero() {
		return math.LegacyZeroDec()
	}
	if lastTwapPrice.Timestamp <= firstTwapPrice.Timestamp {
		panic("twap price timestamp delta incorrect, time delta <= 0")
	}
	timeDelta := math.LegacyNewDec(int64(lastTwapPrice.Timestamp - firstTwapPrice.Timestamp))
	return num.Quo(timeDelta)
}

func (k Keeper) SetTwapPrices(ctx sdk.Context, currentTwapPrice types.TwapPrice) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(currentTwapPrice.MarketId)...))
	reverseIterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

	var lastTwapPrice types.TwapPrice
	if reverseIterator.Valid() {
		k.cdc.MustUnmarshal(reverseIterator.Value(), &lastTwapPrice)
	}
	_ = reverseIterator.Close()

	if lastTwapPrice.MarketId == 0 {
		currentTwapPrice.CumulativePrice = math.LegacyZeroDec()
	} else {
		// lastPrice×(now−lastUpdate)
		if currentTwapPrice.Timestamp <= lastTwapPrice.Timestamp {
			panic("twap price timestamp delta incorrect, time delta <= 0")
		}
		toAdd := lastTwapPrice.Price.Mul(math.LegacyNewDec(int64(currentTwapPrice.Timestamp - lastTwapPrice.Timestamp)))
		currentTwapPrice.CumulativePrice = lastTwapPrice.CumulativePrice.Add(toAdd)
	}
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetTwapPricesKey(currentTwapPrice.MarketId, currentTwapPrice.Block)
	b := k.cdc.MustMarshal(&currentTwapPrice)
	store.Set(key, b)

	// Setting in transient store for fast access
	tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
	tStore.Set(key, b)

	iteratorForward := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer iteratorForward.Close()

	market, err := k.GetPerpetualMarket(ctx, currentTwapPrice.MarketId)
	if err != nil {
		panic(err)
	}

	for ; iteratorForward.Valid(); iteratorForward.Next() {
		var old types.TwapPrice
		k.cdc.MustUnmarshal(iteratorForward.Value(), &old)
		// usually currentTwapPrice.Timestamp will be equal to ctx.BlockTime.Unix()
		// While init genesis, this will not delete all old twap prices
		if old.Timestamp < currentTwapPrice.Timestamp-market.MaxTwapPricesTime {
			prefixStore.Delete(iteratorForward.Key())
		} else {
			break // store is block-ordered, so stop early
		}
	}
}

func (k Keeper) GetAllTwapPrices(ctx sdk.Context) []types.TwapPrice {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.TwapPricesPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.TwapPrice

	for ; iterator.Valid(); iterator.Next() {
		var val types.TwapPrice
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// GetLastMarketPrice First it checks transient store and then it checks KVstore
func (k Keeper) GetLastMarketPrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	prefixTStore := prefix.NewStore(runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	tStoreIterator := storetypes.KVStoreReversePrefixIterator(prefixTStore, []byte{})

	defer tStoreIterator.Close()

	var lastTwapPrice types.TwapPrice
	lastTwapPrice.Price = math.LegacyZeroDec()
	if tStoreIterator.Valid() {
		k.cdc.MustUnmarshal(tStoreIterator.Value(), &lastTwapPrice)
	}

	if lastTwapPrice.MarketId == 0 {
		prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
		iterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

		defer iterator.Close()

		if iterator.Valid() {
			k.cdc.MustUnmarshal(iterator.Value(), &lastTwapPrice)
		}

		// Setting in transient store for fast access
		tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
		key := types.GetTwapPricesKey(marketId, lastTwapPrice.Block)
		b := k.cdc.MustMarshal(&lastTwapPrice)
		tStore.Set(key, b)
	}

	return lastTwapPrice.Price
}

func (k Keeper) GetHighestBuyPrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	iterator := k.GetBuyOrderIterator(ctx, marketId)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val.Price
	}
	return math.LegacyZeroDec()
}

func (k Keeper) GetLowestSellPrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	iterator := k.GetSellOrderIterator(ctx, marketId)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val.Price
	}
	return math.LegacyZeroDec()
}

func (k Keeper) GetMidPrice(ctx sdk.Context, marketId uint64) (math.LegacyDec, error) {
	highestBuy := k.GetLowestSellPrice(ctx, marketId)
	lowestSell := k.GetLowestSellPrice(ctx, marketId)
	if highestBuy.IsZero() || lowestSell.IsZero() {
		return math.LegacyZeroDec(), errors.New("one side of the orderbook is empty")
	}
	sumPrice := highestBuy.Add(lowestSell)
	return sumPrice.Quo(math.LegacyNewDec(2)), nil
}
