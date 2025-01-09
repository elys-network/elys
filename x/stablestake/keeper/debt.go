package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) getDebt(ctx sdk.Context, addr sdk.AccAddress, poolId uint64) (debt types.Debt) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetDebtKey(addr)
	bz := store.Get(key)
	if len(bz) == 0 {
		return types.Debt{
			Address:               addr.String(),
			Borrowed:              sdkmath.ZeroInt(),
			InterestPaid:          sdkmath.ZeroInt(),
			InterestStacked:       sdkmath.ZeroInt(),
			BorrowTime:            uint64(ctx.BlockTime().Unix()),
			LastInterestCalcTime:  uint64(ctx.BlockTime().Unix()),
			LastInterestCalcBlock: uint64(ctx.BlockHeight()),
			PoolId:                poolId,
		}
	}

	k.cdc.MustUnmarshal(bz, &debt)
	return
}

func (k Keeper) GetDebt(ctx sdk.Context, addr sdk.AccAddress, poolId uint64) types.Debt {
	debt := k.getDebt(ctx, addr, poolId)
	debt.InterestStacked = debt.InterestStacked.Add(k.GetInterestForPool(ctx, debt.LastInterestCalcBlock, debt.LastInterestCalcTime, debt.Borrowed.ToLegacyDec(), debt.PoolId))
	debt.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
	debt.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	return debt
}

func (k Keeper) UpdateInterestAndGetDebt(ctx sdk.Context, addr sdk.AccAddress, poolId uint64) types.Debt {
	debt := k.getDebt(ctx, addr, poolId)
	debt = k.UpdateInterestStacked(ctx, debt)
	return debt
}

func (k Keeper) SetDebt(ctx sdk.Context, debt types.Debt) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetDebtKey(debt.GetOwnerAccount())
	bz := k.cdc.MustMarshal(&debt)
	store.Set(key, bz)
}

func (k Keeper) DeleteDebt(ctx sdk.Context, debt types.Debt) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetDebtKey(debt.GetOwnerAccount())
	store.Delete(key)
}

func (k Keeper) GetAllDebts(ctx sdk.Context) []types.Debt {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, types.DebtPrefixKey)
	defer iterator.Close()

	debts := []types.Debt{}
	for ; iterator.Valid(); iterator.Next() {
		debt := types.Debt{}
		k.cdc.MustUnmarshal(iterator.Value(), &debt)

		debts = append(debts, debt)
	}
	return debts
}

func (k Keeper) SetInterestForPool(ctx sdk.Context, poolId uint64, block uint64, interest types.InterestBlock) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetInterestKey(poolId))
	if store.Has(sdk.Uint64ToBigEndian(block - 1)) {
		lastBlock := types.InterestBlock{}
		bz := store.Get(sdk.Uint64ToBigEndian(block - 1))
		k.cdc.MustUnmarshal(bz, &lastBlock)
		interest.InterestRate = interest.InterestRate.Add(lastBlock.InterestRate)

		bz = k.cdc.MustMarshal(&interest)
		store.Set(sdk.Uint64ToBigEndian(block), bz)
	} else {
		bz := k.cdc.MustMarshal(&interest)
		store.Set(sdk.Uint64ToBigEndian(block), bz)
	}
}

func (k Keeper) DeleteInterestForPool(ctx sdk.Context, delBlock int64, poolId uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetInterestKey(poolId))
	key := sdk.Uint64ToBigEndian(uint64(delBlock))
	if store.Has(key) {
		store.Delete([]byte(key))
	}
}

func (k Keeper) GetAllInterestForPool(ctx sdk.Context, poolId uint64) []types.InterestBlock {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetInterestKey(poolId))
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	interests := []types.InterestBlock{}
	for ; iterator.Valid(); iterator.Next() {
		interest := types.InterestBlock{}
		k.cdc.MustUnmarshal(iterator.Value(), &interest)

		// FIXME: remove this in the next upgrade
		if interest.BlockHeight == 0 {
			block := iterator.Key()
			interest.BlockHeight = sdk.BigEndianToUint64(block)
		}

		interests = append(interests, interest)
	}
	return interests
}

