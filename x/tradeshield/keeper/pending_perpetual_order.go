package keeper

import (
	"encoding/binary"
	"errors"

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

func (k Keeper) PerpetualBinarySearch(ctx sdk.Context, orderPrice sdk.Dec, orders []uint64) (int, error) {
	low, high := 0, len(orders)
	for low < high {
		mid := (low + high) / 2
		// Get order price
		order, found := k.GetPendingPerpetualOrder(ctx, orders[mid])
		if !found {
			return 0, types.ErrOrderNotFound
		}
		if order.TriggerPrice.Rate.LT(orderPrice) {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low, nil
}

func (k Keeper) InsertPerptualSortedOrder(ctx sdk.Context, newOrder types.PerpetualOrder) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SortedPerpetualOrderKey)

	key, err := types.GenPerpKey(newOrder)
	if err != nil {
		return err
	}
	bz := store.Get([]byte(key))

	var orderIds []uint64
	if bz != nil {
		orderIds, err = types.DecodeUint64Slice(bz)
		if err != nil {
			return err
		}
	}

	index, err := k.PerpetualBinarySearch(ctx, newOrder.TriggerPrice.Rate, orderIds)
	if err != nil {
		return err
	}

	if len(orderIds) <= index {
		orderIds = append(orderIds, newOrder.OrderId)
	} else {
		orderIds = append(orderIds[:index+1], orderIds[index:]...)
		orderIds[index] = newOrder.OrderId
	}

	bz = types.EncodeUint64Slice(orderIds)

	store.Set([]byte(key), bz)
	return nil
}

// RemoveSortedOrder removes an order from the sorted order list.
func (k Keeper) RemovePerpetualSortedOrder(ctx sdk.Context, orderID uint64, positionID uint64) error {
	order, found := k.GetPendingPerpetualOrder(ctx, orderID)
	if !found {
		return types.ErrOrderNotFound
	}

	// Generate the key for the order
	key, err := types.GenPerpKey(order)
	if err != nil {
		return err
	}

	// Load the sorted order IDs using the key
	sortedStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SortedPerpetualOrderKey)
	bz := sortedStore.Get([]byte(key))

	if bz == nil {
		return errors.New("sorted order IDs not found")
	}

	orderIds, err := types.DecodeUint64Slice(bz)
	if err != nil {
		return err
	}

	// Find the index of the order ID in the sorted order list
	index, err := k.PerpetualBinarySearch(ctx, order.TriggerPrice.Rate, orderIds)
	if err != nil {
		return err
	}

	sizeOfVec := len(orderIds)
	for index < sizeOfVec && orderIds[index] != orderID {
		index++
	}

	if index >= sizeOfVec {
		return errors.New("order ID not found in sorted order list")
	}

	// Remove the order ID from the list
	orderIds = append(orderIds[:index], orderIds[index+1:]...)

	// Save the updated list back to storage
	encodedOrderIds := types.EncodeUint64Slice(orderIds)

	sortedStore.Set([]byte(key), encodedOrderIds)
	return nil
}
