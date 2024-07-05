<!--
order: 6
-->

# Functions

## BeginBlocker

The `BeginBlocker` function is called at the beginning of each block to perform necessary updates and maintenance for the `stablestake` module. It updates interest rates and recalculates interest for all debts if an epoch has passed.

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

### Borrow

The `Borrow` function allows a user to borrow a specified amount of tokens, updating the debt and transferring the borrowed tokens to the user's account.

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
```

### Repay

The `Repay` function allows a user to repay a specified amount of borrowed tokens, updating the debt and handling the repayment of interest and principal amounts.

```go
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

### UpdateInterestStacked

The `UpdateInterestStacked` function updates the stacked interest for a given debt based on the current interest rate and the time elapsed since the last interest calculation.

```go
func (k Keeper) UpdateInterestStacked(ctx sdk.Context, debt types.Debt) types.Debt {
    params := k.GetParams(ctx)
    newInterest := sdk.NewDecFromInt(debt.Borrowed).
        Mul(params.InterestRate).
        Mul(sdk.NewDec(ctx.BlockTime().Unix() - int64(debt.LastInterestCalcTime))).
        Quo(sdk.NewDec(86400 * 365)).
        RoundInt()

    debt.InterestStacked = debt.InterestStacked.Add(newInterest)
    debt.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
    k.SetDebt(ctx, debt)

    params.TotalValue = params.TotalValue.Add(newInterest)
    k.SetParams(ctx, params)
    return debt
}
```

### InterestRateComputation

The `InterestRateComputation` function computes the current interest rate based on the network's parameters, total value, and health gain factor.

```go
func (k Keeper) InterestRateComputation(ctx sdk.Context) sdk.Dec {
    params := k.GetParams(ctx)
    if params.TotalValue.IsZero() {
        return params.InterestRate
    }

    interestRateMax := params.InterestRateMax
    interestRateMin := params.InterestRateMin
    interestRateIncrease := params.InterestRateIncrease
    interestRateDecrease := params.InterestRateDecrease
    healthGainFactor := params.HealthGainFactor
    prevInterestRate := params.InterestRate

    moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
    depositDenom := k.GetDepositDenom(ctx)
    balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)
    borrowed := params.TotalValue.Sub(balance.Amount)
    targetInterestRate := healthGainFactor.
        Mul(sdk.NewDecFromInt(borrowed)).
        Quo(sdk.NewDecFromInt(params.TotalValue))

    interestRateChange := targetInterestRate.Sub(prevInterestRate)
    interestRate := prevInterestRate
    if interestRateChange.GTE(interestRateDecrease.Mul(sdk.NewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
        interestRate = targetInterestRate
    } else if interestRateChange.GT(interestRateIncrease) {
        interestRate = prevInterestRate.Add(interestRateIncrease)
    } else if interestRateChange.LT(interestRateDecrease.Mul(sdk.NewDec(-1))) {
        interestRate = prevInterestRate.Sub(interestRateDecrease)
    }

    newInterestRate := interestRate

    if interestRate.GT(interestRateMin) && interestRate.LT(interestRateMax) {
        newInterestRate = interestRate
    } else if interestRate.LTE(interestRateMin) {
        newInterestRate = interestRateMin
    } else if interestRate.GTE(interestRateMax) {
        newInterestRate = interestRateMax
    }

    return newInterestRate
}
```
