package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	atypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

const (
	secondsPerYear = 31536000
)

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, index uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolKey(index)
	store.Delete(key)
}

func (k Keeper) GetBaseCurreny(ctx sdk.Context) (atypes.Entry, bool) {
	baseCurrency, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return atypes.Entry{}, false
	}
	return baseCurrency, true
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolKey(poolId)
	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolKey(pool.AmmPoolId)
	b := k.cdc.MustMarshal(&pool)
	store.Set(key, b)
}

// GetAllPools returns all pool
func (k Keeper) GetAllPools(ctx sdk.Context) (list []types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAllLegacyPools(ctx sdk.Context) (list []types.LegacyPool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetBorrowRate(ctx sdk.Context, block uint64, pool uint64, interest types.InterestBlock) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InterestRatePrefix)
	prev := types.GetInterestRateKey(block-1, pool)
	key := types.GetInterestRateKey(block, pool)
	if store.Has(prev) {
		lastBlock := types.InterestBlock{}
		bz := store.Get(prev)
		k.cdc.MustUnmarshal(bz, &lastBlock)
		interest.InterestRate = interest.InterestRate.Add(lastBlock.InterestRate)

		bz = k.cdc.MustMarshal(&interest)
		store.Set(key, bz)
	} else {
		bz := k.cdc.MustMarshal(&interest)
		store.Set(key, bz)
	}
}

// Test it out
// Deletes all pool blocks at delBlock
func (k Keeper) DeleteBorrowRate(ctx sdk.Context, delBlock uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.InterestRatePrefix, types.GetUint64Bytes(delBlock)...))
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetAllBorrowRate(ctx sdk.Context) []types.InterestBlock {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InterestRatePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	interests := []types.InterestBlock{}
	for ; iterator.Valid(); iterator.Next() {
		interest := types.InterestBlock{}
		k.cdc.MustUnmarshal(iterator.Value(), &interest)

		interests = append(interests, interest)
	}
	return interests
}

func (k Keeper) GetBorrowInterestRate(ctx sdk.Context, startBlock, startTime uint64, poolId uint64, takeProfitBorrowFactor math.LegacyDec) math.LegacyDec {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InterestRatePrefix)
	currentBlockKey := types.GetInterestRateKey(uint64(ctx.BlockHeight()), poolId)
	startBlockKey := types.GetInterestRateKey(startBlock, poolId)

	// note: exclude start block
	if store.Has(startBlockKey) && store.Has(currentBlockKey) && startBlock != uint64(ctx.BlockHeight()) {
		bz := store.Get(startBlockKey)
		startInterestBlock := types.InterestBlock{}
		k.cdc.MustUnmarshal(bz, &startInterestBlock)

		bz = store.Get(currentBlockKey)
		endInterestBlock := types.InterestBlock{}
		k.cdc.MustUnmarshal(bz, &endInterestBlock)

		totalInterestRate := endInterestBlock.InterestRate.Sub(startInterestBlock.InterestRate)
		numberOfBlocks := ctx.BlockHeight() - int64(startBlock)

		finalInterestRate := totalInterestRate.
			MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
			QuoInt64(numberOfBlocks).
			QuoInt64(secondsPerYear)

		return math.LegacyMaxDec(
			finalInterestRate.Mul(takeProfitBorrowFactor),
			k.GetParams(ctx).BorrowInterestRateMin.MulInt64(ctx.BlockTime().Unix()-int64(startTime)).QuoInt64(secondsPerYear),
		)
	}

	if !store.Has(startBlockKey) && store.Has(currentBlockKey) {
		iterator := storetypes.KVStorePrefixIterator(store, nil)
		defer iterator.Close()

		firstStoredBlock := uint64(0)
		if iterator.Valid() {
			interestBlock := types.InterestBlock{}
			firstStoredBlock = sdk.BigEndianToUint64(iterator.Key())
			k.cdc.MustUnmarshal(iterator.Value(), &interestBlock)
		}
		if firstStoredBlock > startBlock {
			bz := store.Get(currentBlockKey)
			endInterestBlock := types.InterestBlock{}
			k.cdc.MustUnmarshal(bz, &endInterestBlock)

			numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1

			totalInterestRate := endInterestBlock.InterestRate
			finalInterestRate := totalInterestRate.
				MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
				QuoInt64(numberOfBlocks).
				QuoInt64(secondsPerYear)

			return math.LegacyMaxDec(
				finalInterestRate.Mul(takeProfitBorrowFactor),
				k.GetParams(ctx).BorrowInterestRateMin.MulInt64(ctx.BlockTime().Unix()-int64(startTime)).QuoInt64(secondsPerYear),
			)
		}
	}
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		// this is handling case of future block
		return math.LegacyZeroDec()
	}
	newInterest := pool.BorrowInterestRate.MulInt64(ctx.BlockTime().Unix() - int64(startTime)).QuoInt64(secondsPerYear)
	return newInterest
}

