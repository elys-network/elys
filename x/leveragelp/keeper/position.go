package keeper

import (
	"fmt"
	"math"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (k Keeper) SetPosition(ctx sdk.Context, position *types.Position) {
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
			debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, old.GetPositionAddress())
			liquidationKey := types.GetLiquidationSortKey(old.AmmPoolId, old.LeveragedLpAmount, debt.Borrowed, old.Id)
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

	// Add position sort keys
	addrId := types.AddressId{
		Id:      position.Id,
		Address: position.Address,
	}
	bz := k.cdc.MustMarshal(&addrId)
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())
	liquidationKey := types.GetLiquidationSortKey(position.AmmPoolId, position.LeveragedLpAmount, debt.Borrowed, position.Id)
	if len(liquidationKey) > 0 {
		store.Set(liquidationKey, bz)
	}
	stopLossKey := types.GetStopLossSortKey(position.AmmPoolId, position.StopLossPrice, position.Id)
	if len(stopLossKey) > 0 {
		store.Set(stopLossKey, bz)
	}
}

func (k Keeper) DestroyPosition(ctx sdk.Context, positionAddress string, id uint64) error {
	key := types.GetPositionKey(positionAddress, id)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.ErrPositionDoesNotExist
	}
	store.Delete(key)

	// Remove position sort keys
	old, err := k.GetPosition(ctx, positionAddress, id)
	if err == nil {
		debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, old.GetPositionAddress())
		liquidationKey := types.GetLiquidationSortKey(old.AmmPoolId, old.LeveragedLpAmount, debt.Borrowed, old.Id)
		if len(liquidationKey) > 0 {
			store.Delete(liquidationKey)
		}
		stopLossKey := types.GetStopLossSortKey(old.AmmPoolId, old.StopLossPrice, old.Id)
		if len(stopLossKey) > 0 {
			store.Delete(stopLossKey)
		}
	}

	// decrement open position count
	openCount := k.GetOpenPositionCount(ctx)
	if openCount != 0 {
		openCount--
	}

	// Set open Position count
	k.SetOpenPositionCount(ctx, openCount)

	return nil
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
		k.cdc.MustUnmarshal(bytesValue, &position)
		positions = append(positions, position)
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
		k.cdc.MustUnmarshal(value, &position)
		positionList = append(positionList, &position)
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
		k.cdc.MustUnmarshal(value, &position)
		if accumulate && position.AmmPoolId == ammPoolId {
			positions = append(positions, &position)
			return true, nil
		}

		return false, nil
	})

	return positions, pageRes, err
}

func (k Keeper) GetPositionsForAddress(ctx sdk.Context, positionAddress sdk.Address, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positions []*types.Position

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

	pageRes, err := query.Paginate(positionStore, pagination, func(key []byte, value []byte) error {
		var position types.Position
		k.cdc.MustUnmarshal(value, &position)
		positions = append(positions, &position)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return positions, pageRes, nil
}

func (k Keeper) GetPositionHealth(ctx sdk.Context, position types.Position, ammPool ammtypes.Pool) (sdk.Dec, error) {
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())
	xl := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}

	commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress().String())
	positionVal := sdk.ZeroDec()
	depositDenom := k.stableKeeper.GetDepositDenom(ctx)
	for _, commitment := range commitments.CommittedTokens {
		cacheCtx, _ := ctx.CacheContext()
		cacheCtx = cacheCtx.WithBlockTime(cacheCtx.BlockTime().Add(time.Hour))
		exitCoins, err := k.amm.ExitPool(cacheCtx, position.GetPositionAddress(), ammPool.PoolId, commitment.Amount, sdk.Coins{}, depositDenom)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		positionVal = positionVal.Add(sdk.NewDecFromInt(exitCoins.AmountOf(depositDenom)))
	}

	lr := positionVal.Quo(sdk.NewDecFromBigInt(xl.BigInt()))
	return lr, nil
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
