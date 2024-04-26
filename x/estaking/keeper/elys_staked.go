package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

// SetElysStaked set a specific elysStaked in the store from its index
func (k Keeper) SetElysStaked(ctx sdk.Context, elysStaked types.ElysStaked) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakedKeyPrefix))
	b := k.cdc.MustMarshal(&elysStaked)
	store.Set(types.ElysStakedKey(elysStaked.Address), b)
}

// GetElysStaked returns a elysStaked from its index
func (k Keeper) GetElysStaked(ctx sdk.Context, address string) math.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakedKeyPrefix))

	b := store.Get(types.ElysStakedKey(address))
	if b == nil {
		return math.ZeroInt()
	}

	val := types.ElysStaked{}
	k.cdc.MustUnmarshal(b, &val)
	return val.Amount
}

// RemoveElysStaked removes a elysStaked from the store
func (k Keeper) RemoveElysStaked(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakedKeyPrefix))
	store.Delete(types.ElysStakedKey(address))
}

// GetAllElysStaked returns all elysStaked
func (k Keeper) GetAllElysStaked(ctx sdk.Context) (list []types.ElysStaked) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakedKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ElysStaked
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
