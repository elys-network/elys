package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidateMergeMtp(ctx sdk.Context, poolId uint64, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// If Take Profit Price is allowed when adding a new position, the new price for the entire position should be calculated as a weighted Take Profit Price, weighted by the respective positions.
	// If the previous position is 100 ATOM with a Take Profit Price of 10, and the new position is 50 ATOM with a Take Profit Price of 7, the weighted Take Profit Price should be calculated as:
	// (100 * 10 + 50 * 7) / (100 + 50) = 9
	if !types.IsTakeProfitPriceInifite(newMtp) {
		existingCustodyAmt := existingMtp.Custody.ToLegacyDec()
		newCustodyAmt := newMtp.Custody.ToLegacyDec()

		// check no division by zero
		if existingCustodyAmt.Add(newCustodyAmt).IsPositive() {
			existingMtp.TakeProfitPrice = existingMtp.TakeProfitPrice.Mul(existingCustodyAmt).Add(newMtp.TakeProfitPrice.Mul(newCustodyAmt)).Quo(existingCustodyAmt.Add(newCustodyAmt))
		}
	}

	// Merge MTPs
	existingMtp.Collateral = existingMtp.Collateral.Add(newMtp.Collateral)
	existingMtp.Custody = existingMtp.Custody.Add(newMtp.Custody)
	existingMtp.Liabilities = existingMtp.Liabilities.Add(newMtp.Liabilities)
	// Set existing MTP
	if err := k.OpenDefineAssetsChecker.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	// Destroy new MTP
	if err := k.OpenDefineAssetsChecker.DestroyMTP(ctx, newMtp.GetAccountAddress(), newMtp.Id); err != nil {
		return nil, err
	}

	return existingMtp, nil
}
