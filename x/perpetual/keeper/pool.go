package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, index uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(index))
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllPool returns all pool
func (k Keeper) GetAllPools(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetBorrowRate(ctx sdk.Context, block uint64, pool uint64, interest types.InterestBlock) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestRatePrefix)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.InterestRatePrefix, types.GetUint64Bytes(delBlock)...))
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetAllBorrowRate(ctx sdk.Context) []types.InterestBlock {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestRatePrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	interests := []types.InterestBlock{}
	for ; iterator.Valid(); iterator.Next() {
		interest := types.InterestBlock{}
		k.cdc.MustUnmarshal(iterator.Value(), &interest)

		interests = append(interests, interest)
	}
	return interests
}

func (k Keeper) GetBorrowRate(ctx sdk.Context, startBlock uint64, pool uint64, startTime uint64, borrowed sdk.Dec) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestRatePrefix)
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
		numberOfBlocks := ctx.BlockHeight() - int64(startBlock)

		newInterest := borrowed.
			Mul(totalInterest).
			Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(sdk.NewDec(numberOfBlocks)).
			Quo(sdk.NewDec(86400 * 365)).
			RoundInt()
		return newInterest
	}

	if !store.Has(startBlockKey) && store.Has(currentBlockKey) {
		iterator := sdk.KVStorePrefixIterator(store, nil)
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
			numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1

			newInterest := borrowed.Mul(totalInterest).
				Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
				Quo(sdk.NewDec(numberOfBlocks)).
				Quo(sdk.NewDec(86400 * 365)).
				RoundInt()
			return newInterest
		}
	}
	params, found := k.GetPool(ctx, pool)
	if !found {
		return sdk.ZeroInt()
	}
	newInterest := borrowed.
		Mul(params.BorrowInterestRate).
		Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
		Quo(sdk.NewDec(86400 * 365)).
		RoundInt()
	return newInterest
}

func (k Keeper) SetFundingRate(ctx sdk.Context, block uint64, pool uint64, funding types.FundingRateBlock) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingRatePrefix)
	prev := types.GetFundingRateKey(block-1, pool)
	key := types.GetFundingRateKey(block, pool)
	if store.Has(prev) {
		lastBlock := types.FundingRateBlock{}
		bz := store.Get(prev)
		k.cdc.MustUnmarshal(bz, &lastBlock)
		funding.FundingRate = funding.FundingRate.Add(lastBlock.FundingRate)

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.FundingRatePrefix, types.GetUint64Bytes(delBlock)...))
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) GetAllFundingRate(ctx sdk.Context) []types.FundingRateBlock {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingRatePrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	fundings := []types.FundingRateBlock{}
	for ; iterator.Valid(); iterator.Next() {
		funding := types.FundingRateBlock{}
		k.cdc.MustUnmarshal(iterator.Value(), &funding)

		fundings = append(fundings, funding)
	}
	return fundings
}

func (k Keeper) GetFundingRate(ctx sdk.Context, startBlock uint64, pool uint64, startTime uint64, custody sdk.Dec) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingRatePrefix)
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
		numberOfBlocks := ctx.BlockHeight() - int64(startBlock)

		newFunding := custody.
			Mul(totalFunding).
			Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(sdk.NewDec(numberOfBlocks)).
			Quo(sdk.NewDec(86400 * 365)).
			RoundInt()
		return newFunding
	}

	if !store.Has(startBlockKey) && store.Has(currentBlockKey) {
		iterator := sdk.KVStorePrefixIterator(store, nil)
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

			totalFunding := endFundingBlock.FundingRate
			numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1

			newFunding := custody.Mul(totalFunding).
				Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
				Quo(sdk.NewDec(numberOfBlocks)).
				Quo(sdk.NewDec(86400 * 365)).
				RoundInt()
			return newFunding
		}
	}
	params, found := k.GetPool(ctx, pool)
	if !found {
		return sdk.ZeroInt()
	}
	newFunding := custody.
		Mul(params.BorrowInterestRate).
		Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
		Quo(sdk.NewDec(86400 * 365)).
		RoundInt()
	return newFunding
}
