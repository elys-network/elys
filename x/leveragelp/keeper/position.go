package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	cosmosMath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPosition(ctx sdk.Context, positionAddress sdk.AccAddress, id uint64) (types.Position, error) {
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

// remove after migration
func (k Keeper) DeleteLegacyPosition(ctx sdk.Context, positionAddress string, id uint64) {
	key := types.GetLegacyPositionKey(positionAddress, id)
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) SetPosition(ctx sdk.Context, position *types.Position) {
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)

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

func (k Keeper) SetOffset(ctx sdk.Context, offset uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OffsetKeyPrefix, types.GetUint64Bytes(offset))
}

func (k Keeper) GetOffset(ctx sdk.Context) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.OffsetKeyPrefix) {
		res := store.Get(types.OffsetKeyPrefix)
		return types.GetUint64FromBytes(res), true
	} else {
		return 0, false
	}
}

func (k Keeper) DeleteOffset(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.OffsetKeyPrefix)
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

func (k Keeper) GetPositions(ctx sdk.Context, pagination *query.PageRequest) ([]*types.Position, *query.PageResponse, error) {
	var positionList []*types.Position
	store := ctx.KVStore(k.storeKey)
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
		debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())
		position.Liabilities = debt.GetTotalLiablities()
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

func (k Keeper) GetPositionsForAddress(ctx sdk.Context, positionAddress sdk.AccAddress, pagination *query.PageRequest) ([]*types.PositionAndInterest, *query.PageResponse, error) {
	var positions []*types.PositionAndInterest

	store := ctx.KVStore(k.storeKey)
	positionStore := prefix.NewStore(store, types.GetPositionPrefixForAddress(positionAddress))

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
		debt := k.stableKeeper.GetDebt(ctx, positionAndInterest.Position.GetPositionAddress())
		positionAndInterest.Position.Liabilities = debt.GetTotalLiablities()
		positions = append(positions, &positionAndInterest)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return positions, pageRes, nil
}

func (k Keeper) GetPositionHealth(ctx sdk.Context, position types.Position) (sdk.Dec, error) {
	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress())
	debtAmount := debt.GetTotalLiablities()
	if debtAmount.IsZero() {
		return sdk.ZeroDec(), nil
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return sdk.Dec{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	leveragedLpAmount := sdk.ZeroInt()
	commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress())

	for _, commitment := range commitments.CommittedTokens {
		leveragedLpAmount = leveragedLpAmount.Add(commitment.Amount)
	}

	exitCoinsAfterFee, _, err := k.amm.ExitPoolEst(ctx, position.GetAmmPoolId(), leveragedLpAmount, baseCurrency)
	if err != nil {
		return sdk.Dec{}, err
	}

	exitAmountAfterFee := exitCoinsAfterFee.AmountOf(baseCurrency)

	health := exitAmountAfterFee.ToLegacyDec().Quo(debtAmount.ToLegacyDec())

	return health, nil
}

func (k Keeper) GetPositionWithId(ctx sdk.Context, positionAddress sdk.AccAddress, Id uint64) (*types.Position, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPositionKey(positionAddress, Id)
	if !store.Has(key) {
		return nil, false
	}
	res := store.Get(key)
	var position types.Position
	k.cdc.MustUnmarshal(res, &position)
	return &position, true
}

// TODO: remove all functions below after upgrade
func (k Keeper) DeleteCorruptedKeys(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PositionPrefix)
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
		if err != nil || position.Id == 0 {
			store.Delete(iterator.Key())
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
