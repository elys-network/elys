package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetDebtWithoutUpdatedInterestStacked(ctx sdk.Context, addr sdk.AccAddress) types.Debt {
	debt := types.Debt{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DebtPrefixKey)
	bz := store.Get([]byte(addr.String()))
	if len(bz) == 0 {
		return types.Debt{
			Address:               addr.String(),
			Borrowed:              sdk.ZeroInt(),
			InterestPaid:          sdk.ZeroInt(),
			InterestStacked:       sdk.ZeroInt(),
			BorrowTime:            uint64(ctx.BlockTime().Unix()),
			LastInterestCalcTime:  uint64(ctx.BlockTime().Unix()),
			LastInterestCalcBlock: uint64(ctx.BlockHeight()),
		}
	}

	k.cdc.MustUnmarshal(bz, &debt)
	return debt
}

func (k Keeper) GetDebtWithUpdatedInterestStacked(ctx sdk.Context, addr sdk.AccAddress) types.Debt {
	debt := k.GetDebtWithoutUpdatedInterestStacked(ctx, addr)
	debt = k.UpdateInterestStacked(ctx, debt)
	return debt
}

func (k Keeper) SetInterest(ctx sdk.Context, block uint64, interest types.InterestBlock) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestPrefixKey)
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

func (k Keeper) DeleteInterest(ctx sdk.Context, delBlock int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestPrefixKey)
	key := sdk.Uint64ToBigEndian(uint64(delBlock))
	if store.Has(key) {
		store.Delete([]byte(key))
	}
}

func (k Keeper) GetAllInterest(ctx sdk.Context) []types.InterestBlock {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestPrefixKey)
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

func (k Keeper) GetInterest(ctx sdk.Context, startBlock uint64, startTime uint64, borrowed sdk.Dec) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestPrefixKey)
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
	params := k.GetParams(ctx)
	newInterest := borrowed.
		Mul(params.InterestRate).
		Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
		Quo(sdk.NewDec(86400 * 365)).
		RoundInt()
	return newInterest
}

func (k Keeper) SetDebt(ctx sdk.Context, debt types.Debt) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DebtPrefixKey)
	bz := k.cdc.MustMarshal(&debt)
	store.Set([]byte(debt.Address), bz)
}

func (k Keeper) DeleteDebt(ctx sdk.Context, debt types.Debt) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DebtPrefixKey)
	store.Delete([]byte(debt.Address))
}

func (k Keeper) AllDebts(ctx sdk.Context) []types.Debt {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DebtPrefixKey)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	debts := []types.Debt{}
	for ; iterator.Valid(); iterator.Next() {
		debt := types.Debt{}
		k.cdc.MustUnmarshal(iterator.Value(), &debt)

		debts = append(debts, debt)
	}
	return debts
}

func (k Keeper) UpdateInterestStacked(ctx sdk.Context, debt types.Debt) types.Debt {
	params := k.GetParams(ctx)
	newInterest := k.GetInterest(ctx, debt.LastInterestCalcBlock, debt.LastInterestCalcTime, debt.Borrowed.ToLegacyDec())

	debt.InterestStacked = debt.InterestStacked.Add(newInterest)
	debt.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
	debt.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	k.SetDebt(ctx, debt)

	params.TotalValue = params.TotalValue.Add(newInterest)
	k.SetParams(ctx, params)
	return debt
}

func (k Keeper) Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	depositDenom := k.GetDepositDenom(ctx)
	if depositDenom != amount.Denom {
		return types.ErrInvalidBorrowDenom
	}

	// For security reasons, we should avoid borrowing more than 90% in total to the stablestake pool.
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	params := k.GetParams(ctx)
	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)

	borrowed := params.TotalValue.Sub(balance.Amount).ToLegacyDec().Add(amount.Amount.ToLegacyDec())
	maxAllowed := params.TotalValue.ToLegacyDec().Mul(sdk.NewDec(9)).Quo(sdk.NewDec(10))
	if borrowed.GT(maxAllowed) {
		return types.ErrMaxBorrowAmount
	}

	debt := k.GetDebtWithUpdatedInterestStacked(ctx, addr)
	debt.Borrowed = debt.Borrowed.Add(amount.Amount)
	k.SetDebt(ctx, debt)
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{amount})
}

func (k Keeper) Repay(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	depositDenom := k.GetDepositDenom(ctx)
	if depositDenom != amount.Denom {
		return types.ErrInvalidBorrowDenom
	}

	err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}

	// calculate latest interest stacked
	debt := k.GetDebtWithUpdatedInterestStacked(ctx, addr)

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
