package keeper

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k Keeper) GetUserData(ctx sdk.Context, address string, vaultId uint64) (userData types.UserData, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetUserDataKey(address, vaultId))
	if b == nil {
		return types.UserData{
			TotalDepositsUsd:    sdkmath.LegacyZeroDec(),
			TotalWithdrawalsUsd: sdkmath.LegacyZeroDec(),
			EdenUsdValue:        sdkmath.LegacyZeroDec(),
			EdenAmount:          sdkmath.ZeroInt(),
		}, false
	}

	k.cdc.MustUnmarshal(b, &userData)
	return userData, true
}

func (k Keeper) GetAllUserDatas(ctx sdk.Context) []types.UserData {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.UserDataKeyPrefix)
	defer iterator.Close()

	var userDatas []types.UserData
	for ; iterator.Valid(); iterator.Next() {
		var userData types.UserData
		k.cdc.MustUnmarshal(iterator.Value(), &userData)
		userDatas = append(userDatas, userData)
	}
	return userDatas
}

func (k Keeper) SetUserData(ctx sdk.Context, address string, vaultId uint64, userData types.UserData) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&userData)
	store.Set(types.GetUserDataKey(address, vaultId), b)
	return nil
}
