package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

const MaxElysAmount = 3_218_460_000_000
const MaxEdenAmount = 3_441_097_000_000

func (k Keeper) GetAtomStaker(ctx sdk.Context, address sdk.AccAddress) (val types.AtomStaker) {
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

func (k Keeper) SetAtomStaker(ctx sdk.Context, val types.AtomStaker) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetAtomStakerKey(sdk.MustAccAddressFromBech32(val.Address))
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}

func (k Keeper) GetAllAtomStakers(ctx sdk.Context) (list []*types.AtomStaker) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.AtomStakersKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AtomStaker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetGovernor(ctx sdk.Context, address sdk.AccAddress) (val types.Governor) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetGovernorKey(address))

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.Address = address.String()
		val.Amount = math.ZeroInt()
	}
	return
}

func (k Keeper) SetGovernor(ctx sdk.Context, val types.Governor) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetGovernorKey(sdk.MustAccAddressFromBech32(val.Address))
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}

func (k Keeper) GetAllGovernors(ctx sdk.Context) (list []*types.Governor) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.GovernorKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Governor
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetNFTHolder(ctx sdk.Context, address sdk.AccAddress) (val types.NftHolder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetNFTHolderKey(address))

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.Address = address.String()
		val.Amount = math.ZeroInt()
	}
	return
}

func (k Keeper) SetNFTHodler(ctx sdk.Context, val types.NftHolder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetNFTHolderKey(sdk.MustAccAddressFromBech32(val.Address))
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}

func (k Keeper) GetAllNFTHolders(ctx sdk.Context) (list []*types.NftHolder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.NFTHoldersKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NftHolder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetCadet(ctx sdk.Context, address sdk.AccAddress) (val types.Cadet) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetCadetKey(address))

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.Address = address.String()
		val.Amount = math.ZeroInt()
	}
	return
}

func (k Keeper) SetCadet(ctx sdk.Context, val types.Cadet) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetCadetKey(sdk.MustAccAddressFromBech32(val.Address))
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}

func (k Keeper) GetAllCadets(ctx sdk.Context) (list []*types.Cadet) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CadetsKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Cadet
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetAirdropClaimed(ctx sdk.Context, address sdk.AccAddress) (val types.AirdropClaimed) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetAirdropClaimedKey(address))

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.Claimed = false
		val.Address = address.String()
	}
	return
}

func (k Keeper) SetAirdropClaimed(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetAirdropClaimedKey(address)
	val := types.AirdropClaimed{
		Address: address.String(),
		Claimed: true,
	}
	b := k.cdc.MustMarshal(&val)
	store.Set(key, b)
}

func (k Keeper) GetAllAirdropClaimed(ctx sdk.Context) (list []*types.AirdropClaimed) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.AirdropClaimedKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AirdropClaimed
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetTotalClaimed(ctx sdk.Context) (val types.TotalClaimed) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.TotalClaimedKeyPrefix)

	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
	} else {
		val.TotalEdenClaimed = math.ZeroInt()
		val.TotalElysClaimed = math.ZeroInt()
	}
	return
}

func (k Keeper) SetTotalClaimed(ctx sdk.Context, totalClaimed types.TotalClaimed) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&totalClaimed)
	store.Set(types.TotalClaimedKeyPrefix, b)
}