func (k Keeper) GetAllInterest(ctx sdk.Context) []types.InterestBlock {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InterestPrefixKey)
	iterator := storetypes.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	interests := []types.InterestBlock{}
	for ; iterator.Valid(); iterator.Next() {
		interest := types.InterestBlock{}
		k.cdc.MustUnmarshal(iterator.Value(), &interest)

		// FIXME: remove this in the next upgrade
		if interest.BlockHeight == 0 {
			block := iterator.Key()
			interest.BlockHeight = sdk.BigEndianToUint64(block)
		}

		interests = append(interests, interest)
	}
	return interests
}

func (k Keeper) GetInterestForPool(ctx sdk.Context, startBlock uint64, startTime uint64, borrowed sdkmath.LegacyDec, poolId uint64) sdkmath.Int {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetInterestKey(poolId))
	currentBlockKey := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	startBlockKey := sdk.Uint64ToBigEndian(startBlock)

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

		newInterest := borrowed.
			Mul(totalInterestRate).
			Mul(sdkmath.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(sdkmath.LegacyNewDec(numberOfBlocks)).
			Quo(sdkmath.LegacyNewDec(86400 * 365)).
			RoundInt()
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
			numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1

			newInterest := borrowed.Mul(totalInterest).
				Mul(sdkmath.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
				Quo(sdkmath.LegacyNewDec(numberOfBlocks)).
				Quo(sdkmath.LegacyNewDec(86400 * 365)).
				RoundInt()
			return newInterest
		}
	}
	params := k.GetParams(ctx)
	newInterest := borrowed.
		Mul(params.InterestRate).
		Mul(sdkmath.LegacyNewDec(ctx.BlockTime().Unix() - int64(startTime))).
		Quo(sdkmath.LegacyNewDec(86400 * 365)).
		RoundInt()
	return newInterest
}

func (k Keeper) UpdateInterestStacked(ctx sdk.Context, debt types.Debt) types.Debt {
	params := k.GetParams(ctx)
	newInterest := k.GetInterestForPool(ctx, debt.LastInterestCalcBlock, debt.LastInterestCalcTime, debt.Borrowed.ToLegacyDec(), debt.PoolId)

	debt.InterestStacked = debt.InterestStacked.Add(newInterest)
	debt.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
	debt.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	k.SetDebt(ctx, debt)

	params.TotalValue = params.TotalValue.Add(newInterest)
	k.SetParams(ctx, params)
	return debt
}

func (k Keeper) Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.ErrPoolNotFound
	}
	depositDenom := pool.GetDepositDenom()
	if depositDenom != amount.Denom {
		return types.ErrInvalidBorrowDenom
	}

	// For security reasons, we should avoid borrowing more than 90% in total to the stablestake pool.
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	params := k.GetParams(ctx)
	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)

	borrowed := params.TotalValue.Sub(balance.Amount).ToLegacyDec().Add(amount.Amount.ToLegacyDec())
	maxAllowed := params.TotalValue.ToLegacyDec().Mul(params.MaxLeverageRatio)
	if borrowed.GT(maxAllowed) {
		return types.ErrMaxBorrowAmount
	}

	debt := k.UpdateInterestAndGetDebt(ctx, addr, poolId)
	debt.Borrowed = debt.Borrowed.Add(amount.Amount)
	k.SetDebt(ctx, debt)
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{amount})
}

func (k Keeper) Repay(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.ErrPoolNotFound
	}
	depositDenom := pool.GetDepositDenom()
	if depositDenom != amount.Denom {
		return types.ErrInvalidBorrowDenom
	}

	err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}

	// calculate latest interest stacked
	debt := k.UpdateInterestAndGetDebt(ctx, addr, poolId)

	// repay interest
	interestPayAmount := debt.InterestStacked.Sub(debt.InterestPaid)
	if interestPayAmount.GT(amount.Amount) {
		interestPayAmount = amount.Amount
	}

	// repay borrowed
	repayAmount := amount.Amount.Sub(interestPayAmount)
	debt.Borrowed = debt.Borrowed.Sub(repayAmount)
	debt.InterestPaid = debt.InterestPaid.Add(interestPayAmount)

	if debt.Borrowed.IsNegative() {
		return types.ErrNegativeBorrowed
	}

	if debt.Borrowed.IsZero() {
		k.DeleteDebt(ctx, debt)
	} else {
		k.SetDebt(ctx, debt)
	}
	return nil
}
