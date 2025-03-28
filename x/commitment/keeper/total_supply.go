package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) GetTotalSupply(ctx sdk.Context) (val types.TotalSupply) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := store.Get(types.TotalSupplyKeyPrefix)

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.TotalEdenSupply = math.ZeroInt()
		val.TotalEdenbSupply = math.ZeroInt()
		val.TotalEdenVested = math.ZeroInt()
	}
	return
}

func (k Keeper) SetTotalSupply(ctx sdk.Context, val types.TotalSupply) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&val)
	store.Set(types.TotalSupplyKeyPrefix, b)
}

func (k Keeper) V10_SetEdenEdenBSupply(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

	defer iterator.Close()

	totalEden := math.ZeroInt()
	totalEdenB := math.ZeroInt()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		for _, token := range val.CommittedTokens {
			if token.Denom == ptypes.Eden {
				totalEden = totalEden.Add(token.Amount)
			}
			if token.Denom == ptypes.EdenB {
				totalEdenB = totalEdenB.Add(token.Amount)
			}
		}

		totalEden = totalEden.Add(val.Claimed.AmountOf(ptypes.Eden))
		totalEdenB = totalEdenB.Add(val.Claimed.AmountOf(ptypes.EdenB))
	}

	k.SetTotalSupply(ctx, types.TotalSupply{
		TotalEdenSupply:  totalEden,
		TotalEdenbSupply: totalEdenB,
	})
}

func (k Keeper) V10_SetEdenVested(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

	defer iterator.Close()

	totalEdenVested := math.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		for _, token := range val.VestingTokens {
			totalEdenVested = totalEdenVested.Add(token.TotalAmount)
		}
	}

	total := k.GetTotalSupply(ctx)

	k.SetTotalSupply(ctx, types.TotalSupply{
		TotalEdenSupply:  total.TotalEdenSupply,
		TotalEdenbSupply: total.TotalEdenbSupply,
		TotalEdenVested:  totalEdenVested,
	})
}
