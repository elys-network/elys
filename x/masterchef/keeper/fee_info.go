package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) AddFeeInfo(ctx sdk.Context, lp, stakers, protocol sdk.Dec, gas bool) {
	// Get the current block time and format it as a string
	currentTime := ctx.BlockTime()
	dateString := currentTime.Format("2006-01-02") // YYYY-MM-DD format

	info := k.GetFeeInfo(ctx, dateString)

	if gas {
		info.GasLp = info.GasLp.Add(lp.TruncateInt())
		info.GasStakers = info.GasStakers.Add(stakers.TruncateInt())
		info.GasProtocol = info.GasProtocol.Add(protocol.TruncateInt())
	} else {
		info.DexLp = info.DexLp.Add(lp.TruncateInt())
		info.DexStakers = info.DexStakers.Add(stakers.TruncateInt())
		info.DexProtocol = info.DexProtocol.Add(protocol.TruncateInt())
	}

	k.SetFeeInfo(ctx, info, dateString)
}

func (k Keeper) AddEdenInfo(ctx sdk.Context, eden sdk.Dec) {
	// Get the current block time and format it as a string
	currentTime := ctx.BlockTime()
	dateString := currentTime.Format("2006-01-02") // YYYY-MM-DD format

	info := k.GetFeeInfo(ctx, dateString)

	info.EdenLp = info.EdenLp.Add(eden.TruncateInt())

	k.SetFeeInfo(ctx, info, dateString)
}

func (k Keeper) SetFeeInfo(ctx sdk.Context, info types.FeeInfo, timestamp string) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetFeeInfoKey(timestamp)
	b := k.cdc.MustMarshal(&info)
	store.Set(key, b)
}

func (k Keeper) GetFeeInfo(ctx sdk.Context, timestamp string) (val types.FeeInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetFeeInfoKey(timestamp)

	b := store.Get(key)
	if b == nil {
		return types.FeeInfo{
			GasLp:        sdk.ZeroInt(),
			GasStakers:   sdk.ZeroInt(),
			GasProtocol:  sdk.ZeroInt(),
			DexLp:        sdk.ZeroInt(),
			DexStakers:   sdk.ZeroInt(),
			DexProtocol:  sdk.ZeroInt(),
			PerpLp:       sdk.ZeroInt(),
			PerpStakers:  sdk.ZeroInt(),
			PerpProtocol: sdk.ZeroInt(),
			EdenLp:       sdk.ZeroInt(),
		}
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
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
