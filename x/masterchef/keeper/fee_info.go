package keeper

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/masterchef/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) AddFeeInfo(ctx sdk.Context, lp, stakers, protocol osmomath.BigDec, gas bool) {
	// Get the current block time and format it as a string
	currentTime := ctx.BlockTime()
	dateString := currentTime.Format("2006-01-02") // YYYY-MM-DD format

	info := k.GetFeeInfo(ctx, dateString)

	if gas {
		info.GasLp = info.GasLp.Add(lp.Dec().TruncateInt())
		info.GasStakers = info.GasStakers.Add(stakers.Dec().TruncateInt())
		info.GasProtocol = info.GasProtocol.Add(protocol.Dec().TruncateInt())
	} else {
		info.DexLp = info.DexLp.Add(lp.Dec().TruncateInt())
		info.DexStakers = info.DexStakers.Add(stakers.Dec().TruncateInt())
		info.DexProtocol = info.DexProtocol.Add(protocol.Dec().TruncateInt())
	}

	k.SetFeeInfo(ctx, info, dateString)
}

func (k Keeper) AddEdenInfo(ctx sdk.Context, eden osmomath.BigDec) {
	// Get the current block time and format it as a string
	currentTime := ctx.BlockTime()
	dateString := currentTime.Format("2006-01-02") // YYYY-MM-DD format

	info := k.GetFeeInfo(ctx, dateString)

	info.EdenLp = info.EdenLp.Add(eden.Dec().TruncateInt())

	k.SetFeeInfo(ctx, info, dateString)
}

// Deletes fee info that is older than 8 days
func (k Keeper) DeleteFeeInfo(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	currentTime := ctx.BlockTime().AddDate(0, 0, -8)
	dateString := currentTime.Format("2006-01-02") // YYYY-MM-DD format
	key := types.GetFeeInfoKey(dateString)

	if store.Has(key) {
		store.Delete(key)
	}
}

func (k Keeper) SetFeeInfo(ctx sdk.Context, info types.FeeInfo, timestamp string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	key := types.GetFeeInfoKey(timestamp)
	b := k.cdc.MustMarshal(&info)
	store.Set(key, b)
}

func (k Keeper) GetFeeInfo(ctx sdk.Context, timestamp string) (val types.FeeInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetFeeInfoKey(timestamp)

	b := store.Get(key)
	if b == nil {
		return types.FeeInfo{
			GasLp:        sdkmath.ZeroInt(),
			GasStakers:   sdkmath.ZeroInt(),
			GasProtocol:  sdkmath.ZeroInt(),
			DexLp:        sdkmath.ZeroInt(),
			DexStakers:   sdkmath.ZeroInt(),
			DexProtocol:  sdkmath.ZeroInt(),
			PerpLp:       sdkmath.ZeroInt(),
			PerpStakers:  sdkmath.ZeroInt(),
			PerpProtocol: sdkmath.ZeroInt(),
			EdenLp:       sdkmath.ZeroInt(),
		}
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}

func (k Keeper) RemoveFeeInfo(ctx sdk.Context, timestamp string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetFeeInfoKey(timestamp)
	if store.Has(key) {
		store.Delete(key)
	}
}

func (k Keeper) GetAllFeeInfos(ctx sdk.Context) (list []types.FeeInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.FeeInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.FeeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Returns last 7 days average of staker fees collected
func (k Keeper) GetAvgStakerFeesCollected(ctx sdk.Context, days int) osmomath.BigDec {
	start := ctx.BlockTime()
	count := sdkmath.ZeroInt()
	total := sdkmath.ZeroInt()

	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i*-1).Format("2006-01-02")
		info := k.GetFeeInfo(ctx, date)

		collected := info.DexStakers.Add(info.GasStakers).Add(info.PerpStakers)
		if collected.IsPositive() {
			total = total.Add(collected)
			count = count.Add(sdkmath.OneInt())
		}
	}

	if count.IsZero() {
		return osmomath.ZeroBigDec()
	}
	return osmomath.BigDecFromSDKInt(total).Quo(osmomath.BigDecFromSDKInt(count))
}
