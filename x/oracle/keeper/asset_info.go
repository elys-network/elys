package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/oracle/types"
)

// SetAssetInfo set a specific assetInfo in the store from its index
func (k Keeper) SetAssetInfo(ctx sdk.Context, assetInfo types.AssetInfo) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AssetInfoKeyPrefix))
	bz := k.cdc.MustMarshal(&assetInfo)
	store.Set(types.AssetInfoKey(assetInfo.Denom), bz)
}

// GetAssetInfo returns a assetInfo from its index
func (k Keeper) GetAssetInfo(ctx sdk.Context, denom string) (val types.AssetInfo, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AssetInfoKeyPrefix))
	bz := store.Get(types.AssetInfoKey(denom))
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveAssetInfo removes a assetInfo from the store
func (k Keeper) RemoveAssetInfo(ctx sdk.Context, denom string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AssetInfoKeyPrefix))
	store.Delete(types.AssetInfoKey(denom))
}

// GetAllAssetInfo returns all assetInfo
func (k Keeper) GetAllAssetInfo(ctx sdk.Context) (list []types.AssetInfo) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AssetInfoKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
