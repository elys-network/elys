package keeper

import (
	"encoding/binary"
	"math"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// GetPendingSpotOrderCount get the total number of pendingSpotOrder
func (k Keeper) GetPendingSpotOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte{})
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
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte{})
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

	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	appendedValue := k.cdc.MustMarshal(&pendingSpotOrder)
	store.Set(GetPendingSpotOrderIDBytes(pendingSpotOrder.OrderId), appendedValue)

	// Update pendingSpotOrder count
	k.SetPendingSpotOrderCount(ctx, count+1)

	return count
}

// SetPendingSpotOrder set a specific pendingSpotOrder in the store
func (k Keeper) SetPendingSpotOrder(ctx sdk.Context, pendingSpotOrder types.SpotOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	b := k.cdc.MustMarshal(&pendingSpotOrder)
	store.Set(GetPendingSpotOrderIDBytes(pendingSpotOrder.OrderId), b)
}

// GetPendingSpotOrder returns a pendingSpotOrder from its id
func (k Keeper) GetPendingSpotOrder(ctx sdk.Context, id uint64) (val types.SpotOrder, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	b := store.Get(GetPendingSpotOrderIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetPendingSpotOrdersForAddress(ctx sdk.Context, address string, status *types.Status, pagination *query.PageRequest) ([]types.SpotOrder, *query.PageResponse, error) {
	var orders []types.SpotOrder

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	orderStore := prefix.NewStore(store, types.PendingSpotOrderKey)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.FilteredPaginate(orderStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var order types.SpotOrder
		err := k.cdc.Unmarshal(value, &order)
		if err == nil {
			if accumulate && order.OwnerAddress == address && (*status == types.Status_ALL || order.Status == *status) {
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

// RemovePendingSpotOrder removes a pendingSpotOrder from the store
func (k Keeper) RemovePendingSpotOrder(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	store.Delete(GetPendingSpotOrderIDBytes(id))
}

// GetAllPendingSpotOrder returns all pendingSpotOrder
func (k Keeper) GetAllPendingSpotOrder(ctx sdk.Context) (list []types.SpotOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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

// GetAllPendingSpotOrder returns all pendingSpotOrder
func (k Keeper) GetAllSortedSpotOrder(ctx sdk.Context) (list [][]uint64, err error) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.SortedSpotOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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

// DeleteAllPendingSpotOrder deleted all pendingSpotOrder
func (k Keeper) DeleteAllPendingSpotOrder(ctx sdk.Context) (list []types.SpotOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	return
}

// ExecuteStopLossSpotOrder executes a stop loss order
func (k Keeper) ExecuteStopLossOrder(ctx sdk.Context, order types.SpotOrder) error {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.OrderPrice.BaseDenom, order.OrderPrice.QuoteDenom)
	if err != nil {
		return err
	}

	if marketPrice.GT(order.OrderPrice.Rate) {
		// skip the order
		return nil
	}

	// send the order amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.OrderAmount))
	if err != nil {
		return err
	}

	// Swap the order amount with the target denom
	_, err = k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderPrice.BaseDenom,
		DenomOut:  order.OrderPrice.QuoteDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingSpotOrder(ctx, order.OrderId)

	return nil
}

// ExecuteLimitSellOrder executes a limit sell order
func (k Keeper) ExecuteLimitSellOrder(ctx sdk.Context, order types.SpotOrder) error {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.OrderPrice.BaseDenom, order.OrderPrice.QuoteDenom)
	if err != nil {
		return err
	}

	if marketPrice.LT(order.OrderPrice.Rate) {
		// skip the order
		return nil
	}

	// send the order amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.OrderAmount))
	if err != nil {
		return err
	}

	// Swap the order amount with the target denom
	_, err = k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderPrice.BaseDenom,
		DenomOut:  order.OrderPrice.QuoteDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingSpotOrder(ctx, order.OrderId)

	return nil
}

// ExecuteLimitBuyOrder executes a limit buy order
func (k Keeper) ExecuteLimitBuyOrder(ctx sdk.Context, order types.SpotOrder) error {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.OrderPrice.BaseDenom, order.OrderPrice.QuoteDenom)
	if err != nil {
		return err
	}

	if marketPrice.GT(order.OrderPrice.Rate) {
		// skip the order
		return nil
	}

	// send the order amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.OrderAmount))
	if err != nil {
		return err
	}

	// Swap the order amount with the target denom
	_, err = k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderPrice.BaseDenom,
		DenomOut:  order.OrderPrice.QuoteDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingSpotOrder(ctx, order.OrderId)

	return nil
}

// ExecuteMarketBuyOrder executes a market buy order
func (k Keeper) ExecuteMarketBuyOrder(ctx sdk.Context, order types.SpotOrder) error {
	// Swap the order amount with the target denom
	_, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderAmount.Denom,
		DenomOut:  order.OrderTargetDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
	})
	if err != nil {
		return err
	}

	return nil
}
