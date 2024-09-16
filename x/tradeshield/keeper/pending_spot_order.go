package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// GetPendingSpotOrderCount get the total number of pendingSpotOrder
func (k Keeper) GetPendingSpotOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PendingSpotOrderCountKey
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetPendingSpotOrderCount set the total number of pendingSpotOrder
func (k Keeper) SetPendingSpotOrderCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PendingSpotOrderCountKey
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPendingSpotOrder appends a pendingSpotOrder in the store with a new id and update the count
func (k Keeper) AppendPendingSpotOrder(
	ctx sdk.Context,
	pendingSpotOrder types.SpotOrder,
) uint64 {
	// Create the pendingSpotOrder
	count := k.GetPendingSpotOrderCount(ctx)

	// Set the ID of the appended value
	pendingSpotOrder.OrderId = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingSpotOrderKey)
	appendedValue := k.cdc.MustMarshal(&pendingSpotOrder)
	store.Set(GetPendingSpotOrderIDBytes(pendingSpotOrder.OrderId), appendedValue)

	// Update pendingSpotOrder count
	k.SetPendingSpotOrderCount(ctx, count+1)

	return count
}

// SetPendingSpotOrder set a specific pendingSpotOrder in the store
func (k Keeper) SetPendingSpotOrder(ctx sdk.Context, pendingSpotOrder types.SpotOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingSpotOrderKey)
	b := k.cdc.MustMarshal(&pendingSpotOrder)
	store.Set(GetPendingSpotOrderIDBytes(pendingSpotOrder.OrderId), b)
}

// GetPendingSpotOrder returns a pendingSpotOrder from its id
func (k Keeper) GetPendingSpotOrder(ctx sdk.Context, id uint64) (val types.SpotOrder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingSpotOrderKey)
	b := store.Get(GetPendingSpotOrderIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePendingSpotOrder removes a pendingSpotOrder from the store
func (k Keeper) RemovePendingSpotOrder(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingSpotOrderKey)
	store.Delete(GetPendingSpotOrderIDBytes(id))
}

// GetAllPendingSpotOrder returns all pendingSpotOrder
func (k Keeper) GetAllPendingSpotOrder(ctx sdk.Context) (list []types.SpotOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingSpotOrderKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SpotOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPendingSpotOrderIDBytes returns the byte representation of the ID
func GetPendingSpotOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetPendingSpotOrderIDFromBytes returns ID in uint64 format from a byte array
func GetPendingSpotOrderIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
