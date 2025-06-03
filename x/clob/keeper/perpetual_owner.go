package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetualOwner(ctx sdk.Context, owner sdk.AccAddress, subAccountId, marketId, perpetualId uint64) (types.PerpetualOwner, bool) {
	key := types.GetPerpetualOwnerKey(owner, subAccountId, marketId, perpetualId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.PerpetualOwner{}, false
	}

	var v types.PerpetualOwner
	k.cdc.MustUnmarshal(b, &v)
	return v, true
}

func (k Keeper) GetAllPerpetualOwners(ctx sdk.Context) []types.PerpetualOwner {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOwnerPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualOwner

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOwner
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetualOwner(ctx sdk.Context, v types.PerpetualOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOwnerKey(v.GetOwnerAccAddress(), v.SubAccountId, v.MarketId, v.PerpetualId)
	b := k.cdc.MustMarshal(&v)
	store.Set(key, b)
}

func (k Keeper) DeletePerpetualOwner(ctx sdk.Context, v types.PerpetualOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOwnerKey(v.GetOwnerAccAddress(), v.SubAccountId, v.MarketId, v.PerpetualId)
	store.Delete(key)
}

// CheckAndGetPerpetualOwner Since we follow a net accounting model,
// for each market a subaccount will always have only one position
// The total number of perpetuals per subaccount is limited by total number of markets
func (k Keeper) CheckAndGetPerpetualOwner(ctx sdk.Context, subAccount types.SubAccount, marketId uint64) (types.PerpetualOwner, bool) {
	prefixKey := types.GetPerpetualOwnerAddressKey(subAccount.GetOwnerAccAddress())
	prefixKey = append(prefixKey, sdk.Uint64ToBigEndian(subAccount.Id)...)
	prefixKey = append(prefixKey, []byte("/")...)

	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), prefixKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOwner
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.MarketId == marketId {
			return val, true
		}
	}

	return types.PerpetualOwner{}, false
}
