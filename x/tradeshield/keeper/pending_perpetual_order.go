package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// GetPendingPerpetualOrderCount get the total number of pendingPerpetualOrder
func (k Keeper) GetPendingPerpetualOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PendingPerpetualOrderCountKey
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetPendingPerpetualOrderCount set the total number of pendingPerpetualOrder
func (k Keeper) SetPendingPerpetualOrderCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PendingPerpetualOrderCountKey
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPendingPerpetualOrder appends a pendingPerpetualOrder in the store with a new id and update the count
func (k Keeper) AppendPendingPerpetualOrder(
	ctx sdk.Context,
	pendingPerpetualOrder types.PerpetualOrder,
) uint64 {
	// Create the pendingPerpetualOrder
	count := k.GetPendingPerpetualOrderCount(ctx)

	// Set the ID of the appended value
	pendingPerpetualOrder.OrderId = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingPerpetualOrderKey)
	appendedValue := k.cdc.MustMarshal(&pendingPerpetualOrder)
	store.Set(GetPendingPerpetualOrderIDBytes(pendingPerpetualOrder.OrderId), appendedValue)

	// Update pendingPerpetualOrder count
	k.SetPendingPerpetualOrderCount(ctx, count+1)

	return count
}

// SetPendingPerpetualOrder set a specific pendingPerpetualOrder in the store
func (k Keeper) SetPendingPerpetualOrder(ctx sdk.Context, pendingPerpetualOrder types.PerpetualOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingPerpetualOrderKey)
	b := k.cdc.MustMarshal(&pendingPerpetualOrder)
	store.Set(GetPendingPerpetualOrderIDBytes(pendingPerpetualOrder.OrderId), b)
}

// GetPendingPerpetualOrder returns a pendingPerpetualOrder from its id
func (k Keeper) GetPendingPerpetualOrder(ctx sdk.Context, id uint64) (val types.PerpetualOrder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingPerpetualOrderKey)
	b := store.Get(GetPendingPerpetualOrderIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePendingPerpetualOrder removes a pendingPerpetualOrder from the store
func (k Keeper) RemovePendingPerpetualOrder(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingPerpetualOrderKey)
	store.Delete(GetPendingPerpetualOrderIDBytes(id))
}

// GetAllPendingPerpetualOrder returns all pendingPerpetualOrder
func (k Keeper) GetAllPendingPerpetualOrder(ctx sdk.Context) (list []types.PerpetualOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PendingPerpetualOrderKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPendingPerpetualOrderIDBytes returns the byte representation of the ID
func GetPendingPerpetualOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetPendingPerpetualOrderIDFromBytes returns ID in uint64 format from a byte array
func GetPendingPerpetualOrderIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
