package keeper

import (
	"encoding/binary"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
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

// GetAllSortedSpotOrder returns all sortedSpotOrder
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

// SetAllLegacySpotOrderPriceToNewOrderPriceStructure set all legacy spot order price to new order price structure
func (k Keeper) SetAllLegacySpotOrderPriceToNewOrderPriceStructure(ctx sdk.Context) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingSpotOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var order types.SpotOrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		order.OrderPrice = order.LegacyOrderPriceV1.Rate
		order.LegacyOrderPriceV1 = types.LegacyOrderPriceV1{}
		store.Set(iterator.Key(), k.cdc.MustMarshal(&order))
	}
}

// ExecuteStopLossOrder executes a stop loss order
func (k Keeper) ExecuteStopLossOrder(ctx sdk.Context, order types.SpotOrder) (*ammtypes.MsgSwapByDenomResponse, error) {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.OrderAmount.Denom, order.OrderTargetDenom)
	if err != nil {
		return nil, err
	}
	if marketPrice.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrZeroMarketPrice, "denom in: %s, denom out: %s", order.OrderAmount.Denom, order.OrderTargetDenom)
	}

	if marketPrice.GT(order.GetBigDecOrderPrice()) {
		// skip the order
		return nil, nil
	}

	// send the order amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.OrderAmount))
	if err != nil {
		return nil, err
	}

	// Swap the order amount with the target denom
	res, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderAmount.Denom,
		DenomOut:  order.OrderTargetDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
		MaxAmount: order.OrderAmount,
	})
	if err != nil {
		return res, err
	}

	// Remove the order from the pending order list
	k.RemovePendingSpotOrder(ctx, order.OrderId)

	// emit the event
	ctx.EventManager().EmitEvent(types.NewExecuteStopLossSpotOrderEvt(order, res))

	return res, nil
}

// ExecuteLimitSellOrder executes a limit sell order
func (k Keeper) ExecuteLimitSellOrder(ctx sdk.Context, order types.SpotOrder) (*ammtypes.MsgSwapByDenomResponse, error) {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.OrderAmount.Denom, order.OrderTargetDenom)
	if err != nil {
		return nil, err
	}
	if marketPrice.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrZeroMarketPrice, "denom in: %s, denom out: %s", order.OrderAmount.Denom, order.OrderTargetDenom)
	}

	if marketPrice.LT(order.GetBigDecOrderPrice()) {
		// skip the order
		return nil, nil
	}

	// send the order amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.OrderAmount))
	if err != nil {
		return nil, err
	}

	// Swap the order amount with the target denom
	res, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderAmount.Denom,
		DenomOut:  order.OrderTargetDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
		MaxAmount: order.OrderAmount,
	})
	if err != nil {
		return res, err
	}

	params := k.GetParams(ctx)
	expectedAmount := marketPrice.Mul(osmomath.BigDecFromSDKInt(order.OrderAmount.Amount))
	gotAmount := osmomath.BigDecFromSDKInt(res.Amount.Amount)
	tolerance := osmomath.ZeroBigDec()

	if gotAmount.LT(expectedAmount) {
		tolerance = (expectedAmount.Sub(gotAmount)).Quo(expectedAmount)
	}

	if tolerance.GT(params.GetBigDecTolerance()) {
		return res, errorsmod.Wrapf(types.ErrHighTolerance, "tolerance: %s", tolerance)
	}

	// Remove the order from the pending order list
	k.RemovePendingSpotOrder(ctx, order.OrderId)

	// emit the event
	ctx.EventManager().EmitEvent(types.NewExecuteLimitSellSpotOrderEvt(order, res))

	return res, nil
}

// ExecuteLimitBuyOrder executes a limit buy order
func (k Keeper) ExecuteLimitBuyOrder(ctx sdk.Context, order types.SpotOrder) (*ammtypes.MsgSwapByDenomResponse, error) {
	marketPrice, err := k.GetAssetPriceFromDenomInToDenomOut(ctx, order.OrderTargetDenom, order.OrderAmount.Denom)
	if err != nil {
		return nil, err
	}
	if marketPrice.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrZeroMarketPrice, "denom in: %s, denom out: %s", order.OrderAmount.Denom, order.OrderTargetDenom)
	}

	if marketPrice.GT(order.GetBigDecOrderPrice()) {
		// skip the order
		return nil, nil
	}

	// send the order amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.OrderAmount))
	if err != nil {
		return nil, err
	}

	// Swap the order amount with the target denom
	res, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderAmount.Denom,
		DenomOut:  order.OrderTargetDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
		MaxAmount: order.OrderAmount,
	})
	if err != nil {
		return res, err
	}

	params := k.GetParams(ctx)
	expectedAmount := osmomath.BigDecFromSDKInt(order.OrderAmount.Amount).Quo(marketPrice)
	gotAmount := osmomath.BigDecFromSDKInt(res.Amount.Amount)
	tolerance := osmomath.ZeroBigDec()

	if gotAmount.LT(expectedAmount) {
		tolerance = (expectedAmount.Sub(gotAmount)).Quo(expectedAmount)
	}

	if tolerance.GT(params.GetBigDecTolerance()) {
		return res, errorsmod.Wrapf(types.ErrHighTolerance, "tolerance: %s", tolerance)
	}

	// Remove the order from the pending order list
	k.RemovePendingSpotOrder(ctx, order.OrderId)

	// emit the event
	ctx.EventManager().EmitEvent(types.NewExecuteLimitBuySpotOrderEvt(order, res))

	return res, nil
}

// ExecuteMarketBuyOrder executes a market buy order
func (k Keeper) ExecuteMarketBuyOrder(ctx sdk.Context, order types.SpotOrder) (*ammtypes.MsgSwapByDenomResponse, error) {
	// Swap the order amount with the target denom
	res, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
		Sender:    order.OwnerAddress,
		Recipient: order.OwnerAddress,
		Amount:    order.OrderAmount,
		DenomIn:   order.OrderAmount.Denom,
		DenomOut:  order.OrderTargetDenom,
		MinAmount: sdk.NewCoin(order.OrderTargetDenom, sdkmath.ZeroInt()),
		MaxAmount: order.OrderAmount,
	})
	if err != nil {
		return res, err
	}

	// emit the event
	ctx.EventManager().EmitEvent(types.NewExecuteMarketBuySpotOrderEvt(order, res))

	return res, nil
}
