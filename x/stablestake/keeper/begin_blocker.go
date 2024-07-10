package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		// divide them in blocks, update values
		params := k.GetParams(ctx)
		rate := k.InterestRateComputation(ctx)
		params.InterestRate = rate
		k.SetParams(ctx, params)

		debts := k.AllDebts(ctx)
		for _, debt := range debts {
			old := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
			k.UpdateInterestStacked(ctx, debt)
			k.hooks.AfterUpdateInterestStacked(ctx, debt.Address, old, debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid))
		}
	}
}
