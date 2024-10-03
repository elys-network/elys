package keeper

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// GetPendingPerpetualOrderCount get the total number of pendingPerpetualOrder
func (k Keeper) GetPendingPerpetualOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.PendingPerpetualOrderCountKey
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	// Set count
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

	if count == 0 {
		k.SetPendingPerpetualOrderCount(ctx, uint64(1))
		count = 1
	}

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

func (k Keeper) GetPendingPerpetualOrdersForAddress(ctx sdk.Context, address string, pagination *query.PageRequest) ([]types.PerpetualOrder, *query.PageResponse, error) {
	var orders []types.PerpetualOrder

	store := ctx.KVStore(k.storeKey)
	orderStore := prefix.NewStore(store, types.PendingPerpetualOrderKey)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.FilteredPaginate(orderStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var order types.PerpetualOrder
		err := k.cdc.Unmarshal(value, &order)
		if err == nil {
			if accumulate && order.OwnerAddress == address {
				orders = append(orders, order)
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, nil, err
	}

	return orders, pageRes, nil
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
			return 0, types.ErrPerpetualOrderNotFound
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

func (k Keeper) GetAllSortedPerpetualOrder(ctx sdk.Context) (list [][]uint64, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SortedPerpetualOrderKey)
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
func (k Keeper) RemovePerpetualSortedOrder(ctx sdk.Context, orderID uint64) error {
	order, found := k.GetPendingPerpetualOrder(ctx, orderID)
	if !found {
		return types.ErrPerpetualOrderNotFound
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

// ExecuteLimitOpenOrder executes a limit open order
func (k Keeper) ExecuteLimitOpenOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.TriggerPrice.BaseDenom, order.TriggerPrice.QuoteDenom)
	if err != nil {
		return err
	}

	found := false

	switch order.Position {
	case types.PerpetualPosition_LONG:
		if marketPrice.LTE(order.TriggerPrice.Rate) {
			_, err := k.perpetual.Open(ctx, &perpetualtypes.MsgOpen{
				Creator:         order.OwnerAddress,
				Position:        perpetualtypes.Position(order.Position),
				Leverage:        order.Leverage,
				TradingAsset:    order.TradingAsset,
				Collateral:      order.Collateral,
				TakeProfitPrice: order.TakeProfitPrice,
				StopLossPrice:   order.StopLossPrice,
			}, false)
			if err != nil {
				return err
			}

			found = true
		}
	case types.PerpetualPosition_SHORT:
		if marketPrice.GTE(order.TriggerPrice.Rate) {
			_, err := k.perpetual.Open(ctx, &perpetualtypes.MsgOpen{
				Creator:         order.OwnerAddress,
				Position:        perpetualtypes.Position(order.Position),
				Leverage:        order.Leverage,
				TradingAsset:    order.TradingAsset,
				Collateral:      order.Collateral,
				TakeProfitPrice: order.TakeProfitPrice,
				StopLossPrice:   order.StopLossPrice,
			}, false)
			if err != nil {
				return err
			}

			found = true
		}
	}

	if found {
		// Remove the order from the pending order list
		k.RemovePendingPerpetualOrder(ctx, order.OrderId)

		return nil
	}

	// skip the order
	return nil
}

// ExecuteLimitCloseOrder executes a limit close order
func (k Keeper) ExecuteLimitCloseOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.TriggerPrice.BaseDenom, order.TriggerPrice.QuoteDenom)
	if err != nil {
		return err
	}

	found := false

	switch order.Position {
	case types.PerpetualPosition_LONG:
		if marketPrice.GTE(order.TriggerPrice.Rate) {
			_, err := k.perpetual.Close(ctx, &perpetualtypes.MsgClose{
				Creator: order.OwnerAddress,
				Id:      order.PositionId,
				Amount:  sdk.ZeroInt(),
			})
			if err != nil {
				return err
			}

			found = true
		}
	case types.PerpetualPosition_SHORT:
		if marketPrice.LTE(order.TriggerPrice.Rate) {
			_, err := k.perpetual.Close(ctx, &perpetualtypes.MsgClose{
				Creator: order.OwnerAddress,
				Id:      order.PositionId,
				Amount:  sdk.ZeroInt(),
			})
			if err != nil {
				return err
			}

			found = true
		}
	}

	if found {
		// Remove the order from the pending order list
		k.RemovePendingPerpetualOrder(ctx, order.OrderId)

		return nil
	}

	// skip the order
	return nil
}

// ExecuteMarketOpenOrder executes a market open order
func (k Keeper) ExecuteMarketOpenOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	_, err := k.perpetual.Open(ctx, &perpetualtypes.MsgOpen{
		Creator:         order.OwnerAddress,
		Position:        perpetualtypes.Position(order.Position),
		Leverage:        order.Leverage,
		TradingAsset:    order.TradingAsset,
		Collateral:      order.Collateral,
		TakeProfitPrice: order.TakeProfitPrice,
		StopLossPrice:   order.StopLossPrice,
	}, false)
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, order.OrderId)

	return nil
}

// ExecuteMarketCloseOrder executes a market close order
func (k Keeper) ExecuteMarketCloseOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	_, err := k.perpetual.Close(ctx, &perpetualtypes.MsgClose{
		Creator: order.OwnerAddress,
		Id:      order.PositionId,
		Amount:  sdk.ZeroInt(),
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, order.OrderId)

	return nil
}
