package keeper

import (
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPBorrowInterestLiabilities(ctx sdk.Context, mtp *types.MTP, borrowInterestRate sdk.Dec, epochPosition, epochLength int64, ammPool ammtypes.Pool, baseCurrency string) (math.Int, error) {
	// Ensure borrow interest rate or liabilities are not zero to avoid division by zero
	if borrowInterestRate.IsZero() || mtp.Liabilities.IsZero() {
		return sdk.ZeroInt(), types.ErrAmountTooLow
	}

	var borrowInterestRational, liabilitiesRational, rate, epochPositionRational, epochLengthRational big.Rat

	rate.SetFloat64(borrowInterestRate.MustFloat64())

	unpaidCollateral := sdk.ZeroInt()
	// Calculate collateral borrow interests in base currency
	if mtp.CollateralAsset == baseCurrency {
		unpaidCollateral = unpaidCollateral.Add(mtp.BorrowInterestUnpaidCollateral)
	} else {
		// Liability is in base currency, so convert it to base currency
		unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
		C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroInt(), err
		}

		unpaidCollateral = unpaidCollateral.Add(C)
	}

	liabilitiesRational.SetInt(mtp.Liabilities.BigInt().Add(mtp.Liabilities.BigInt(), unpaidCollateral.BigInt()))
	borrowInterestRational.Mul(&rate, &liabilitiesRational)

	if epochPosition > 0 { // prorate borrow interest if within epoch
		epochPositionRational.SetInt64(epochPosition)
		epochLengthRational.SetInt64(epochLength)
		epochPositionRational.Quo(&epochPositionRational, &epochLengthRational)
		borrowInterestRational.Mul(&borrowInterestRational, &epochPositionRational)
	}

	borrowInterestNew := borrowInterestRational.Num().Quo(borrowInterestRational.Num(), borrowInterestRational.Denom())

	borrowInterestNewInt := sdk.NewIntFromBigInt(borrowInterestNew.Add(borrowInterestNew, unpaidCollateral.BigInt()))
	// round up to lowest digit if borrow interest too low and rate not 0
	if borrowInterestNewInt.IsZero() && !borrowInterestRate.IsZero() {
		borrowInterestNewInt = sdk.NewInt(1)
	}

	// apply take profit borrow rate to borrow interest
	borrowInterestNewInt = sdk.NewDecFromInt(borrowInterestNewInt).Mul(mtp.TakeProfitBorrowRate).TruncateInt()

	return borrowInterestNewInt, nil
}
