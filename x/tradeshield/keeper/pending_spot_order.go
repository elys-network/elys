package keeper

import (
	"encoding/binary"
	"errors"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// GetPendingSpotOrderCount get the total number of pendingSpotOrder
func (k Keeper) GetPendingSpotOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PendingSpotOrderCountKey
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	// Set count
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

	if count == 0 {
		k.SetPendingSpotOrderCount(ctx, uint64(1))
		count = 1
	}

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

func (k Keeper) SpotBinarySearch(ctx sdk.Context, orderPrice sdk.Dec, orders []uint64) (int, error) {
	low, high := 0, len(orders)
	for low < high {
		mid := (low + high) / 2
		// Get order price
		order, found := k.GetPendingSpotOrder(ctx, orders[mid])
		if !found {
			return 0, types.ErrSpotOrderNotFound
		}
		if order.OrderPrice.Rate.LT(orderPrice) {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low, nil
}

func (k Keeper) InsertSpotSortedOrder(ctx sdk.Context, newOrder types.SpotOrder) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SortedSpotOrderKey)

	key, err := types.GenSpotKey(newOrder)
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

	index, err := k.SpotBinarySearch(ctx, newOrder.OrderPrice.Rate, orderIds)
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

// GetAllPendingSpotOrder returns all pendingSpotOrder
func (k Keeper) GetAllSortedSpotOrder(ctx sdk.Context) (list [][]uint64, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SortedSpotOrderKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var orderIds []uint64
		orderIds, err := types.DecodeUint64Slice(iterator.Value())
		if err != nil {
			return nil, err
		}
		list = append(list, orderIds)
	}

	return
}

// RemoveSortedOrder removes an order from the sorted order list.
func (k Keeper) RemoveSpotSortedOrder(ctx sdk.Context, orderID uint64) error {
	order, found := k.GetPendingSpotOrder(ctx, orderID)
	if !found {
		return types.ErrSpotOrderNotFound
	}

	// Generate the key for the order
	key, err := types.GenSpotKey(order)
	if err != nil {
		return err
	}

	// Load the sorted order IDs using the key
	sortedStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SortedSpotOrderKey)
	bz := sortedStore.Get([]byte(key))

	if bz == nil {
		return errors.New("sorted order IDs not found")
	}

	orderIds, err := types.DecodeUint64Slice(bz)
	if err != nil {
		return err
	}

	// Find the index of the order ID in the sorted order list
	index, err := k.SpotBinarySearch(ctx, order.OrderPrice.Rate, orderIds)
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

// ExecuteStopLossSpotOrder executes a stop loss order
func (k Keeper) ExecuteStopLossSpotOrder(ctx sdk.Context, order types.SpotOrder) error {
	// throws not implemented error
	return errors.New("not implemented")
}

// ExecuteLimitSellOrder executes a limit sell order
func (k Keeper) ExecuteLimitSellOrder(ctx sdk.Context, order types.SpotOrder) error {
	// throws not implemented error
	return errors.New("not implemented")
}

// ExecuteLimitBuyOrder executes a limit buy order
func (k Keeper) ExecuteLimitBuyOrder(ctx sdk.Context, order types.SpotOrder) error {
	// throws not implemented error
	return errors.New("not implemented")
}

// ExecuteMarketBuyOrder executes a market buy order
func (k Keeper) ExecuteMarketBuyOrder(ctx sdk.Context, order types.SpotOrder) error {
	// Swap the order amount with the target denom
	k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderAmount.Denom,
		DenomOut:  order.OrderTargetDenom,
		Discount:  sdk.ZeroDec(),
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdk.ZeroInt()),
	})

	return nil
}
