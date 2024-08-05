package keeper

import (
	errorsmod "cosmossdk.io/errors"
	cosmosMath "cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
)

func (k Keeper) GetPosition(ctx sdk.Context, positionAddress string, id uint64) (types.Position, error) {
	var position types.Position
	key := types.GetPositionKey(positionAddress, id)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return position, types.ErrPositionDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &position)
	return position, nil
}

func (k Keeper) SetPosition(ctx sdk.Context, position *types.Position, oldDebt sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	count := k.GetPositionCount(ctx)
	openCount := k.GetOpenPositionCount(ctx)

	if position.Id == 0 {
		// increment global id count
		count++
		position.Id = count
		k.SetPositionCount(ctx, count)
		// increment open position count
		openCount++
		k.SetOpenPositionCount(ctx, openCount)
	} else {
		old, err := k.GetPosition(ctx, position.Address, position.Id)
		if err == nil {
			// Make sure liability changes are handled properly here, this should always be updated whenever liability is changed
			liquidationKey := types.GetLiquidationSortKey(old.AmmPoolId, old.LeveragedLpAmount, oldDebt, old.Id)
			if len(liquidationKey) > 0 {
				store.Delete(liquidationKey)
			}
			stopLossKey := types.GetStopLossSortKey(old.AmmPoolId, old.StopLossPrice, old.Id)
			if len(stopLossKey) > 0 {
				store.Delete(stopLossKey)
			}
		}
	}

	key := types.GetPositionKey(position.Address, position.Id)
	store.Set(key, k.cdc.MustMarshal(position))

	// for stablestake hook
	store.Set([]byte(position.GetPositionAddress()), key)

	// Add position sort keys
	addrId := types.AddressId{
		Id:      position.Id,
		Address: position.Address,
	}
	bz := k.cdc.MustMarshal(&addrId)
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())
	liquidationKey := types.GetLiquidationSortKey(position.AmmPoolId, position.LeveragedLpAmount, debt.Borrowed.Sub(debt.InterestPaid).Add(debt.InterestStacked), position.Id)
	if len(liquidationKey) > 0 {
		store.Set(liquidationKey, bz)
	}
	stopLossKey := types.GetStopLossSortKey(position.AmmPoolId, position.StopLossPrice, position.Id)
	if len(stopLossKey) > 0 {
		store.Set(stopLossKey, bz)
	}
}

func (k Keeper) DestroyPosition(ctx sdk.Context, positionAddress string, id uint64, oldDebt sdk.Int) error {
	store := ctx.KVStore(k.storeKey)

	// Remove position sort keys
	old, err := k.GetPosition(ctx, positionAddress, id)
	if err == nil {
		debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, old.GetPositionAddress())
		liquidationKey := types.GetLiquidationSortKey(old.AmmPoolId, old.LeveragedLpAmount, debt.Borrowed.Sub(debt.InterestPaid).Add(debt.InterestStacked), old.Id)
		if len(liquidationKey) > 0 {
			store.Delete(liquidationKey)
		}
		stopLossKey := types.GetStopLossSortKey(old.AmmPoolId, old.StopLossPrice, old.Id)
		if len(stopLossKey) > 0 {
			store.Delete(stopLossKey)
		}
		store.Delete(old.GetPositionAddress())
	}

	key := types.GetPositionKey(positionAddress, id)
	if !store.Has(key) {
		return types.ErrPositionDoesNotExist
	}
	store.Delete(key)

	// decrement open position count
	openCount := k.GetOpenPositionCount(ctx)
	if openCount != 0 {
		openCount--
	}

	// Set open Position count
	k.SetOpenPositionCount(ctx, openCount)

	return nil
}

