package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CalcMTPTakeProfitBorrowRate(ctx sdk.Context, mtp *types.MTP) (sdk.Dec, error) {
	// Check if there are any takeProfitCustodies to avoid division by zero
	if len(mtp.TakeProfitCustodies) == 0 {
		return sdk.ZeroDec(), nil
	}

	var totalTakeProfitBorrowRate sdk.Dec = sdk.ZeroDec()
	for takeProfitCustodyIndex, takeProfitCustody := range mtp.TakeProfitCustodies {
		// Calculate the borrow rate for this takeProfitCustody
		takeProfitBorrowRateInt := takeProfitCustody.Amount.Quo(mtp.Custodies[takeProfitCustodyIndex].Amount)

		// Convert takeProfitBorrowRateInt from sdk.Int to sdk.Dec
		takeProfitBorrowRateDec := sdk.NewDecFromInt(takeProfitBorrowRateInt)

		// Add this take profit borrow rate to the total
		totalTakeProfitBorrowRate = totalTakeProfitBorrowRate.Add(takeProfitBorrowRateDec)
	}

	// Calculate the average take profit borrow rate
	averageTakeProfitBorrowRate := totalTakeProfitBorrowRate.Quo(sdk.NewDec(int64(len(mtp.TakeProfitCustodies))))

	// Get Margin Params
	params := k.GetParams(ctx)

	// Use TakeProfitBorrowInterestRateMin param as minimum take profit borrow rate
	averageTakeProfitBorrowRate = sdk.MaxDec(averageTakeProfitBorrowRate, params.TakeProfitBorrowInterestRateMin)

	return averageTakeProfitBorrowRate, nil
}
