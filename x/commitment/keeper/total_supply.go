package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k Keeper) GetTotalSupply(ctx sdk.Context) (val types.TotalSupply) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := store.Get(types.TotalSupplyKeyPrefix)

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.TotalEdenSupply = math.ZeroInt()
		val.TotalEdenbSupply = math.ZeroInt()
	}
	return
}

func (k Keeper) SetTotalSupply(ctx sdk.Context, val types.TotalSupply) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&val)
	store.Set(types.TotalSupplyKeyPrefix, b)
}