// Set sorted liquidation
func (k Keeper) SetSortedLiquidation(ctx sdk.Context, address string, old sdk.Int, new sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(address)) {
		key := store.Get([]byte(address))
		if !store.Has(key) {
			return
		}
		res := store.Get(key)
		var position types.Position
		k.cdc.MustUnmarshal(res, &position)
		// Make sure liability changes are handled properly here, this should always be updated whenever liability is changed
		liquidationKey := types.GetLiquidationSortKey(position.AmmPoolId, position.LeveragedLpAmount, old, position.Id)
		if len(liquidationKey) > 0 && store.Has(liquidationKey) {
			store.Delete(liquidationKey)
		}

		// Add position sort keys
		addrId := types.AddressId{
			Id:      position.Id,
			Address: position.Address,
		}
		bz := k.cdc.MustMarshal(&addrId)
		liquidationKey = types.GetLiquidationSortKey(position.AmmPoolId, position.LeveragedLpAmount, new, position.Id)
		if len(liquidationKey) > 0 {
			store.Set(liquidationKey, bz)
		}
	}
}

// Change DS for migration
func (k Keeper) SetSortedLiquidationAndStopLoss(ctx sdk.Context, position types.Position) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPositionKey(position.Address, position.Id)
	// for stablestake hook
	store.Set([]byte(position.GetPositionAddress()), key)

	// Add position sort keys
	addrId := types.AddressId{
		Id:      position.Id,
		Address: position.Address,
	}
	bz := k.cdc.MustMarshal(&addrId)
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())
	liquidationKey := types.GetLiquidationSortKey(position.AmmPoolId, position.LeveragedLpAmount, debt.Borrowed.Sub(debt.InterestPaid).Add(debt.InterestStacked), position.Id)
	if len(liquidationKey) > 0 {
		store.Set(liquidationKey, bz)
	}
	stopLossKey := types.GetStopLossSortKey(position.AmmPoolId, position.StopLossPrice, position.Id)
	if len(stopLossKey) > 0 {
		store.Set(stopLossKey, bz)
	}
}

// Set Open Position count
func (k Keeper) SetOpenPositionCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OpenPositionCountPrefix, types.GetUint64Bytes(count))
}

// Set Position count
func (k Keeper) SetPositionCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PositionCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetPositionCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := ctx.KVStore(k.storeKey).Get(types.PositionCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) GetOpenPositionCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := ctx.KVStore(k.storeKey).Get(types.OpenPositionCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) GetPositionIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PositionPrefix)
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []types.Position {
	var positions []types.Position
	iterator := k.GetPositionIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &position)
		if err == nil {
			positions = append(positions, position)
		}
	}
	return positions
}

func (k Keeper) IteratePoolPosIdsLiquidationSorted(ctx sdk.Context, poolId uint64, fn func(posId types.AddressId) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidationSortPrefix(poolId))
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		addrId := types.AddressId{}
		k.cdc.MustUnmarshal(iterator.Value(), &addrId)
		stop := fn(addrId)
		if stop {
			return
		}
	}
}

func (k Keeper) IteratePoolPosIdsStopLossSorted(ctx sdk.Context, poolId uint64, fn func(posId types.AddressId) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetStopLossSortPrefix(poolId))
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		addrId := types.AddressId{}
		k.cdc.MustUnmarshal(iterator.Value(), &addrId)
		stop := fn(addrId)
		if stop {
			return
		}
	}
}

func (k Keeper) DeletePoolPosIdsLiquidationSorted(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetLiquidationSortPrefix(poolId))
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) DeletePoolPosIdsStopLossSorted(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetStopLossSortPrefix(poolId))
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetPositions(ctx sdk.Context, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positionList []*types.Position
	store := ctx.KVStore(k.storeKey)
	positionStore := prefix.NewStore(store, types.PositionPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(positionStore, pagination, func(key []byte, value []byte) error {
		var position types.Position
		err := k.cdc.Unmarshal(value, &position)
		if err == nil {
			positionList = append(positionList, &position)
		}
		return nil
	})

	return positionList, pageRes, err
}

func (k Keeper) GetPositionsForPool(ctx sdk.Context, ammPoolId uint64, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positions []*types.Position

	store := ctx.KVStore(k.storeKey)
	positionStore := prefix.NewStore(store, types.PositionPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}
	pageRes, err := query.FilteredPaginate(positionStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var position types.Position
		err := k.cdc.Unmarshal(value, &position)
		if err == nil {
			if accumulate && position.AmmPoolId == ammPoolId {
				positions = append(positions, &position)
				return true, nil
			}
		}
		return false, nil
	})

	return positions, pageRes, err
}

