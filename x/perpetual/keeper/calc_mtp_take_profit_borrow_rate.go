package keeper

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
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

func (k Keeper) CalcMTPTakeProfitBorrowFactor(mtp types.MTP) (sdkmath.LegacyDec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return sdkmath.LegacyZeroDec(), types.ErrZeroCustodyAmount
	}

	// infinite for long, 0 for short
	if types.IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return sdk.OneDec(), nil
	}

	takeProfitBorrowFactor := math.LegacyOneDec()
	if mtp.Position == types.Position_LONG {
		// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
		takeProfitBorrowFactor = sdk.OneDec().Sub(mtp.Liabilities.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec().Mul(mtp.TakeProfitPrice)))
	} else {
		// takeProfitBorrowFactor = 1 - ((liabilities  * take profit price) / custody)
		takeProfitBorrowFactor = sdk.OneDec().Sub((mtp.Liabilities.ToLegacyDec().Mul(mtp.TakeProfitPrice)).Quo(mtp.Custody.ToLegacyDec()))

	}

	return takeProfitBorrowFactor, nil
}
