package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

func (k Keeper) TrackWeightBreakingSlippage(ctx sdk.Context, poolId uint64, token sdk.Coin) {
	price := k.oracleKeeper.GetDenomPrice(ctx, token.Denom)
	track := types.WeightBreakingSlippage{
		PoolId: poolId,
		Date:   ctx.BlockTime().Format("2006-01-02"),
		Amount: price.Mul(osmomath.BigDecFromSDKInt(token.Amount)).Dec(),
	}
	k.AddWeightAndSlippageFee(ctx, track)
}

// Returns last x days avg for weight breaking and slippage
func (k Keeper) GetWeightBreakingSlippageAvg(ctx sdk.Context, poolId uint64, days int) osmomath.BigDec {
	start := ctx.BlockTime()
	count := math.ZeroInt()
	total := osmomath.ZeroBigDec()

	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i*-1).Format("2006-01-02")
		info := k.GetWeightAndSlippageFee(ctx, poolId, date)

		if info.Amount.IsPositive() {
			total = total.Add(osmomath.BigDecFromDec(info.Amount))
			count = count.Add(math.OneInt())
		}
	}

	if count.IsZero() {
		return osmomath.ZeroBigDec()
	}
	return total.Quo(osmomath.BigDecFromSDKInt(count))
}
