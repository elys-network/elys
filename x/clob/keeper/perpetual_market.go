package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

func (k Keeper) GetPerpetualMarket(ctx sdk.Context, id uint64) (types.PerpetualMarket, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketKey(id)
	b := store.Get(key)
	if b == nil {
		return types.PerpetualMarket{}, types.ErrPerpetualMarketNotFound
	}

	var v types.PerpetualMarket
	if err := k.cdc.Unmarshal(b, &v); err != nil {
		return types.PerpetualMarket{}, errors.Wrapf(err, "failed to unmarshal perpetual market for id: %d", id)
	}
	return v, nil
}

func (k Keeper) SetPerpetualMarket(ctx sdk.Context, p types.PerpetualMarket) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketKey(p.Id)
	b, err := k.cdc.Marshal(&p)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal perpetual market for id: %d", p.Id)
	}
	store.Set(key, b)
	return nil
}

func (k Keeper) GetAllPerpetualMarket(ctx sdk.Context) []types.PerpetualMarket {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualMarket

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualMarket
		if err := k.cdc.Unmarshal(iterator.Value(), &val); err != nil {
			ctx.Logger().Error("failed to unmarshal perpetual market", "error", err)
			continue
		}
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
		if err := k.cdc.Unmarshal(iterator.Value(), &val); err != nil {
			ctx.Logger().Error("failed to unmarshal perpetual market in CheckPerpetualMarketAlreadyExists", "error", err)
			continue
		}
		if val.QuoteDenom == quoteDenom && val.BaseDenom == baseDenom {
			return true
		}
	}

	return false
}
