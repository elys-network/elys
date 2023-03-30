package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// SetAssetInfo set a specific assetInfo in the store from its index
func (k Keeper) SetAssetInfo(ctx sdk.Context, assetInfo types.AssetInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
	bz := k.cdc.MustMarshal(&assetInfo)
	store.Set(types.AssetInfoKey(assetInfo.Denom), bz)
}

// GetAssetInfo returns a assetInfo from its index
func (k Keeper) GetAssetInfo(
	ctx sdk.Context,
	denom string,
) (val types.AssetInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
	bz := store.Get(types.AssetInfoKey(denom))
	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveAssetInfo removes a assetInfo from the store
func (k Keeper) RemoveAssetInfo(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
	store.Delete(types.AssetInfoKey(denom))
}

// GetAllAssetInfo returns all assetInfo
func (k Keeper) GetAllAssetInfo(ctx sdk.Context) (list []types.AssetInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
