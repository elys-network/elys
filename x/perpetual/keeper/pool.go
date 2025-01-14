package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	"github.com/elys-network/elys/x/perpetual/types"
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

// GetAllLegacyPools returns all legacy pool
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

	blocksPerYear := math.LegacyNewDec(int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear))

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
			Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(math.LegacyNewDec(numberOfBlocks)).
			Quo(blocksPerYear)

		return math.LegacyMaxDec(
			finalInterestRate.Mul(takeProfitBorrowFactor),
			k.GetParams(ctx).BorrowInterestRateMin.Mul(math.LegacyNewDec(ctx.BlockTime().Unix()-int64(startTime))).Quo(blocksPerYear),
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
				Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
				Quo(math.LegacyNewDec(numberOfBlocks)).
				Quo(blocksPerYear)

			return math.LegacyMaxDec(
				finalInterestRate.Mul(takeProfitBorrowFactor),
				k.GetParams(ctx).BorrowInterestRateMin.Mul(math.LegacyNewDec(ctx.BlockTime().Unix()-int64(startTime))).Quo(blocksPerYear),
			)
		}
	}
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		// this is handling case of future block
		return math.LegacyZeroDec()
	}
	newInterest := pool.BorrowInterestRate.Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).Quo(blocksPerYear)
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

		funding.FundingAmountLong = funding.FundingAmountLong.Add(lastBlock.FundingAmountLong)
		funding.FundingAmountShort = funding.FundingAmountShort.Add(lastBlock.FundingAmountShort)

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

	blocksPerYear := math.LegacyNewDec(int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear))

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
			Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(math.LegacyNewDec(numberOfBlocks)).
			Quo(blocksPerYear)
		totalFundingShort := endFundingBlock.FundingRateShort.Sub(startFundingBlock.FundingRateShort).
			Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(math.LegacyNewDec(numberOfBlocks)).
			Quo(blocksPerYear)
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

			return endFundingBlock.FundingRateLong.Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
					Quo(math.LegacyNewDec(numberOfBlocks)).
					Quo(blocksPerYear),
				endFundingBlock.FundingRateShort.Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
					Quo(math.LegacyNewDec(numberOfBlocks)).
					Quo(blocksPerYear)
		}
	}
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		panic("perpetual pool not found")
	}

	if pool.BorrowInterestRate.IsPositive() {
		return pool.FundingRate.Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(blocksPerYear), math.LegacyZeroDec()
	} else {
		return math.LegacyZeroDec(), pool.FundingRate.Mul(math.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(blocksPerYear)
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

		totalCustodyLong := endFundingBlock.FundingAmountLong.Sub(startFundingBlock.FundingAmountLong)
		totalCustodyShort := endFundingBlock.FundingAmountShort.Sub(startFundingBlock.FundingAmountShort)

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

			totalCustodyLong := endFundingBlock.FundingAmountLong
			totalCustodyShort := endFundingBlock.FundingAmountShort
			return totalCustodyLong, totalCustodyShort
		}
	}

	return math.LegacyZeroDec(), math.LegacyZeroDec()
}

// Deletes all pool blocks at delBlock
func (k Keeper) DeleteAllInterestRate(ctx sdk.Context) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InterestRatePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