func (k Keeper) GetPositionsForAddress(ctx sdk.Context, positionAddress sdk.Address, pagination *query.PageRequest) ([]*types.PositionAndInterest, *query.PageResponse, error) {
	var positions []*types.PositionAndInterest

	store := ctx.KVStore(k.storeKey)
	positionStore := prefix.NewStore(store, types.GetPositionPrefixForAddress(positionAddress.String()))

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: types.MaxPageLimit,
		}
	}

	if pagination.Limit > types.MaxPageLimit {
		return nil, nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	params := k.stableKeeper.GetParams(ctx)
	hours := cosmosMath.LegacyNewDec(365 * 24)
	pageRes, err := query.Paginate(positionStore, pagination, func(key []byte, value []byte) error {
		var p types.Position
		k.cdc.MustUnmarshal(value, &p)
		var positionAndInterest types.PositionAndInterest
		positionAndInterest.Position = &p
		price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, p.Collateral.Denom)
		interestRateHour := params.InterestRate.Quo(hours)
		positionAndInterest.InterestRateHour = interestRateHour
		positionAndInterest.InterestRateHourUsd = interestRateHour.Mul(cosmosMath.LegacyDec(p.Liabilities.Mul(price.RoundInt())))
		debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, positionAndInterest.Position.GetPositionAddress())
		positionAndInterest.Position.Liabilities = debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
		positions = append(positions, &positionAndInterest)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return positions, pageRes, nil
}

func (k Keeper) GetPositionHealth(ctx sdk.Context, position types.Position) (sdk.Dec, error) {
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())
	debtAmount := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return sdk.Dec{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	ammPool, err := k.GetAmmPool(ctx, position.GetAmmPoolId())
	if err != nil {
		return sdk.Dec{}, nil
	}

	leveragedLpAmount := sdk.ZeroDec()
	commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress().String())

	for _, commitment := range commitments.CommittedTokens {
		leveragedLpAmount = leveragedLpAmount.Add(commitment.Amount.ToLegacyDec())
	}

	exitCoins, _, err := ammPool.CalcExitPoolCoinsFromShares(ctx, k.oracleKeeper, k.accountedPoolKeeper, leveragedLpAmount.TruncateInt(), baseCurrency)
	if err != nil {
		return sdk.Dec{}, err
	}

	exitFeeCoins := ammkeeper.PortionCoins(exitCoins, ammPool.PoolParams.ExitFee)
	exitAmountAfterFee := exitFeeCoins.AmountOf(baseCurrency)

	health := exitAmountAfterFee.ToLegacyDec().Quo(debtAmount.ToLegacyDec())

	return health, nil
}

func (k Keeper) GetPositionWithId(ctx sdk.Context, positionAddress sdk.Address, Id uint64) (*types.Position, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPositionKey(positionAddress.String(), Id)
	if !store.Has(key) {
		return nil, false
	}
	res := store.Get(key)
	var position types.Position
	k.cdc.MustUnmarshal(res, &position)
	return &position, true
}

// FIXME: currently we only avoid the error while loading and not entirely delete the corrupted key, value
// This is done to ensure we don't delete everything.
// After the upgrade that comes after 0.39.0
// We need to uncomment it and remove any corrupted data with this logic

// func (k Keeper) MigrateKeys(ctx sdk.Context) {
// 	store := ctx.KVStore(k.storeKey)
// 	iterator := sdk.KVStorePrefixIterator(store, types.PositionPrefix)
// 	defer func(iterator sdk.Iterator) {
// 		err := iterator.Close()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}(iterator)

// 	for ; iterator.Valid(); iterator.Next() {
// 		var position types.Position
// 		bytesValue := iterator.Value()
// 		err := k.cdc.Unmarshal(bytesValue, &position)
// 		if err != nil {
// 			store.Delete(iterator.Key())
// 		}
// 	}
// }
