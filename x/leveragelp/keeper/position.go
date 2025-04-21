package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPosition(ctx sdk.Context, positionAddress sdk.AccAddress, id uint64) (types.Position, error) {
	var position types.Position
	key := types.GetPositionKey(positionAddress, id)
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
	count := k.GetPositionCount(ctx)
	openCount := k.GetOpenPositionCount(ctx)
	creator := sdk.MustAccAddressFromBech32(position.Address)

	if position.Id == 0 {
		// increment global id count
		count++
		position.Id = count
		k.SetPositionCount(ctx, count)
		// increment open position count
		openCount++
		k.SetOpenPositionCount(ctx, openCount)
	}

	key := types.GetPositionKey(creator, position.Id)
	store.Set(key, k.cdc.MustMarshal(position))
}

func (k Keeper) DestroyPosition(ctx sdk.Context, positionAddress sdk.AccAddress, id uint64) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

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

// Set Open Position count
func (k Keeper) SetOpenPositionCount(ctx sdk.Context, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.OpenPositionCountPrefix, types.GetUint64Bytes(count))
}

// Set Position count
func (k Keeper) SetPositionCount(ctx sdk.Context, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.PositionCountPrefix, types.GetUint64Bytes(count))
}

func (k Keeper) GetPositionCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.PositionCountPrefix)
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

func (k Keeper) GetOpenPositionCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.OpenPositionCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) GetPositionIterator(ctx sdk.Context) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return storetypes.KVStorePrefixIterator(store, types.PositionPrefix)
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []types.Position {
	var positions []types.Position
	iterator := k.GetPositionIterator(ctx)
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

func (k Keeper) GetPositions(ctx sdk.Context, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positionList []*types.Position
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
		positionList = append(positionList, &position)
		return nil
	})

	return positionList, pageRes, err
}

func (k Keeper) GetPositionsForPool(ctx sdk.Context, ammPoolId uint64, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positions []*types.Position

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

func (k Keeper) GetPositionsForAddress(ctx sdk.Context, positionAddress sdk.AccAddress, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positions []*types.Position

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.GetPositionPrefixForAddress(positionAddress))

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
		debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
		position.Liabilities = debt.GetTotalLiablities()
		positions = append(positions, &position)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return positions, pageRes, nil
}

// GetPositionHealth Should not be used in queries as UpdateInterestAndGetDebt updates KVStore as well
func (k Keeper) GetPositionHealth(ctx sdk.Context, position types.Position) (sdkmath.LegacyDec, error) {
	if position.LeveragedLpAmount.IsZero() {
		return sdkmath.LegacyZeroDec(), nil
	}
	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
	debtAmount := debt.GetTotalLiablities()
	if debtAmount.IsZero() {
		return sdkmath.LegacyMaxSortableDec, nil
	}

	debtDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)
	debtValue := debtAmount.ToLegacyDec().Mul(debtDenomPrice)

	ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	ammTVL, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	positionValue := position.LeveragedLpAmount.ToLegacyDec().Mul(ammTVL).Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	health := positionValue.Quo(debtValue)

	return health, nil
}

func (k Keeper) GetPositionWithId(ctx sdk.Context, positionAddress sdk.AccAddress, Id uint64) (*types.Position, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPositionKey(positionAddress, Id)
	if !store.Has(key) {
		return nil, false
	}
	res := store.Get(key)
	var position types.Position
	k.cdc.MustUnmarshal(res, &position)
	return &position, true
}

func (k Keeper) MigrateData(ctx sdk.Context) {
	iterator := k.GetPositionIterator(ctx)
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
			leveragedLpAmount := sdkmath.ZeroInt()
			commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress())

			for _, commitment := range commitments.CommittedTokens {
				leveragedLpAmount = leveragedLpAmount.Add(commitment.Amount)
			}
			pool, found := k.GetPool(ctx, position.AmmPoolId)
			if found {
				pool.LeveragedLpAmount = pool.LeveragedLpAmount.Add(leveragedLpAmount)
				pool.Health = k.CalculatePoolHealth(ctx, &pool)
				k.SetPool(ctx, pool)
			}

			// Repay any balance, delete position
			debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
			repayAmount := debt.GetTotalLiablities()

			// Check if position has enough coins to repay else repay partial
			bal := k.bankKeeper.GetBalance(ctx, position.GetPositionAddress(), position.Collateral.Denom)
			userAmount := sdkmath.ZeroInt()
			if bal.Amount.LT(repayAmount) {
				repayAmount = bal.Amount
			} else {
				userAmount = bal.Amount.Sub(repayAmount)
			}

			if repayAmount.IsPositive() {
				k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, repayAmount), position.BorrowPoolId, position.AmmPoolId)
			} else {
				userAmount = bal.Amount
			}

			positionOwner := sdk.MustAccAddressFromBech32(position.Address)
			if userAmount.IsPositive() {
				k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), positionOwner, sdk.Coins{sdk.NewCoin(position.Collateral.Denom, userAmount)})
			}

			if leveragedLpAmount.IsZero() {
				// Repay any balance, delete position
				k.DestroyPosition(ctx, positionOwner, position.Id)
			} else {
				// Repay any balance and update position value
				position.LeveragedLpAmount = leveragedLpAmount
				k.SetPosition(ctx, &position)
			}
		}
	}
}

func (k Keeper) SetAllPositions(ctx sdk.Context) {
	iterator := k.GetPositionIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		leveragedLpAmount := sdkmath.ZeroInt()
		commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress())

		poolDenom := ammtypes.GetPoolShareDenom(position.AmmPoolId)
		for _, commitment := range commitments.CommittedTokens {
			if poolDenom == commitment.Denom {
				leveragedLpAmount = leveragedLpAmount.Add(commitment.Amount)
			}
		}

		// Set correct lev amount
		position.LeveragedLpAmount = leveragedLpAmount
		position.BorrowPoolId = stabletypes.UsdcPoolId
		k.SetPosition(ctx, &position)
	}

	// Pool liabilities are reset in stablestake migration
	k.V18MigratonPoolLiabilities(ctx)

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		// After setting pool liabilities
		balance := k.bankKeeper.GetBalance(ctx, position.GetPositionAddress(), position.Collateral.Denom)
		if balance.IsPositive() {
			debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)
			totalLiab := debt.GetTotalLiablities()
			if totalLiab.GT(balance.Amount) {
				totalLiab = balance.Amount
			}
			if totalLiab.IsPositive() {
				k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, totalLiab), position.AmmPoolId, position.BorrowPoolId)
			}
			if balance.Amount.GT(totalLiab) {
				payToUser := balance.Amount.Sub(totalLiab)
				k.bankKeeper.SendCoins(ctx, position.GetPositionAddress(), sdk.MustAccAddressFromBech32(position.Address), sdk.Coins{sdk.NewCoin(position.Collateral.Denom, payToUser)})
			}
		}

		if position.LeveragedLpAmount.IsZero() {
			k.DestroyPosition(ctx, sdk.MustAccAddressFromBech32(position.Address), position.Id)
		}
	}
	return
}

func (k Keeper) V18MigratonPoolLiabilities(ctx sdk.Context) {
	iterator := k.GetPositionIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		debt := k.stableKeeper.GetDebtWithoutUpdatedInterest(ctx, position.GetPositionAddress(), stabletypes.UsdcPoolId)
		k.stableKeeper.AddPoolLiabilities(ctx, position.AmmPoolId, sdk.NewCoin(position.Collateral.Denom, debt.GetTotalLiablities()))
		k.SetPosition(ctx, &position)
	}
	return
}
