package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (mtp *MTP) UpdateMTPTakeProfitBorrowFactor() error {
	takeProfitBorrowFactor, err := mtp.CalcMTPTakeProfitBorrowFactor()
	if err != nil {
		return err
	}
	mtp.TakeProfitBorrowFactor = takeProfitBorrowFactor
	return nil
}

func (mtp MTP) CalcMTPTakeProfitBorrowFactor() (sdk.Dec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return sdk.ZeroDec(), ErrZeroCustodyAmount
	}

	// infinite for long, 0 for short
	if IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return sdk.OneDec(), nil
	}

	takeProfitBorrowFactor := math.LegacyOneDec()
	if mtp.Position == Position_LONG {
		// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
		takeProfitBorrowFactor = sdk.OneDec().Sub(mtp.Liabilities.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec().Mul(mtp.TakeProfitPrice)))
	} else {
		// takeProfitBorrowFactor = 1 - ((liabilities  * take profit price) / custody)
		takeProfitBorrowFactor = sdk.OneDec().Sub((mtp.Liabilities.ToLegacyDec().Mul(mtp.TakeProfitPrice)).Quo(mtp.Custody.ToLegacyDec()))
	}

	return takeProfitBorrowFactor, nil
}
