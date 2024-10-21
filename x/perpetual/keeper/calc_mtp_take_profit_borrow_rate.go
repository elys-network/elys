package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateMTPTakeProfitBorrowFactor(ctx sdk.Context, mtp *types.MTP) error {
	takeProfitBorrowFactor, err := k.CalcMTPTakeProfitBorrowFactor(*mtp)
	if err != nil {
		return err
	}
	mtp.TakeProfitBorrowFactor = takeProfitBorrowFactor
	return nil
}

func (k Keeper) CalcMTPTakeProfitBorrowFactor(mtp types.MTP) (sdk.Dec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return sdk.ZeroDec(), types.ErrZeroCustodyAmount
	}

	if types.IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return sdk.OneDec(), nil
	}

	// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
	takeProfitBorrowFactor := sdk.OneDec().Sub(mtp.Liabilities.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec().Mul(mtp.TakeProfitPrice)))

	return takeProfitBorrowFactor, nil
}
