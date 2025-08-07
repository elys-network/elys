package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPosition(ctx sdk.Context, poolId uint64, positionAddress sdk.AccAddress, id uint64) (types.Position, error) {
	var position types.Position
	key := types.GetPositionKey(poolId, positionAddress, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(key)
	if bz == nil {
		return position, types.ErrPositionDoesNotExist
	}
	k.cdc.MustUnmarshal(bz, &position)
	k.UpdatePositionLiabilties(ctx, &position)
	return position, nil
}

func (k Keeper) UpdatePositionLiabilties(ctx sdk.Context, position *types.Position) {
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
	position.Liabilities = debt.GetTotalLiablities()
}

func (k Keeper) SetPosition(ctx sdk.Context, position *types.Position) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	creator := sdk.MustAccAddressFromBech32(position.Address)

	if position.Id == 0 {
		// Use atomic increment to prevent race conditions
		position.Id = k.IncrementPositionCounter(ctx, position.AmmPoolId)
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

func (k Keeper) GetAllPositions(ctx sdk.Context) []types.Position {
	var positions []types.Position
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PositionPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &position)
		if err == nil {
			k.UpdatePositionLiabilties(ctx, &position)
			positions = append(positions, position)
		}
	}
	return positions
}

func (k Keeper) GetAllPositionsForPool(ctx sdk.Context, poolId uint64) []types.Position {
	var positions []types.Position
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.GetPoolPrefixKey(poolId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &position)
		if err == nil {
			k.UpdatePositionLiabilties(ctx, &position)
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
		k.UpdatePositionLiabilties(ctx, &position)
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

	pageRes, err := query.Paginate(positionStore, pagination, func(key []byte, value []byte) error {
		var position types.Position
		err := k.cdc.Unmarshal(value, &position)
		if err != nil {
			return err
		}
		k.UpdatePositionLiabilties(ctx, &position)
		positions = append(positions, position)
		return nil
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
		k.UpdatePositionLiabilties(ctx, &position)
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
			k.UpdatePositionLiabilties(ctx, &position)
			positionList = append(positionList, position)
		}
		iterator.Close()
	}

	return positionList
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
	k.UpdatePositionLiabilties(ctx, &position)
	return &position, true
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

	// Check for division by zero
	if ammPool.TotalShares.Amount.IsZero() {
		return osmomath.ZeroBigDec(), fmt.Errorf("amm pool %d has zero total shares", position.AmmPoolId)
	}
	positionValue := position.GetBigDecLeveragedLpAmount().Mul(ammTVL).Quo(osmomath.BigDecFromSDKInt(ammPool.TotalShares.Amount))

	health := positionValue.Quo(debtValue)

	return health, nil
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

func (k Keeper) DeleteLegacyPosition(ctx sdk.Context, position types.Position) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	key := types.GetLegacyPositionKey(position.GetOwnerAddress(), position.Id)
	store.Delete(key)
}

func (k Keeper) MigratePositionsToNewKeys(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.LegacyPositionPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)

		k.SetPosition(ctx, &position)

		k.DeleteLegacyPosition(ctx, position)

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

	store.Delete(types.LegacyPositionCountPrefix)
	store.Delete(types.LegacyOpenPositionCountPrefix)
}
