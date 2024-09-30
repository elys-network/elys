package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPTakeProfitBorrowRate(ctx sdk.Context, mtp *types.MTP) (sdkmath.LegacyDec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return sdkmath.LegacyZeroDec(), types.ErrAmountTooLow
	}

	// Calculate the borrow rate for this takeProfitCustody
	takeProfitBorrowRateInt := mtp.TakeProfitCustody.Quo(mtp.Custody)

	// Convert takeProfitBorrowRateInt from math.Int to sdkmath.LegacyDec
	takeProfitBorrowRateDec := sdkmath.LegacyNewDecFromInt(takeProfitBorrowRateInt)

	// Get Perpetual Params
	params := k.GetParams(ctx)

	// Use TakeProfitBorrowInterestRateMin param as minimum take profit borrow rate
	takeProfitBorrowRate := sdkmath.LegacyMaxDec(takeProfitBorrowRateDec, params.TakeProfitBorrowInterestRateMin)

	return takeProfitBorrowRate, nil
}
