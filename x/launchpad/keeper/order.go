package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

func (k Keeper) GetOrder(ctx sdk.Context, orderId uint64) types.Purchase {
	order := types.Purchase{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PurchasePrefixKey)
	bz := store.Get(sdk.Uint64ToBigEndian(orderId))
	if len(bz) == 0 {
		return types.Purchase{
			OrderId:            0,
			OrderMaker:         "",
			SpendingToken:      "",
			TokenAmount:        sdk.ZeroInt(),
			ElysAmount:         sdk.ZeroInt(),
			ReturnedElysAmount: sdk.ZeroInt(),
			BonusAmount:        sdk.ZeroInt(),
		}
	}

	k.cdc.MustUnmarshal(bz, &order)
	return order
}

func (k Keeper) SetOrder(ctx sdk.Context, order types.Purchase) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PurchasePrefixKey)
	bz := k.cdc.MustMarshal(&order)
	store.Set(sdk.Uint64ToBigEndian(order.OrderId), bz)
}

func (k Keeper) DeleteOrder(ctx sdk.Context, order types.Purchase) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PurchasePrefixKey)
	store.Delete(sdk.Uint64ToBigEndian(order.OrderId))
}

func (k Keeper) GetAllOrders(ctx sdk.Context) []types.Purchase {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PurchasePrefixKey)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	orders := []types.Purchase{}
	for ; iterator.Valid(); iterator.Next() {
		order := types.Purchase{}
		k.cdc.MustUnmarshal(iterator.Value(), &order)

		orders = append(orders, order)
	}
	return orders
}

func (k Keeper) LastOrderId(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PurchasePrefixKey)

	iterator := sdk.KVStoreReversePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		order := types.Purchase{}
		k.cdc.MustUnmarshal(iterator.Value(), &order)

		return order.OrderId
	}
	return 0
}
