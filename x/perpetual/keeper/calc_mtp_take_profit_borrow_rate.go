package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPTakeProfitBorrowRate(ctx sdk.Context, mtp *types.MTP) (sdk.Dec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return sdk.ZeroDec(), types.ErrAmountTooLow
	}

	// Calculate the borrow rate for this takeProfitCustody
	takeProfitBorrowRateInt := mtp.TakeProfitCustody.Quo(mtp.Custody)

	// Convert takeProfitBorrowRateInt from math.Int to sdk.Dec
	takeProfitBorrowRateDec := sdk.NewDecFromInt(takeProfitBorrowRateInt)

	// Get Perpetual Params
	params := k.GetParams(ctx)

	// Use TakeProfitBorrowInterestRateMin param as minimum take profit borrow rate
	takeProfitBorrowRate := sdk.MaxDec(takeProfitBorrowRateDec, params.TakeProfitBorrowInterestRateMin)

	return takeProfitBorrowRate, nil
}
