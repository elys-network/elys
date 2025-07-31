package keeper

import (
	"encoding/binary"
	"math"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	perpetualtypes "github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

// GetPendingPerpetualOrderCount get the total number of pendingPerpetualOrder
func (k Keeper) GetPendingPerpetualOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte{})
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
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte{})
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

	k.SetPendingPerpetualOrder(ctx, pendingPerpetualOrder)

	// Update pendingPerpetualOrder count
	k.SetPendingPerpetualOrderCount(ctx, count+1)

	return count
}

// SetPendingPerpetualOrder set a specific pendingPerpetualOrder in the store
func (k Keeper) SetPendingPerpetualOrder(ctx sdk.Context, pendingPerpetualOrder types.PerpetualOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	b := k.cdc.MustMarshal(&pendingPerpetualOrder)
	store.Set(types.GetPendingPerpetualOrderKeyBytes(sdk.MustAccAddressFromBech32(pendingPerpetualOrder.OwnerAddress), pendingPerpetualOrder.PoolId, pendingPerpetualOrder.OrderId), b)
}

// GetPendingPerpetualOrder returns a pendingPerpetualOrder from its id
func (k Keeper) GetPendingPerpetualOrder(ctx sdk.Context, user sdk.AccAddress, poolId uint64, orderId uint64) (val types.PerpetualOrder, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	b := store.Get(types.GetPendingPerpetualOrderKeyBytes(user, poolId, orderId))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DeletePendingPerpetualOrdersForAddressAndPool deletes all pending perpetual orders for a given address, pool id and position id
func (k Keeper) DeletePendingPerpetualOrdersForAddressAndPool(ctx sdk.Context, user sdk.AccAddress, poolId uint64, positionId uint64) ([]types.PerpetualOrder, error) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetPendingPerpetualOrderAddressPoolKey(user, poolId))
	defer iterator.Close()

	var orders []types.PerpetualOrder
	for ; iterator.Valid(); iterator.Next() {
		var order types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		if order.Status == types.Status_PENDING && order.PositionId == positionId &&
			(order.PerpetualOrderType == types.PerpetualOrderType_LIMITCLOSE || order.PerpetualOrderType == types.PerpetualOrderType_STOPLOSSPERP) {
			store.Delete(types.GetPendingPerpetualOrderKeyBytes(user, poolId, order.OrderId))
			ctx.EventManager().EmitEvent(types.NewDeletePendingPerpetualOrderEvt(order))
		}
	}

	return orders, nil
}

