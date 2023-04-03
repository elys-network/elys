package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

// SetAirdrop set a specific airdrop in the store from its index
func (k Keeper) SetAirdrop(ctx sdk.Context, airdrop types.Airdrop) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))
	b := k.cdc.MustMarshal(&airdrop)
	store.Set(types.AirdropKey(
		airdrop.Intent,
	), b)
}

// GetAirdrop returns a airdrop from its index
func (k Keeper) GetAirdrop(
	ctx sdk.Context,
	intent string,

) (val types.Airdrop, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))

	b := store.Get(types.AirdropKey(
		intent,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAirdrop removes a airdrop from the store
func (k Keeper) RemoveAirdrop(
	ctx sdk.Context,
	intent string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))
	store.Delete(types.AirdropKey(
		intent,
	))
}

// GetAllAirdrop returns all airdrop
func (k Keeper) GetAllAirdrop(ctx sdk.Context) (list []types.Airdrop) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Airdrop
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
