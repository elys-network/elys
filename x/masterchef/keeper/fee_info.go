package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetFeeInfo(ctx sdk.Context, info types.FeeInfo) {
	store := ctx.KVStore(k.storeKey)

	// Get the current block time and format it as a string
	currentTime := ctx.BlockTime()
	dateString := currentTime.Format("2006-01-02") // YYYY-MM-DD format

	key := types.GetFeeInfoKey(dateString)
	b := k.cdc.MustMarshal(&info)
	store.Set(key, b)
}

func (k Keeper) GetFeeInfo(ctx sdk.Context, timestamp string) (val types.FeeInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetFeeInfoKey(timestamp)

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveFeeInfo(ctx sdk.Context, timestamp string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetFeeInfoKey(timestamp)
	if store.Has(key) {
		store.Delete(key)
	}
}

func (k Keeper) GetAllFeeInfos(ctx sdk.Context) (list []types.FeeInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FeeInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.FeeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