func (k Keeper) SetFundingRate(ctx sdk.Context, block uint64, pool uint64, funding types.FundingRateBlock) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.FundingRatePrefix)
	prev := types.GetFundingRateKey(block-1, pool)
	key := types.GetFundingRateKey(block, pool)
	if store.Has(prev) {
		lastBlock := types.FundingRateBlock{}
		bz := store.Get(prev)
		k.cdc.MustUnmarshal(bz, &lastBlock)
		funding.FundingRateLong = funding.FundingRateLong.Add(lastBlock.FundingRateLong)
		funding.FundingRateShort = funding.FundingRateShort.Add(lastBlock.FundingRateShort)

		funding.FundingShareLong = funding.FundingShareLong.Add(lastBlock.FundingShareLong)
		funding.FundingShareShort = funding.FundingShareShort.Add(lastBlock.FundingShareShort)

		bz = k.cdc.MustMarshal(&funding)
		store.Set(key, bz)
	} else {
		bz := k.cdc.MustMarshal(&funding)
		store.Set(key, bz)
	}
}

// Test it out
// Deletes all pool blocks at delBlock
func (k Keeper) DeleteFundingRate(ctx sdk.Context, delBlock uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), append(types.FundingRatePrefix, types.GetUint64Bytes(delBlock)...))
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetAllFundingRate(ctx sdk.Context) []types.FundingRateBlock {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.FundingRatePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	fundings := []types.FundingRateBlock{}
	for ; iterator.Valid(); iterator.Next() {
		funding := types.FundingRateBlock{}
		k.cdc.MustUnmarshal(iterator.Value(), &funding)

		fundings = append(fundings, funding)
	}
	return fundings
}

