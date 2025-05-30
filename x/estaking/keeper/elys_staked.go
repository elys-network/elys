package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/estaking/types"
)

// SetElysStaked set a specific elysStaked in the store from its index
func (k Keeper) SetElysStaked(ctx sdk.Context, elysStaked types.ElysStaked) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&elysStaked)
	key := types.GetElysStakedKey(elysStaked.GetAccountAddress())
	store.Set(key, b)
}

// GetElysStaked returns a elysStaked from its index
func (k Keeper) GetElysStaked(ctx sdk.Context, address sdk.AccAddress) types.ElysStaked {
	key := types.GetElysStakedKey(address)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(key)
	if bz == nil {
		return types.ElysStaked{
			Address: address.String(),
			Amount:  sdkmath.ZeroInt(),
		}
	}
	var val types.ElysStaked
	k.cdc.MustUnmarshal(bz, &val)
	return val
}

// RemoveElysStaked removes a elysStaked from the store
func (k Keeper) RemoveElysStaked(ctx sdk.Context, acc sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetElysStakedKey(acc))
}

// GetAllElysStaked returns all elysStaked
func (k Keeper) GetAllElysStaked(ctx sdk.Context) (list []types.ElysStaked) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.ElysStakedKeyPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ElysStaked
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
