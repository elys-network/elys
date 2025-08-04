package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) GetCurrentTwapPrice(ctx sdk.Context, marketId uint64) (math.LegacyDec, error) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	reverseIterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})

	var lastTwapPrice types.TwapPrice
	if reverseIterator.Valid() {
		if err := k.cdc.Unmarshal(reverseIterator.Value(), &lastTwapPrice); err != nil {
			ctx.Logger().Error("failed to unmarshal last twap price", "marketId", marketId, "error", err)
			return math.LegacyZeroDec(), err
		}
	} else {
		lastTwapPrice = types.TwapPrice{
			MarketId:          marketId,
			Block:             uint64(ctx.BlockHeight()),
			AverageTradePrice: math.LegacyZeroDec(),
			TotalVolume:       math.LegacyZeroDec(),
			CumulativePrice:   math.LegacyZeroDec(),
			Timestamp:         uint64(ctx.BlockTime().Unix()),
		}
	}

	_ = reverseIterator.Close()

	iteratorForward := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
	defer iteratorForward.Close()

	var firstTwapPrice types.TwapPrice
	if iteratorForward.Valid() {
		k.cdc.MustUnmarshal(iteratorForward.Value(), &firstTwapPrice)
	} else {
		firstTwapPrice = types.TwapPrice{
			MarketId:          marketId,
			Block:             uint64(ctx.BlockHeight()),
			AverageTradePrice: math.LegacyZeroDec(),
			TotalVolume:       math.LegacyZeroDec(),
			CumulativePrice:   math.LegacyZeroDec(),
			Timestamp:         uint64(ctx.BlockTime().Unix()),
		}
	}

	if lastTwapPrice.Timestamp < firstTwapPrice.Timestamp {
		return math.LegacyZeroDec(), errors.New("twap price timestamp delta incorrect, time delta < 0")
	}
	// Handles the case when no or only 1 twap price is present
	// Returning average price, because to calculate unrealized PnL, mark price should not be 0 otherwise losses shown will be too high
	if lastTwapPrice.Timestamp == firstTwapPrice.Timestamp {
		return lastTwapPrice.AverageTradePrice, nil
	}

	num := lastTwapPrice.CumulativePrice.Sub(firstTwapPrice.CumulativePrice)
	if num.IsZero() {
		return math.LegacyZeroDec(), nil
	}
	timeDelta := math.LegacyNewDec(int64(lastTwapPrice.Timestamp - firstTwapPrice.Timestamp))
	return num.Quo(timeDelta), nil
}

// SetTwapPricesStruct Should only be called by Import or Init Genesis
func (k Keeper) SetTwapPricesStruct(ctx sdk.Context, twapPrice types.TwapPrice) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetTwapPricesKey(twapPrice.MarketId, twapPrice.Block)
	store.Set(key, k.cdc.MustMarshal(&twapPrice))
}

