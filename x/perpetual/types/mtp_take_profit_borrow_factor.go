package types

import (
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (mtp *MTP) UpdateMTPTakeProfitBorrowFactor() error {
	takeProfitBorrowFactor, err := mtp.CalcMTPTakeProfitBorrowFactor()
	if err != nil {
		return err
	}
	mtp.TakeProfitBorrowFactor = takeProfitBorrowFactor.Dec()
	return nil
}

func (mtp MTP) CalcMTPTakeProfitBorrowFactor() (osmomath.BigDec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return osmomath.ZeroBigDec(), ErrZeroCustodyAmount
	}

	// infinite for long, 0 for short
	if IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return osmomath.OneBigDec(), nil
	}

	takeProfitBorrowFactor := osmomath.OneBigDec()
	if mtp.Position == Position_LONG {
		// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
		takeProfitBorrowFactor = osmomath.OneBigDec().Sub(mtp.GetBigDecLiabilities().Quo(mtp.GetBigDecCustody().MulDec(mtp.TakeProfitPrice)))
	} else {
		// takeProfitBorrowFactor = 1 - ((liabilities  * take profit price) / custody)
		takeProfitBorrowFactor = osmomath.OneBigDec().Sub((mtp.GetBigDecLiabilities().MulDec(mtp.TakeProfitPrice)).Quo(mtp.GetBigDecCustody()))
	}

	return takeProfitBorrowFactor, nil
}
