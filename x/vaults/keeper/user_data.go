package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k Keeper) SetUserData(ctx sdk.Context, userData types.UserData) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetUserDataKey(userData.User), k.cdc.MustMarshal(&userData))
}

func (k Keeper) GetUserData(ctx sdk.Context, user string) (types.UserData, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := store.Get(types.GetUserDataKey(user))
	if b == nil {
		return types.UserData{}, false
	}

	var userData types.UserData
	k.cdc.MustUnmarshal(b, &userData)
	return userData, true
}
