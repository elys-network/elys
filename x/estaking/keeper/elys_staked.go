package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

// SetElysStaked set a specific elysStaked in the store from its index
func (k Keeper) SetElysStaked(ctx sdk.Context, elysStaked types.ElysStaked) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&elysStaked)
	key := types.GetElysStakedKey(elysStaked.GetAccountAddress())
	store.Set(key, b)
}

// GetElysStaked returns a elysStaked from its index
func (k Keeper) GetElysStaked(ctx sdk.Context, address sdk.AccAddress) types.ElysStaked {
	key := types.GetElysStakedKey(address)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return types.ElysStaked{
			Address: address.String(),
			Amount:  sdk.ZeroInt(),
		}
	}
	var val types.ElysStaked
	k.cdc.MustUnmarshal(bz, &val)
	return val
}

// RemoveElysStaked removes a elysStaked from the store
func (k Keeper) RemoveElysStaked(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetElysStakedKey(acc))
}

// GetAllElysStaked returns all elysStaked
func (k Keeper) GetAllElysStaked(ctx sdk.Context) (list []types.ElysStaked) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ElysStakedKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ElysStaked
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// remove after migration
func (k Keeper) GetAllLegacyElysStaked(ctx sdk.Context) (list []types.ElysStaked) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyElysStakedKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ElysStaked
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// TODO: remove all legacy prefixes and functions after migration
func (k Keeper) DeleteLegacyElysStaked(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyElysStakedKeyPrefix))
	store.Delete(types.LegacyElysStakedKey(address))
}
