package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CalcMTPInterestLiabilities(mtp *types.MTP, interestRate sdk.Dec, epochPosition, epochLength int64) sdk.Int {
	var interestRational, liabilitiesRational, rate, epochPositionRational, epochLengthRational big.Rat

	rate.SetFloat64(interestRate.MustFloat64())

	liabilitiesRational.SetInt(mtp.Liabilities.BigInt().Add(mtp.Liabilities.BigInt(), mtp.InterestUnpaidCollateral.BigInt()))
	interestRational.Mul(&rate, &liabilitiesRational)

	if epochPosition > 0 { // prorate interest if within epoch
		epochPositionRational.SetInt64(epochPosition)
		epochLengthRational.SetInt64(epochLength)
		epochPositionRational.Quo(&epochPositionRational, &epochLengthRational)
		interestRational.Mul(&interestRational, &epochPositionRational)
	}

	interestNew := interestRational.Num().Quo(interestRational.Num(), interestRational.Denom())

	interestNewInt := sdk.NewIntFromBigInt(interestNew.Add(interestNew, mtp.InterestUnpaidCollateral.BigInt()))
	// round up to lowest digit if interest too low and rate not 0
	if interestNewInt.IsZero() && !interestRate.IsZero() {
		interestNewInt = sdk.NewInt(1)
	}

	return interestNewInt
}