func (k Keeper) GetPendingPerpetualOrdersForAddress(ctx sdk.Context, address string, status *types.Status, pagination *query.PageRequest) ([]types.PerpetualOrder, *query.PageResponse, error) {
	var orders []types.PerpetualOrder

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPendingPerpetualOrderAddressKey(sdk.MustAccAddressFromBech32(address))
	orderStore := prefix.NewStore(store, key)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.FilteredPaginate(orderStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var order types.PerpetualOrder
		err := k.cdc.Unmarshal(value, &order)
		if err == nil {
			if accumulate && (*status == types.Status_ALL || order.Status == *status) {
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
func (k Keeper) RemovePendingPerpetualOrder(ctx sdk.Context, user sdk.AccAddress, poolId uint64, orderId uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	store.Delete(types.GetPendingPerpetualOrderKeyBytes(user, poolId, orderId))
}

// LegacyRemovePendingPerpetualOrder removes a pendingPerpetualOrder from the store
func (k Keeper) LegacyRemovePendingPerpetualOrder(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	store.Delete(GetPendingPerpetualOrderIDBytes(id))
}

// GetAllPendingPerpetualOrder returns all pendingPerpetualOrder
func (k Keeper) GetAllPendingPerpetualOrder(ctx sdk.Context) (list []types.PerpetualOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAllLegacyPendingPerpetualOrder(ctx sdk.Context) (list []types.LegacyPerpetualOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyPerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// DeleteAllPendingPerpetualOrder returns all pendingPerpetualOrder
func (k Keeper) DeleteAllPendingPerpetualOrder(ctx sdk.Context) (list []types.PerpetualOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	return
}

// Remove after migration
// GetPendingPerpetualOrderIDBytes returns the byte representation of the ID
func GetPendingPerpetualOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

func (k Keeper) MigratePendingOrders(ctx sdk.Context) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.SortedPerpetualOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var order types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		k.LegacyRemovePendingPerpetualOrder(ctx, order.OrderId)

		k.SetPendingPerpetualOrder(ctx, order)
	}
}

// ExecuteLimitOpenOrder executes a limit open order
func (k Keeper) ExecuteLimitOpenOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	tradingAsset, err := k.perpetual.GetTradingAsset(ctx, order.PoolId)
	if err != nil {
		return err
	}

	marketPrice, _, err := k.perpetual.GetAssetPriceAndAssetUsdcDenomRatio(ctx, tradingAsset)
	if err != nil {
		return err
	}

	switch order.Position {
	case types.PerpetualPosition_LONG:
		if marketPrice.GT(order.TriggerPrice) {
			// skip the order
			return nil
		}
	case types.PerpetualPosition_SHORT:
		if marketPrice.LT(order.TriggerPrice) {
			// skip the order
			return nil
		}
	}

	// send the collateral amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err = k.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.Collateral))
	if err != nil {
		return err
	}

	openMsg := perpetualtypes.MsgOpen{
		Creator:         order.OwnerAddress,
		Position:        perpetualtypes.Position(order.Position),
		Leverage:        order.Leverage,
		Collateral:      order.Collateral,
		TakeProfitPrice: order.TakeProfitPrice,
		StopLossPrice:   order.StopLossPrice,
		PoolId:          order.PoolId,
	}

	if err = openMsg.ValidateBasic(); err != nil {
		return err
	}
	res, err := k.perpetual.Open(ctx, &openMsg)
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, ownerAddress, order.PoolId, order.OrderId)

	ctx.EventManager().EmitEvent(types.NewExecuteLimitOpenPerpetualOrderEvt(order, res.Id))

	return nil
}

// ExecuteLimitCloseOrder executes a limit close order
func (k Keeper) ExecuteLimitCloseOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	tradingAsset, err := k.perpetual.GetTradingAsset(ctx, order.PoolId)
	if err != nil {
		return err
	}

	marketPrice, _, err := k.perpetual.GetAssetPriceAndAssetUsdcDenomRatio(ctx, tradingAsset)
	if err != nil {
		return err
	}

	switch order.Position {
	case types.PerpetualPosition_LONG:
		if marketPrice.LT(order.TriggerPrice) {
			// skip the order
			return nil
		}
	case types.PerpetualPosition_SHORT:
		if marketPrice.GT(order.TriggerPrice) {
			// skip the order
			return nil
		}
	}

	// get the position info
	position, err := k.perpetual.GetMTP(ctx, order.PoolId, sdk.MustAccAddressFromBech32(order.OwnerAddress), order.PositionId)
	if err != nil {
		return err
	}

	closePercentage := sdkmath.LegacyNewDec(int64(order.ClosePercentage)).Quo(sdkmath.LegacyNewDec(100))
	closeAmount := closePercentage.Mul(sdkmath.LegacyNewDecFromInt(position.Custody)).TruncateInt()

	closeMsg := perpetualtypes.MsgClose{
		Creator: order.OwnerAddress,
		Id:      order.PositionId,
		Amount:  closeAmount,
		PoolId:  order.PoolId,
	}

	if err = closeMsg.ValidateBasic(); err != nil {
		return err
	}

	_, err = k.perpetual.Close(ctx, &closeMsg)
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, sdk.MustAccAddressFromBech32(order.OwnerAddress), order.PoolId, order.OrderId)

	// emit event for limit close order executed
	ctx.EventManager().EmitEvent(types.NewExecuteLimitClosePerpetualOrderEvt(order, closeAmount.String()))

	return nil
}

