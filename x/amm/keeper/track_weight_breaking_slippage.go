package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) GetWeightAndSlippageFee(ctx sdk.Context, poolId uint64, date string) types.WeightBreakingSlippage {
	track := types.WeightBreakingSlippage{}
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.WeightAndSlippagePrefix))
	bz := store.Get(types.WeightAndSlippageFeeKey(poolId, date))
	if len(bz) == 0 {
		return types.WeightBreakingSlippage{
			PoolId: poolId,
			Date:   date,
			Amount: math.LegacyZeroDec(),
		}
	}

	k.cdc.MustUnmarshal(bz, &track)
	return track
}

func (k Keeper) SetWeightAndSlippageFee(ctx sdk.Context, track types.WeightBreakingSlippage) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.WeightAndSlippagePrefix))
	bz := k.cdc.MustMarshal(&track)
	store.Set(types.WeightAndSlippageFeeKey(track.PoolId, track.Date), bz)
}

func (k Keeper) DeleteWeightAndSlippageFee(ctx sdk.Context, poolId uint64, date string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.WeightAndSlippagePrefix))
	store.Delete(types.WeightAndSlippageFeeKey(poolId, date))
}

func (k Keeper) AddWeightAndSlippageFee(ctx sdk.Context, track types.WeightBreakingSlippage) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.WeightAndSlippagePrefix))

	trackTotal := types.WeightBreakingSlippage{}
	bz := store.Get(types.WeightAndSlippageFeeKey(track.PoolId, track.Date))
	if len(bz) == 0 {
		trackTotal = track
	} else {
		k.cdc.MustUnmarshal(bz, &trackTotal)
		trackTotal.Amount = trackTotal.Amount.Add(track.Amount)
	}

	bz = k.cdc.MustMarshal(&trackTotal)
	store.Set(types.WeightAndSlippageFeeKey(track.PoolId, track.Date), bz)
}

func (k Keeper) TrackWeightBreakingSlippage(ctx sdk.Context, pool types.Pool, slippage sdk.Coin, weightBreakingFee math.LegacyDec) {
	track := types.WeightBreakingSlippage{
		PoolId: pool.PoolId,
		Date:   ctx.BlockTime().Format("2006-01-02"),
		Amount: weightBreakingFee,
	}
	if weightBreakingFee.IsNegative() {
		track.Amount = weightBreakingFee.Abs()
	}
	k.AddWeightAndSlippageFee(ctx, track)
}
