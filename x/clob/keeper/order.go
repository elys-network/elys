package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) GetPerpetualOrder(ctx sdk.Context, orderId types.OrderId) (types.PerpetualOrder, bool) {
	key := types.GetPerpetualOrderKey(orderId.MarketId, orderId.OrderType, orderId.PriceTick, orderId.Counter)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.PerpetualOrder{}, false
	}

	var val types.PerpetualOrder
	if err := k.cdc.Unmarshal(b, &val); err != nil {
		ctx.Logger().Error("failed to unmarshal perpetual order", "key", orderId, "error", err)
		return types.PerpetualOrder{}, false
	}
	return val, true
}

func (k Keeper) GetAllPerpetualOrders(ctx sdk.Context) []types.PerpetualOrder {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOrderPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualOrder

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetualOrder(ctx sdk.Context, v types.PerpetualOrder) error {
	if v.GetOrderType() == types.OrderType_ORDER_TYPE_LIMIT_SELL || v.GetOrderType() == types.OrderType_ORDER_TYPE_LIMIT_BUY || v.GetPriceTick() <= 0 {
		store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		key := types.GetPerpetualOrderKey(v.GetMarketId(), v.GetOrderType(), v.GetPriceTick(), v.GetCounter())
		b := k.cdc.MustMarshal(&v)
		store.Set(key, b)
		return nil
	} else {
		return types.ErrInvalidOrderType
	}
}

func (k Keeper) DeleteOrder(ctx sdk.Context, perpetualOrderOwner types.PerpetualOrderOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOrderKey(perpetualOrderOwner.OrderId.MarketId, perpetualOrderOwner.OrderId.OrderType, perpetualOrderOwner.OrderId.PriceTick, perpetualOrderOwner.OrderId.Counter)
	store.Delete(key)

	k.DeleteOrderOwner(ctx, perpetualOrderOwner)
}

func (k Keeper) GetBuyOrderIterator(ctx sdk.Context, marketId uint64) storetypes.Iterator {
	key := types.GetPerpetualOrderBookIteratorKey(marketId, true)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	return storetypes.KVStoreReversePrefixIterator(store, []byte{})

}

func (k Keeper) GetSellOrderIterator(ctx sdk.Context, marketId uint64) storetypes.Iterator {
	key := types.GetPerpetualOrderBookIteratorKey(marketId, false)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	return storetypes.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) RequiredBalanceForOrder(ctx sdk.Context, order types.PerpetualOrder) (sdk.Coin, error) {
	market, err := k.GetPerpetualMarket(ctx, order.GetMarketId())
	if err != nil {
		return sdk.Coin{}, err
	}

	maxFees := k.CalculateMaxFees(market, order)
	initialMarginRequired := order.UnfilledValue().Mul(market.InitialMarginRatio).RoundInt()

	amount := maxFees.Add(initialMarginRequired)

	return sdk.NewCoin(market.QuoteDenom, amount), nil
}
