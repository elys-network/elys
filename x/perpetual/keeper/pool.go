package keeper

import (
	sdkmath "cosmossdk.io/math"
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

// GetAllPool returns all pool
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

func (k Keeper) GetAllLegacyPools(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.LegacyPoolKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) RemoveLegacyPool(ctx sdk.Context, index uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetLegacyPoolKey(index))
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

func (k Keeper) GetBorrowRate(ctx sdk.Context, startBlock uint64, pool uint64, borrowed sdkmath.LegacyDec) sdkmath.LegacyDec {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InterestRatePrefix)
	currentBlockKey := types.GetInterestRateKey(uint64(ctx.BlockHeight()), pool)
	startBlockKey := types.GetInterestRateKey(startBlock, pool)

	// note: exclude start block
	if store.Has(startBlockKey) && store.Has(currentBlockKey) && startBlock != uint64(ctx.BlockHeight()) {
		bz := store.Get(startBlockKey)
		startInterestBlock := types.InterestBlock{}
		k.cdc.MustUnmarshal(bz, &startInterestBlock)

		bz = store.Get(currentBlockKey)
		endInterestBlock := types.InterestBlock{}
		k.cdc.MustUnmarshal(bz, &endInterestBlock)

		totalInterest := endInterestBlock.InterestRate.Sub(startInterestBlock.InterestRate)

		newInterest := borrowed.Mul(totalInterest)

		return newInterest
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

			totalInterest := endInterestBlock.InterestRate
			newInterest := borrowed.Mul(totalInterest)
			return newInterest
		}
	}
	params, found := k.GetPool(ctx, pool)
	if !found {
		return sdkmath.LegacyZeroDec()
	}
	newInterest := borrowed.Mul(params.BorrowInterestRate)
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
		funding.FundingRate = funding.FundingRate.Add(lastBlock.FundingRate)
		funding.FundingAmountLong = funding.FundingAmountLong.Add(lastBlock.FundingAmountLong)
		funding.FundingAmountShort = funding.FundingAmountShort.Add(lastBlock.FundingAmountShort)

		funding.FundingRateLong = funding.FundingRateLong.Add(lastBlock.FundingRateLong)
		funding.FundingRateShort = funding.FundingRateShort.Add(lastBlock.FundingRateShort)

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

func (k Keeper) GetFundingRate(ctx sdk.Context, startBlock uint64, pool uint64) (net sdkmath.LegacyDec, long sdkmath.LegacyDec, short sdkmath.LegacyDec) {
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

		totalFunding := endFundingBlock.FundingRate.Sub(startFundingBlock.FundingRate)
		totalFundingLong := endFundingBlock.FundingRateLong.Sub(startFundingBlock.FundingRateLong)
		totalFundingShort := endFundingBlock.FundingRateShort.Sub(startFundingBlock.FundingRateShort)
		return totalFunding, totalFundingLong, totalFundingShort
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

			return endFundingBlock.FundingRate, endFundingBlock.FundingRateLong, endFundingBlock.FundingRateShort
		}
	}
	params, found := k.GetPool(ctx, pool)
	if !found {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()
	}

	if params.BorrowInterestRate.IsPositive() {
		return params.FundingRate, params.FundingRate, sdkmath.LegacyZeroDec()
	} else {
		return params.FundingRate, sdkmath.LegacyZeroDec(), params.FundingRate
	}
}

func (k Keeper) GetFundingDistributionValue(ctx sdk.Context, startBlock uint64, pool uint64) (long sdkmath.LegacyDec, short sdkmath.LegacyDec) {
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
		numberOfBlocks := ctx.BlockHeight() - int64(startBlock)

		totalCustodyShort := endFundingBlock.FundingAmountShort.Sub(startFundingBlock.FundingAmountShort)

		return totalCustodyLong.ToLegacyDec().Quo(sdkmath.LegacyNewDec(numberOfBlocks)), totalCustodyShort.ToLegacyDec().Quo(sdkmath.LegacyNewDec(numberOfBlocks))
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

			numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1
			totalCustodyLong := endFundingBlock.FundingAmountLong
			totalCustodyShort := endFundingBlock.FundingAmountShort
			return totalCustodyLong.ToLegacyDec().Quo(sdkmath.LegacyNewDec(numberOfBlocks)), totalCustodyShort.ToLegacyDec().Quo(sdkmath.LegacyNewDec(numberOfBlocks))
		}
	}

	return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()
}
