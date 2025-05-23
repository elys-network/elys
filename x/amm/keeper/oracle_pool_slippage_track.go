package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
)

func (k Keeper) GetSlippageTrack(ctx sdk.Context, poolId uint64, timestamp uint64) types.OraclePoolSlippageTrack {
	track := types.OraclePoolSlippageTrack{}
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.OraclePoolSlippageTrackPrefix))
	bz := store.Get(types.OraclePoolSlippageTrackKey(poolId, timestamp))
	if len(bz) == 0 {
		return types.OraclePoolSlippageTrack{
			PoolId:    poolId,
			Timestamp: uint64(ctx.BlockTime().Unix()),
			Tracked:   sdk.Coins{},
		}
	}

	k.cdc.MustUnmarshal(bz, &track)
	return track
}

func (k Keeper) SetSlippageTrack(ctx sdk.Context, track types.OraclePoolSlippageTrack) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.OraclePoolSlippageTrackPrefix))
	bz := k.cdc.MustMarshal(&track)
	store.Set(types.OraclePoolSlippageTrackKey(track.PoolId, track.Timestamp), bz)
}

func (k Keeper) DeleteSlippageTrack(ctx sdk.Context, track types.OraclePoolSlippageTrack) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.OraclePoolSlippageTrackPrefix))
	store.Delete(types.OraclePoolSlippageTrackKey(track.PoolId, track.Timestamp))
}

func (k Keeper) AllSlippageTracks(ctx sdk.Context) []types.OraclePoolSlippageTrack {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.OraclePoolSlippageTrackPrefix))

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	tracks := []types.OraclePoolSlippageTrack{}
	for ; iterator.Valid(); iterator.Next() {
		track := types.OraclePoolSlippageTrack{}
		k.cdc.MustUnmarshal(iterator.Value(), &track)

		tracks = append(tracks, track)
	}
	return tracks
}

func (k Keeper) GetLastSlippageTrack(ctx sdk.Context, poolId uint64) types.OraclePoolSlippageTrack {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.OraclePoolSlippageTrackPrefix))

	iterator := storetypes.KVStoreReversePrefixIterator(store, sdk.Uint64ToBigEndian(poolId))
	defer iterator.Close()

	track := types.OraclePoolSlippageTrack{}
	if iterator.Valid() {
		k.cdc.MustUnmarshal(iterator.Value(), &track)
		return track
	}

	return types.OraclePoolSlippageTrack{
		PoolId:    poolId,
		Timestamp: uint64(ctx.BlockTime().Unix()),
		Tracked:   sdk.Coins{},
	}
}

func (k Keeper) GetFirstSlippageTrack(ctx sdk.Context, poolId uint64) types.OraclePoolSlippageTrack {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.OraclePoolSlippageTrackPrefix))

	iterator := storetypes.KVStorePrefixIterator(store, sdk.Uint64ToBigEndian(poolId))
	defer iterator.Close()

	track := types.OraclePoolSlippageTrack{}
	if iterator.Valid() {
		k.cdc.MustUnmarshal(iterator.Value(), &track)
		return track
	}

	return types.OraclePoolSlippageTrack{
		PoolId:    poolId,
		Timestamp: uint64(ctx.BlockTime().Unix()),
		Tracked:   sdk.Coins{},
	}
}

func (k Keeper) GetTrackedSlippageDiff(ctx sdk.Context, poolId uint64) sdk.Coins {
	lastTrack := k.GetLastSlippageTrack(ctx, poolId)
	firstTrack := k.GetFirstSlippageTrack(ctx, poolId)
	return lastTrack.Tracked.Sub(firstTrack.Tracked...)
}

func (k Keeper) TrackSlippage(ctx sdk.Context, poolId uint64, amount sdk.Coin) {
	lastTrack := k.GetLastSlippageTrack(ctx, poolId)
	lastTrack.Tracked = lastTrack.Tracked.Add(amount)
	lastTrack.Timestamp = uint64(ctx.BlockTime().Unix())
	k.SetSlippageTrack(ctx, lastTrack)
}
