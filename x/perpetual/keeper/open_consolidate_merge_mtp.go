package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

func (k Keeper) OpenConsolidateMergeMtp(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP) (*types.MTP, error) {
	// If Take Profit Price is allowed when adding a new position, the new price for the entire position should be calculated as a weighted Take Profit Price, weighted by the respective positions.
	// If the previous position is 100 ATOM with a Take Profit Price of 10, and the new position is 50 ATOM with a Take Profit Price of 7, the weighted Take Profit Price should be calculated as:
	// (100 * 10 + 50 * 7) / (100 + 50) = 9
	if !types.IsTakeProfitPriceInfinite(*newMtp) {
		existingCustodyAmt := existingMtp.GetBigDecCustody()
		newCustodyAmt := newMtp.GetBigDecCustody()

		// check no division by zero
		if existingCustodyAmt.Add(newCustodyAmt).IsPositive() {
			existingMtp.TakeProfitPrice = existingMtp.GetBigDecTakeProfitPrice().Mul(existingCustodyAmt).Add(newMtp.GetBigDecTakeProfitPrice().Mul(newCustodyAmt)).Quo(existingCustodyAmt.Add(newCustodyAmt)).Dec()
		}
	}

	// Merge MTPs
	existingMtp.Collateral = existingMtp.Collateral.Add(newMtp.Collateral)
	existingMtp.Custody = existingMtp.Custody.Add(newMtp.Custody)
	existingMtp.Liabilities = existingMtp.Liabilities.Add(newMtp.Liabilities)
	// Set existing MTP
	if err := k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	// Destroy new MTP
	if err := k.DestroyMTP(ctx, newMtp.GetAccountAddress(), newMtp.Id); err != nil {
		return nil, err
	}

	return existingMtp, nil
}
