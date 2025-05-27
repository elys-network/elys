package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
)

func (k Keeper) GetRewardProgram(ctx sdk.Context, address sdk.AccAddress) (val types.RewardProgram) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetRewardProgramKey(address))

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.Address = address.String()
		val.Amount = math.ZeroInt()
	}
	return
}

func (k Keeper) SetRewardProgram(ctx sdk.Context, val types.RewardProgram) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetRewardProgramKey(sdk.MustAccAddressFromBech32(val.Address))
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}

func (k Keeper) GetAllRewardPrograms(ctx sdk.Context) (list []*types.RewardProgram) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.RewardProgramKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RewardProgram
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetTotalRewardProgramClaimed(ctx sdk.Context) (val types.TotalClaimed) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.TotalRewardProgramClaimedKeyPrefix)

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.TotalEdenClaimed = math.ZeroInt()
		val.TotalElysClaimed = math.ZeroInt()
	}
	return
}

func (k Keeper) SetTotalRewardProgramClaimed(ctx sdk.Context, totalClaimed types.TotalClaimed) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&totalClaimed)
	store.Set(types.TotalRewardProgramClaimedKeyPrefix, b)
}
