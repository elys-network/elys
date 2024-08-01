package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetDebt(ctx sdk.Context, addr sdk.AccAddress) types.Debt {
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

func (k Keeper) UpdateInterestStackedByAddress(ctx sdk.Context, addr sdk.AccAddress) types.Debt {
	debt := k.GetDebt(ctx, addr)
	debt = k.UpdateInterestStacked(ctx, debt)
	k.SetDebt(ctx, debt)
	return debt
}

func (k Keeper) SetInterest(ctx sdk.Context, block uint64, interest types.InterestBlock) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestPrefixKey)
	if store.Has(sdk.Uint64ToBigEndian(block - 1)) {
		res := store.Get(sdk.Uint64ToBigEndian(block - 1))
		prev := sdk.MustNewDecFromStr(string(res))
		interest.InterestRate = interest.InterestRate.Add(prev)
		bz := k.cdc.MustMarshal(&interest)
		store.Set(sdk.Uint64ToBigEndian(block), bz)
	} else {
		bz := k.cdc.MustMarshal(&interest)
		store.Set(sdk.Uint64ToBigEndian(block), bz)
	}
}

func (k Keeper) GetInterest(ctx sdk.Context, startBlock uint64, startTime uint64, borrowed sdk.Dec) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InterestPrefixKey)

	if store.Has(sdk.Uint64ToBigEndian(startBlock-1)) && store.Has(sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))) {
		bz := store.Get(sdk.Uint64ToBigEndian(startBlock - 1))
		startInterestBlock := types.InterestBlock{}
		k.cdc.MustUnmarshal(bz, &startInterestBlock)

		bz = store.Get(sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight())))
		endInterestBlock := types.InterestBlock{}
		k.cdc.MustUnmarshal(bz, &endInterestBlock)

		totalInterest := endInterestBlock.InterestRate.Sub(startInterestBlock.InterestRate)
		numberOfBlocks := ctx.BlockHeight() - int64(startBlock) + 1

		newInterest := borrowed.Mul(totalInterest).
			Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(startTime))).
			Quo(sdk.NewDec(numberOfBlocks)).
			Quo(sdk.NewDec(86400 * 365)).
			RoundInt()
		return newInterest
	}
	// if startBlock-1 is not and cur_block is
	// if start is and end is not
	// both are not
	return sdk.NewInt(0)
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
	newInterest := sdk.NewDecFromInt(debt.Borrowed).
		Mul(params.InterestRate).
		Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(debt.LastInterestCalcTime))).
		Quo(sdk.NewDec(86400 * 365)).
		RoundInt()

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

	debt := k.UpdateInterestStackedByAddress(ctx, addr)
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
	debt := k.UpdateInterestStackedByAddress(ctx, addr)

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
