package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsPrefix)
	if b == nil {
		panic("cannot found params")
	}

	var v types.Params
	k.cdc.MustUnmarshal(b, &v)
	return v
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&p)
	store.Set(types.ParamsPrefix, b)
}
