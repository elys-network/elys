<!--
order: 4
-->

# Keeper

## Interest Rate Management

The `stablestake` module's keeper handles the computation and updating of interest rates, ensuring they are adjusted based on the network's parameters and conditions.

### BeginBlocker

The `BeginBlocker` function is invoked at the beginning of each block. It checks if an epoch has passed, updates interest rates, and recalculates the stacked interest for all debts.

```go
func (k Keeper) BeginBlocker(ctx sdk.Context) {
    // check if epoch has passed then execute
    epochLength := k.GetEpochLength(ctx)
    epochPosition := k.GetEpochPosition(ctx, epochLength)

    if epochPosition == 0 { // if epoch has passed
        params := k.GetParams(ctx)
        rate := k.InterestRateComputation(ctx)
        params.InterestRate = rate
        k.SetParams(ctx, params)

        debts := k.AllDebts(ctx)
        for _, debt := range debts {
            k.UpdateInterestStacked(ctx, debt)
        }
    }
}
```

### Borrowing and Repaying

The `Borrow` function allows a user to borrow tokens, while the `Repay` function allows a user to repay borrowed tokens, including any accrued interest.

```go
func (k Keeper) Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
    depositDenom := k.GetDepositDenom(ctx)
    if depositDenom != amount.Denom {
        return types.ErrInvalidBorrowDenom
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

    if !debt.Borrowed.IsPositive() {
        k.DeleteDebt(ctx, debt)
    } else {
        k.SetDebt(ctx, debt)
    }
    return nil
}
```
