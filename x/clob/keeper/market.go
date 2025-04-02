package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetualMarket(ctx sdk.Context, id uint64) (types.PerpetualMarket, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketKey(id)
	b := store.Get(key)
	if b == nil {
		return types.PerpetualMarket{}, types.ErrPerpetualMarketNotFound
	}

	var v types.PerpetualMarket
	k.cdc.MustUnmarshal(b, &v)
	return v, nil
}

func (k Keeper) SetPerpetualMarket(ctx sdk.Context, p types.PerpetualMarket) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketKey(p.Id)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}

func (k Keeper) GetAllPerpetualMarket(ctx sdk.Context) []types.PerpetualMarket {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualMarket

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualMarket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) CountAllPerpetualMarket(ctx sdk.Context) uint64 {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	count := uint64(0)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	return count
}

func (k Keeper) CheckPerpetualMarketAlreadyExists(ctx sdk.Context, baseDenom, quoteDenom string) bool {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualMarket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.QuoteDenom == quoteDenom && val.BaseDenom == baseDenom {
			return true
		}
	}

	return false
}

func (k Keeper) GetTwapPrice(ctx sdk.Context, marketId uint64) (math.Dec, error) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	iterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

	var lastTwapPrice types.TwapPrice
	defer iterator.Close()
	for iterator.Valid() {
		k.cdc.MustUnmarshal(iterator.Value(), &lastTwapPrice)
	}

	prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	forwardIterator := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer forwardIterator.Close()

	var firstTwapPrice types.TwapPrice
	for iterator.Valid() {
		k.cdc.MustUnmarshal(iterator.Value(), &firstTwapPrice)
	}

	num, err := lastTwapPrice.CumulativePrice.Sub(firstTwapPrice.CumulativePrice)
	if err != nil {
		return utils.ZeroDec, err
	}
	timeDelta := math.NewDecFromInt64(int64(lastTwapPrice.Timestamp - firstTwapPrice.Timestamp))
	return num.Quo(timeDelta)
}

func (k Keeper) SetTwapPrices(ctx sdk.Context, p types.TwapPrice) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(p.MarketId)...))
	iterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

	var lastTwapPrice types.TwapPrice
	for iterator.Valid() {
		k.cdc.MustUnmarshal(iterator.Value(), &lastTwapPrice)
	}
	_ = iterator.Close()

	if lastTwapPrice.MarketId == 0 {
		p.CumulativePrice = utils.ZeroDec
	} else {
		// lastPrice×(now−lastUpdate)
		toAdd, err := lastTwapPrice.Price.Mul(math.NewDecFromInt64(int64(p.Timestamp - lastTwapPrice.Timestamp)))
		if err != nil {
			panic(err)
		}
		p.CumulativePrice, err = lastTwapPrice.CumulativePrice.Add(toAdd)
		if err != nil {
			panic(err)
		}
	}
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetTwapPricesKey(p.MarketId, p.Block)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)

	// Setting in transient store for fast access
	tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
	tStore.Set(key, b)

	prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(p.MarketId)...))
	iteratorForward := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer iteratorForward.Close()

	for ; iteratorForward.Valid(); iteratorForward.Next() {
		var old types.TwapPrice
		k.cdc.MustUnmarshal(iteratorForward.Value(), &old)
		if old.Block < p.Block-1000 {
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
func (k Keeper) GetLastMarketPrice(ctx sdk.Context, marketId uint64) math.Dec {
	prefixTStore := prefix.NewStore(runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	tStoreIterator := storetypes.KVStoreReversePrefixIterator(prefixTStore, []byte{})

	defer tStoreIterator.Close()

	var lastTwapPrice types.TwapPrice
	lastTwapPrice.Price = utils.ZeroDec
	for tStoreIterator.Valid() {
		k.cdc.MustUnmarshal(tStoreIterator.Value(), &lastTwapPrice)
	}

	if lastTwapPrice.MarketId == 0 {
		prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
		iterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

		defer iterator.Close()

		for iterator.Valid() {
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

func (k Keeper) GetCurrentTwapPrice(ctx sdk.Context, marketId uint64) (math.Dec, error) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	reverseIterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

	var lastTwapPrice types.TwapPrice
	for reverseIterator.Valid() {
		k.cdc.MustUnmarshal(reverseIterator.Value(), &lastTwapPrice)
	}
	defer reverseIterator.Close()

	prefixStore = prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	iteratorForward := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer iteratorForward.Close()

	var firstTwapPrice types.TwapPrice
	for iteratorForward.Valid() {
		k.cdc.MustUnmarshal(reverseIterator.Value(), &firstTwapPrice)
	}
	num, err := lastTwapPrice.CumulativePrice.Sub(firstTwapPrice.CumulativePrice)
	if err != nil {
		return math.Dec{}, err
	}
	den := math.NewDecFromInt64(int64(lastTwapPrice.Timestamp - firstTwapPrice.Timestamp))
	return num.Quo(den)
}