func (k Keeper) SetTwapPrices(ctx sdk.Context, trade types.Trade) error {
	if trade.Quantity.LTE(math.LegacyZeroDec()) {
		return errors.New("trade quantity cannot be negative or zero")
	}
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetTwapPricesKey(trade.MarketId, uint64(ctx.BlockHeight()))

	currentTwapPrice := types.TwapPrice{
		MarketId:          trade.MarketId,
		Block:             uint64(ctx.BlockHeight()),
		AverageTradePrice: math.LegacyZeroDec(),
		TotalVolume:       math.LegacyZeroDec(),
		CumulativePrice:   math.LegacyZeroDec(),
		Timestamp:         uint64(ctx.BlockTime().Unix()),
	}

	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &currentTwapPrice)

		oldTradeValue := currentTwapPrice.AverageTradePrice.Mul(currentTwapPrice.TotalVolume)
		newAveragePrice := oldTradeValue.Add(trade.GetTradeValue()).Quo(currentTwapPrice.TotalVolume.Add(trade.Quantity))

		currentTwapPrice.AverageTradePrice = newAveragePrice
		currentTwapPrice.TotalVolume = currentTwapPrice.TotalVolume.Add(trade.Quantity.Abs())
		// currentTwapPrice will not change as it is previous trade price multiplied amount of time that price stayed

		b := k.cdc.MustMarshal(&currentTwapPrice)
		store.Set(key, b)

		// Setting in transient store for fast access
		tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
		tStore.Set(key, b)
	} else {
		currentTwapPrice.AverageTradePrice = trade.Price
		currentTwapPrice.TotalVolume = trade.Quantity

		prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(trade.MarketId)...))
		reverseIterator := storetypes.KVStoreReversePrefixIterator(prefixStore, []byte{})
		defer reverseIterator.Close()

		var lastTwapPrice types.TwapPrice
		if reverseIterator.Valid() {
			k.cdc.MustUnmarshal(reverseIterator.Value(), &lastTwapPrice)
		}

		if lastTwapPrice.MarketId == 0 {
			currentTwapPrice.CumulativePrice = math.LegacyZeroDec()
		} else {
			// lastPrice×(now−lastUpdate)
			if currentTwapPrice.Timestamp <= lastTwapPrice.Timestamp {
				return errors.New("twap price timestamp delta incorrect, time delta <= 0")
			}
			toAdd := lastTwapPrice.AverageTradePrice.Mul(math.LegacyNewDec(int64(currentTwapPrice.Timestamp - lastTwapPrice.Timestamp)))
			currentTwapPrice.CumulativePrice = lastTwapPrice.CumulativePrice.Add(toAdd)
		}
		b := k.cdc.MustMarshal(&currentTwapPrice)
		store.Set(key, b)

		// Setting in transient store for fast access
		tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
		tStore.Set(key, b)

		iteratorForward := storetypes.KVStorePrefixIterator(prefixStore, []byte{})
		defer iteratorForward.Close()

		market, err := k.GetPerpetualMarket(ctx, currentTwapPrice.MarketId)
		if err != nil {
			return err
		}

		for ; iteratorForward.Valid(); iteratorForward.Next() {
			var old types.TwapPrice
			k.cdc.MustUnmarshal(iteratorForward.Value(), &old)
			// usually currentTwapPrice.Timestamp will be equal to ctx.BlockTime.Unix()
			// While init genesis, this will not delete all old twap prices
			if old.Timestamp < currentTwapPrice.Timestamp-market.TwapPricesWindow {
				prefixStore.Delete(iteratorForward.Key())
			} else {
				break // store is block-ordered, so stop early
			}
		}
	}

	return nil
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

// GetLastAverageTradePrice First it checks transient store and then it checks KVstore
func (k Keeper) GetLastAverageTradePrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	prefixTStore := prefix.NewStore(runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx)), append(types.TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...))
	tStoreIterator := storetypes.KVStoreReversePrefixIterator(prefixTStore, []byte{})

	defer tStoreIterator.Close()

	var lastTwapPrice types.TwapPrice
	lastTwapPrice.AverageTradePrice = math.LegacyZeroDec()
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

	return lastTwapPrice.AverageTradePrice
}

func (k Keeper) GetHighestBuyPrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	iterator := k.GetBuyOrderIterator(ctx, marketId)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Order
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val.GetPrice()
	}
	return math.LegacyZeroDec()
}

func (k Keeper) GetLowestSellPrice(ctx sdk.Context, marketId uint64) math.LegacyDec {
	iterator := k.GetSellOrderIterator(ctx, marketId)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Order
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val.GetPrice()
	}
	return math.LegacyZeroDec()
}

func (k Keeper) GetMidPrice(ctx sdk.Context, marketId uint64) (math.LegacyDec, error) {
	highestBuy := k.GetHighestBuyPrice(ctx, marketId)
	lowestSell := k.GetLowestSellPrice(ctx, marketId)
	if highestBuy.IsZero() || lowestSell.IsZero() {
		return math.LegacyZeroDec(), errors.New("one side of the orderbook is empty")
	}
	sumPrice := highestBuy.Add(lowestSell)
	return sumPrice.Quo(math.LegacyNewDec(2)), nil
}
