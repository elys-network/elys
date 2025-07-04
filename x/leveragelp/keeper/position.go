package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPosition(ctx sdk.Context, poolId uint64, positionAddress sdk.AccAddress, id uint64) (types.Position, error) {
	var position types.Position
	key := types.GetPositionKey(poolId, positionAddress, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	if !store.Has(key) {
		return position, types.ErrPositionDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &position)
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
	position.Liabilities = debt.GetTotalLiablities()
	return position, nil
}

func (k Keeper) SetPosition(ctx sdk.Context, position *types.Position) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	creator := sdk.MustAccAddressFromBech32(position.Address)

	if position.Id == 0 {
		positionCounter := k.GetPositionCounter(ctx, position.AmmPoolId)
		positionCounter.Counter++
		positionCounter.TotalOpen++
		k.SetPositionCounter(ctx, positionCounter)

		position.Id = positionCounter.Counter
	}

	key := types.GetPositionKey(position.AmmPoolId, creator, position.Id)
	store.Set(key, k.cdc.MustMarshal(position))
}

func (k Keeper) DestroyPosition(ctx sdk.Context, position types.Position) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	key := types.GetPositionKey(position.AmmPoolId, position.GetOwnerAddress(), position.Id)
	if !store.Has(key) {
		return types.ErrPositionDoesNotExist
	}
	store.Delete(key)

	positionCounter := k.GetPositionCounter(ctx, position.AmmPoolId)
	positionCounter.TotalOpen--
	k.SetPositionCounter(ctx, positionCounter)

	return nil
}

// Set Open Position count
func (k Keeper) SetLegacyOpenPositionCount(ctx sdk.Context, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.LegacyOpenPositionCountPrefix, types.GetUint64Bytes(count))
}

// Set Position count
func (k Keeper) SetLegacyPositionCount(ctx sdk.Context, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.LegacyPositionCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetLegacyPositionCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.LegacyPositionCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) SetOffset(ctx sdk.Context, offset uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.OffsetKeyPrefix, types.GetUint64Bytes(offset))
}

func (k Keeper) GetOffset(ctx sdk.Context) (uint64, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	if store.Has(types.OffsetKeyPrefix) {
		res := store.Get(types.OffsetKeyPrefix)
		return types.GetUint64FromBytes(res), true
	} else {
		return 0, false
	}
}

func (k Keeper) DeleteOffset(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.OffsetKeyPrefix)
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []types.Position {
	var positions []types.Position
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PositionPrefix)
	defer func(iterator storetypes.Iterator) {
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
			debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
			position.Liabilities = debt.GetTotalLiablities()
			positions = append(positions, position)
		}
	}
	return positions
}

func (k Keeper) GetPositions(ctx sdk.Context, pagination *query.PageRequest) ([]types.Position, *query.PageResponse, error) {
	var positionList []types.Position
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.PositionPrefix)

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
		err := k.cdc.Unmarshal(value, &position)
		if err != nil {
			return err
		}
		debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
		position.Liabilities = debt.GetTotalLiablities()
		positionList = append(positionList, position)
		return nil
	})

	return positionList, pageRes, err
}

func (k Keeper) GetPositionsForPool(ctx sdk.Context, ammPoolId uint64, pagination *query.PageRequest) ([]types.Position, *query.PageResponse, error) {
	var positions []types.Position

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.GetPoolPrefixKey(ammPoolId))

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: types.MaxPageLimit,
		}
	}

	if pagination.Limit > types.MaxPageLimit {
		return nil, nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	pageRes, err := query.FilteredPaginate(positionStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var position types.Position
		err := k.cdc.Unmarshal(value, &position)
		if err == nil {
			if accumulate && position.AmmPoolId == ammPoolId {
				positions = append(positions, position)
				return true, nil
			}
		}
		return false, nil
	})

	return positions, pageRes, err
}

func (k Keeper) GetPositionsForPoolAndAddress(ctx sdk.Context, poolId uint64, positionAddress sdk.AccAddress) []types.Position {

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.GetPoolCreatorPrefixKey(poolId, positionAddress))
	defer iterator.Close()

	var positionList []types.Position
	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		positionList = append(positionList, position)
	}

	return positionList
}

func (k Keeper) GetPositionsForAddress(ctx sdk.Context, positionAddress sdk.AccAddress) []types.Position {

	allPools := k.GetAllPools(ctx)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	var positionList []types.Position

	for _, pool := range allPools {
		iterator := storetypes.KVStorePrefixIterator(store, types.GetPoolCreatorPrefixKey(pool.AmmPoolId, positionAddress))
		for ; iterator.Valid(); iterator.Next() {
			var position types.Position
			k.cdc.MustUnmarshal(iterator.Value(), &position)
			positionList = append(positionList, position)
		}
		iterator.Close()
	}

	return positionList
}

// GetPositionHealth Should not be used in queries as UpdateInterestAndGetDebt updates KVStore as well
func (k Keeper) GetPositionHealth(ctx sdk.Context, position types.Position) (osmomath.BigDec, error) {
	if position.LeveragedLpAmount.IsZero() {
		return osmomath.ZeroBigDec(), nil
	}
	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
	debtAmount := debt.GetBigDecTotalLiablities()
	if debtAmount.IsZero() {
		maxDec := osmomath.OneBigDec().Quo(osmomath.SmallestBigDec())
		return maxDec, nil
	}

	debtDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, position.Collateral.Denom)
	debtValue := debtAmount.Mul(debtDenomPrice)

	ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
	if err != nil {
		return osmomath.ZeroBigDec(), err
	}

	ammTVL, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return osmomath.ZeroBigDec(), err
	}
	positionValue := position.GetBigDecLeveragedLpAmount().Mul(ammTVL).Quo(osmomath.BigDecFromSDKInt(ammPool.TotalShares.Amount))

	health := positionValue.Quo(debtValue)

	return health, nil
}

func (k Keeper) GetPositionWithId(ctx sdk.Context, poolId uint64, positionAddress sdk.AccAddress, Id uint64) (*types.Position, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPositionKey(poolId, positionAddress, Id)
	if !store.Has(key) {
		return nil, false
	}
	res := store.Get(key)
	var position types.Position
	k.cdc.MustUnmarshal(res, &position)
	return &position, true
}

func (k Keeper) MigrateToNewKeys(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.LegacyPositionPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)

		k.SetPosition(ctx, &position)

		positionCounter := k.GetPositionCounter(ctx, position.AmmPoolId)
		positionCounter.TotalOpen++
		k.SetPositionCounter(ctx, positionCounter)
	}

	count := k.GetLegacyPositionCount(ctx) + 1
	// we don't know whats the latest counter for the pool, we can count though but simplest would be just to set from global value of legacy
	for _, positionCounter := range k.GetAllPositionCounters(ctx) {
		positionCounter.Counter = count
		k.SetPositionCounter(ctx, positionCounter)
	}
}
