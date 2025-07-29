package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsPrefix)
	if b == nil {
		return types.DefaultParams()
	}

	var v types.Params
	if err := k.cdc.Unmarshal(b, &v); err != nil {
		ctx.Logger().Error("failed to unmarshal params, using defaults", "error", err)
		return types.DefaultParams()
	}
	return v
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b, err := k.cdc.Marshal(&p)
	if err != nil {
		return err
	}
	store.Set(types.ParamsPrefix, b)
	return nil
}
