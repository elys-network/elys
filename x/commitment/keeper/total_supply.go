package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k Keeper) GetTotalEdenSupply(ctx sdk.Context, address sdk.AccAddress) (val types.AtomStaker) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetAtomStakerKey(address))

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.Address = address.String()
		val.Amount = math.ZeroInt()
	}
	return
}

func (k Keeper) SetTotalEdenSupply(ctx sdk.Context, val types.AtomStaker) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetAtomStakerKey(sdk.MustAccAddressFromBech32(val.Address))
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}