// ExecuteMarketOpenOrder executes a market open order
func (k Keeper) ExecuteMarketOpenOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	openMsg := perpetualtypes.MsgOpen{
		Creator:         order.OwnerAddress,
		Position:        perpetualtypes.Position(order.Position),
		Leverage:        order.Leverage,
		Collateral:      order.Collateral,
		TakeProfitPrice: order.TakeProfitPrice,
		StopLossPrice:   order.StopLossPrice,
		PoolId:          order.PoolId,
	}
	if err := openMsg.ValidateBasic(); err != nil {
		return err
	}
	_, err := k.perpetual.Open(ctx, &openMsg)
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, sdk.MustAccAddressFromBech32(order.OwnerAddress), order.PoolId, order.OrderId)

	return nil
}

// ExecuteMarketCloseOrder executes a market close order
func (k Keeper) ExecuteMarketCloseOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	closeMsg := perpetualtypes.MsgClose{
		Creator: order.OwnerAddress,
		Id:      order.PositionId,
		Amount:  sdkmath.ZeroInt(),
		PoolId:  order.PoolId,
	}

	if err := closeMsg.ValidateBasic(); err != nil {
		return err
	}

	_, err := k.perpetual.Close(ctx, &closeMsg)
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, sdk.MustAccAddressFromBech32(order.OwnerAddress), order.PoolId, order.OrderId)

	return nil
}

// ConstructPerpetualOrderExtraInfo fills up the extra information of the perpetual order and returns it
func (k Keeper) ConstructPerpetualOrderExtraInfo(ctx sdk.Context, order types.PerpetualOrder) (*types.PerpetualOrderExtraInfo, error) {
	// If position id not set then estimate the info values
	if order.PositionId == 0 {
		res, err := k.perpetual.HandleOpenEstimation(ctx, &perpetualtypes.QueryOpenEstimationRequest{
			Position:        perpetualtypes.Position(order.Position),
			Leverage:        order.Leverage,
			Collateral:      order.Collateral,
			TakeProfitPrice: order.TakeProfitPrice,
			PoolId:          order.PoolId,
			LimitPrice:      order.TriggerPrice,
		})

		// If error use zero values
		if err != nil {
			return &types.PerpetualOrderExtraInfo{
				PerpetualOrder:     &order,
				PositionSize:       sdk.Coin{},
				LiquidationPrice:   sdkmath.LegacyZeroDec(),
				FundingRate:        sdkmath.LegacyZeroDec(),
				BorrowInterestRate: sdkmath.LegacyZeroDec(),
			}, nil
		}

		return &types.PerpetualOrderExtraInfo{
			PerpetualOrder:     &order,
			PositionSize:       res.PositionSize,
			LiquidationPrice:   res.LiquidationPrice,
			FundingRate:        res.FundingRate,
			BorrowInterestRate: res.BorrowInterestRate,
		}, nil
	}

	// otherwise retrieve the position info from existing position
	mtp, err := k.perpetual.GetMTP(ctx, order.PoolId, sdk.AccAddress(order.OwnerAddress), order.PositionId)
	if err != nil {
		return nil, err
	}

	pool, found := k.perpetual.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, perpetualtypes.ErrPoolDoesNotExist
	}

	res, err := k.perpetual.HandleCloseEstimation(ctx, &perpetualtypes.QueryCloseEstimationRequest{
		Address:    order.OwnerAddress,
		PositionId: order.PositionId,
		PoolId:     mtp.AmmPoolId,
	})
	if err != nil {
		return nil, err
	}

	return &types.PerpetualOrderExtraInfo{
		PerpetualOrder:     &order,
		PositionSize:       res.PositionSize,
		LiquidationPrice:   res.LiquidationPrice,
		FundingRate:        pool.FundingRate,
		BorrowInterestRate: pool.BorrowInterestRate,
	}, nil
}