func (k Keeper) GetFundingRate(ctx sdk.Context, startBlock uint64, startTime uint64, poolId uint64) (long math.LegacyDec, short math.LegacyDec) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.FundingRatePrefix)
	currentBlockKey := types.GetFundingRateKey(uint64(ctx.BlockHeight()), poolId)
	startBlockKey := types.GetFundingRateKey(startBlock, poolId)

	// note: exclude start block
	if store.Has(startBlockKey) && store.Has(currentBlockKey) && startBlock != uint64(ctx.BlockHeight()) {
		bz := store.Get(startBlockKey)
		startFundingBlock := types.FundingRateBlock{}
		k.cdc.MustUnmarshal(bz, &startFundingBlock)

		bz = store.Get(currentBlockKey)
		endFundingBlock := types.FundingRateBlock{}
		k.cdc.MustUnmarshal(bz, &endFundingBlock)

		numberOfBlocks := ctx.BlockHeight() - int64(startBlock)

		totalFundingLong := endFundingBlock.FundingRateLong.Sub(startFundingBlock.FundingRateLong).
			MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
			QuoInt64(numberOfBlocks).
			QuoInt64(secondsPerYear)
		totalFundingShort := endFundingBlock.FundingRateShort.Sub(startFundingBlock.FundingRateShort).
			MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
			QuoInt64(numberOfBlocks).
			QuoInt64(secondsPerYear)
		return totalFundingLong, totalFundingShort
	}

	if !store.Has(startBlockKey) && store.Has(currentBlockKey) {
		iterator := storetypes.KVStorePrefixIterator(store, nil)
		defer iterator.Close()

		firstStoredBlock := uint64(0)
		if iterator.Valid() {
			fundingBlock := types.FundingRateBlock{}
			firstStoredBlock = sdk.BigEndianToUint64(iterator.Key())
			k.cdc.MustUnmarshal(iterator.Value(), &fundingBlock)
		}
		if firstStoredBlock > startBlock {
			bz := store.Get(currentBlockKey)
			endFundingBlock := types.FundingRateBlock{}
			numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1
			k.cdc.MustUnmarshal(bz, &endFundingBlock)

			return endFundingBlock.FundingRateLong.MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
					QuoInt64(numberOfBlocks).
					QuoInt64(secondsPerYear),
				endFundingBlock.FundingRateShort.MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
					QuoInt64(numberOfBlocks).
					QuoInt64(secondsPerYear)
		}
	}
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		panic("perpetual pool not found")
	}

	if pool.BorrowInterestRate.IsPositive() {
		return pool.FundingRate.MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
			QuoInt64(secondsPerYear), math.LegacyZeroDec()
	} else {
		return math.LegacyZeroDec(), pool.FundingRate.MulInt64(ctx.BlockTime().Unix() - int64(startTime)).
			QuoInt64(secondsPerYear)
	}
}

// Deletes all pool blocks at delBlock
func (k Keeper) DeleteAllFundingRate(ctx sdk.Context) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.FundingRatePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetFundingDistributionValue(ctx sdk.Context, startBlock uint64, pool uint64) (long math.LegacyDec, short math.LegacyDec) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.FundingRatePrefix)
	currentBlockKey := types.GetFundingRateKey(uint64(ctx.BlockHeight()), pool)
	startBlockKey := types.GetFundingRateKey(startBlock, pool)

	// note: exclude start block
	if store.Has(startBlockKey) && store.Has(currentBlockKey) && startBlock != uint64(ctx.BlockHeight()) {
		bz := store.Get(startBlockKey)
		startFundingBlock := types.FundingRateBlock{}
		k.cdc.MustUnmarshal(bz, &startFundingBlock)

		bz = store.Get(currentBlockKey)
		endFundingBlock := types.FundingRateBlock{}
		k.cdc.MustUnmarshal(bz, &endFundingBlock)

		totalCustodyLong := endFundingBlock.FundingShareLong.Sub(startFundingBlock.FundingShareLong)
		totalCustodyShort := endFundingBlock.FundingShareShort.Sub(startFundingBlock.FundingShareShort)

		return totalCustodyLong, totalCustodyShort
	}

	if !store.Has(startBlockKey) && store.Has(currentBlockKey) {
		iterator := storetypes.KVStorePrefixIterator(store, nil)
		defer iterator.Close()

		firstStoredBlock := uint64(0)
		if iterator.Valid() {
			fundingBlock := types.FundingRateBlock{}
			firstStoredBlock = sdk.BigEndianToUint64(iterator.Key())
			k.cdc.MustUnmarshal(iterator.Value(), &fundingBlock)
		}
		if firstStoredBlock > startBlock {
			bz := store.Get(currentBlockKey)
			endFundingBlock := types.FundingRateBlock{}
			k.cdc.MustUnmarshal(bz, &endFundingBlock)

			totalCustodyLong := endFundingBlock.FundingShareLong
			totalCustodyShort := endFundingBlock.FundingShareShort
			return totalCustodyLong, totalCustodyShort
		}
	}

	return math.LegacyZeroDec(), math.LegacyZeroDec()
}

func (k Keeper) GetTradingAsset(ctx sdk.Context, poolId uint64) (string, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return "", errors.New("pool not found")
	}
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return "", errorsmod.Wrapf(atypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	tradingAsset, err := pool.GetTradingAsset(baseCurrency)
	if err != nil {
		return "", err
	}
	return tradingAsset, nil
}
