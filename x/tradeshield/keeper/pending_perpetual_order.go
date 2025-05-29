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
	perpetualtypes "github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/elys-network/elys/v5/x/tradeshield/types"
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

	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	appendedValue := k.cdc.MustMarshal(&pendingPerpetualOrder)
	store.Set(GetPendingPerpetualOrderIDBytes(pendingPerpetualOrder.OrderId), appendedValue)

	// Update pendingPerpetualOrder count
	k.SetPendingPerpetualOrderCount(ctx, count+1)

	return count
}

// SetPendingPerpetualOrder set a specific pendingPerpetualOrder in the store
func (k Keeper) SetPendingPerpetualOrder(ctx sdk.Context, pendingPerpetualOrder types.PerpetualOrder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	b := k.cdc.MustMarshal(&pendingPerpetualOrder)
	store.Set(GetPendingPerpetualOrderIDBytes(pendingPerpetualOrder.OrderId), b)
}

// GetPendingPerpetualOrder returns a pendingPerpetualOrder from its id
func (k Keeper) GetPendingPerpetualOrder(ctx sdk.Context, id uint64) (val types.PerpetualOrder, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	b := store.Get(GetPendingPerpetualOrderIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetPendingPerpetualOrdersForAddress(ctx sdk.Context, address string, status *types.Status, pagination *query.PageRequest) ([]types.PerpetualOrder, *query.PageResponse, error) {
	var orders []types.PerpetualOrder

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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

// RemovePendingPerpetualOrder removes a pendingPerpetualOrder from the store
func (k Keeper) RemovePendingPerpetualOrder(ctx sdk.Context, id uint64) {
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

// SetAllLegacyPerpetualTriggerPriceToNewTriggerPriceStructure set all legacy perpetual trigger price to new trigger price structure
func (k Keeper) SetAllLegacyPerpetualTriggerPriceToNewTriggerPriceStructure(ctx sdk.Context) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PendingPerpetualOrderKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var order types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		order.TriggerPrice = order.LegacyTriggerPriceV1.Rate
		order.LegacyTriggerPriceV1 = types.LegacyTriggerPriceV1{}
		store.Set(iterator.Key(), k.cdc.MustMarshal(&order))
	}
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

func (k Keeper) GetAllSortedPerpetualOrder(ctx sdk.Context) (list [][]uint64, err error) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.SortedPerpetualOrderKey)
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

// ExecuteLimitOpenOrder executes a limit open order
func (k Keeper) ExecuteLimitOpenOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	marketPrice, _, err := k.perpetual.GetAssetPriceAndAssetUsdcDenomRatio(ctx, order.TradingAsset)
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

	res, err := k.perpetual.Open(ctx, &perpetualtypes.MsgOpen{
		Creator:         order.OwnerAddress,
		Position:        perpetualtypes.Position(order.Position),
		Leverage:        order.Leverage,
		TradingAsset:    order.TradingAsset,
		Collateral:      order.Collateral,
		TakeProfitPrice: order.TakeProfitPrice,
		StopLossPrice:   order.StopLossPrice,
		PoolId:          order.PoolId,
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, order.OrderId)

	ctx.EventManager().EmitEvent(types.NewExecuteLimitOpenPerpetualOrderEvt(order, res.Id))

	return nil
}

// ExecuteLimitCloseOrder executes a limit close order
func (k Keeper) ExecuteLimitCloseOrder(ctx sdk.Context, order types.PerpetualOrder) error {
	marketPrice, _, err := k.perpetual.GetAssetPriceAndAssetUsdcDenomRatio(ctx, order.TradingAsset)
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

	_, err = k.perpetual.Close(ctx, &perpetualtypes.MsgClose{
		Creator: order.OwnerAddress,
		Id:      order.PositionId,
		Amount:  sdkmath.ZeroInt(),
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, order.OrderId)

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
		PoolId:          order.PoolId,
	})
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
		Amount:  sdkmath.ZeroInt(),
	})
	if err != nil {
		return err
	}

	// Remove the order from the pending order list
	k.RemovePendingPerpetualOrder(ctx, order.OrderId)

	return nil
}

// ConstructPerpetualOrderExtraInfo fills up the extra information of the perpetual order and returns it
func (k Keeper) ConstructPerpetualOrderExtraInfo(ctx sdk.Context, order types.PerpetualOrder) (*types.PerpetualOrderExtraInfo, error) {
	// If position id not set then estimate the info values
	if order.PositionId == 0 {
		res, err := k.perpetual.HandleOpenEstimation(ctx, &perpetualtypes.QueryOpenEstimationRequest{
			Position:        perpetualtypes.Position(order.Position),
			Leverage:        order.Leverage,
			TradingAsset:    order.TradingAsset,
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
	mtp, err := k.perpetual.GetMTP(ctx, sdk.AccAddress(order.OwnerAddress), order.PositionId)
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
