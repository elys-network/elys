package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetDebt(ctx sdk.Context, addr sdk.AccAddress) types.Debt {
	debt := types.Debt{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DebtPrefixKey)
	bz := store.Get([]byte(addr.String()))
	if len(bz) == 0 {
		return types.Debt{
			Address: addr.String(),
			Debt:    sdk.ZeroInt(),
		}
	}

	k.cdc.MustUnmarshal(bz, &debt)
	return debt
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

func (k Keeper) Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	params := k.GetParams(ctx)
	if params.DepositDenom != amount.Denom {
		return types.ErrInvalidBorrowDenom
	}
	debt := k.GetDebt(ctx, addr)
	debt.Debt = debt.Debt.Add(amount.Amount)
	k.SetDebt(ctx, debt)
	return k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{amount})
}

func (k Keeper) Repay(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	params := k.GetParams(ctx)
	if params.DepositDenom != amount.Denom {
		return types.ErrInvalidBorrowDenom
	}

	err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}

	debt := k.GetDebt(ctx, addr)
	debt.Debt = debt.Debt.Sub(amount.Amount)

	if !debt.Debt.IsPositive() {
		k.DeleteDebt(ctx, debt)
	} else {
		k.SetDebt(ctx, debt)
	}
	return nil
}
