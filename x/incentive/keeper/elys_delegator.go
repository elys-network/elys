package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

// SetElysDelegator set a specific elysDelegator in the store from its index
func (k Keeper) SetElysDelegator(ctx sdk.Context, elysDelegator types.ElysDelegator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysDelegatorKeyPrefix))
	b := k.cdc.MustMarshal(&elysDelegator)
	store.Set(types.ElysDelegatorKey(
		elysDelegator.Index,
	), b)
}

// GetElysDelegator returns a elysDelegator from its index
func (k Keeper) GetElysDelegator(
	ctx sdk.Context,
	index string,

) (val types.ElysDelegator, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysDelegatorKeyPrefix))

	b := store.Get(types.ElysDelegatorKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveElysDelegator removes a elysDelegator from the store
func (k Keeper) RemoveElysDelegator(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysDelegatorKeyPrefix))
	store.Delete(types.ElysDelegatorKey(
		index,
	))
}

// GetAllElysDelegator returns all elysDelegator
func (k Keeper) GetAllElysDelegator(ctx sdk.Context) (list []types.ElysDelegator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysDelegatorKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ElysDelegator
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetElysDelegator returns a elysDelegator from delegator & validator address pair
func (k Keeper) GetElysDelegatorFromAddresses(
	ctx sdk.Context,
	delegator string,
	validator string,

) bool {
	delegators := k.GetAllElysDelegator(ctx)
	for _, d := range delegators {
		if strings.EqualFold(d.DelegatorAddr, delegator) && strings.EqualFold(d.ValidatorAddr, validator) {
			return true
		}
	}

	return false
}

// GetTotalItemsCount return total number of items
func (k Keeper) GetTotalElysDelegationItemCount(
	ctx sdk.Context,
) int {
	listElysDelegators := k.GetAllElysDelegator(ctx)

	return len(listElysDelegators)
}
