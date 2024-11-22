package types

import (
	sdkmath "cosmossdk.io/math"
)

func (mtp *MTP) UpdateMTPTakeProfitBorrowFactor() error {
	takeProfitBorrowFactor, err := mtp.CalcMTPTakeProfitBorrowFactor()
	if err != nil {
		return err
	}
	mtp.TakeProfitBorrowFactor = takeProfitBorrowFactor
	return nil
}

func (mtp MTP) CalcMTPTakeProfitBorrowFactor() (sdkmath.LegacyDec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return sdkmath.LegacyZeroDec(), ErrZeroCustodyAmount
	}

	// infinite for long, 0 for short
	if IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return sdkmath.LegacyOneDec(), nil
	}

	takeProfitBorrowFactor := sdkmath.LegacyOneDec()
	if mtp.Position == Position_LONG {
		// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
		takeProfitBorrowFactor = sdkmath.LegacyOneDec().Sub(mtp.Liabilities.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec().Mul(mtp.TakeProfitPrice)))
	} else {
		// takeProfitBorrowFactor = 1 - ((liabilities  * take profit price) / custody)
		takeProfitBorrowFactor = sdkmath.LegacyOneDec().Sub((mtp.Liabilities.ToLegacyDec().Mul(mtp.TakeProfitPrice)).Quo(mtp.Custody.ToLegacyDec()))
	}

	return takeProfitBorrowFactor, nil
}
