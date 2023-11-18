package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) HandleBorrowInterestPayment(ctx sdk.Context, collateralAsset string, custodyAsset string, borrowInterestPayment sdk.Int, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) sdk.Int {
	incrementalBorrowInterestPaymentEnabled := k.GetIncrementalBorrowInterestPaymentEnabled(ctx)
	// if incremental payment on, pay interest
	if incrementalBorrowInterestPaymentEnabled {
		finalBorrowInterestPayment, err := k.IncrementalBorrowInterestPayment(ctx, collateralAsset, custodyAsset, borrowInterestPayment, mtp, pool, ammPool)
		if err != nil {
			ctx.Logger().Error(sdkerrors.Wrap(err, "error executing incremental borrow interest payment").Error())
		} else {
			return finalBorrowInterestPayment
		}
	} else { // else update unpaid mtp interest
		collateralIndex, _ := types.GetMTPAssetIndex(mtp, collateralAsset, "")
		if collateralIndex < 0 {
			return sdk.ZeroInt()
		}

		// collateralAsset is in base currency
		if mtp.Collaterals[collateralIndex].Denom == baseCurrency {
			mtp.BorrowInterestUnpaidCollaterals[collateralIndex].Amount = borrowInterestPayment
		} else {
			// swap
			amtTokenIn := sdk.NewCoin(baseCurrency, borrowInterestPayment)
			borrowInterestPayment, err := k.EstimateSwap(ctx, amtTokenIn, collateralAsset, ammPool) // may need spot price here to not deduct fee
			if err != nil {
				return sdk.ZeroInt()
			}

			mtp.BorrowInterestUnpaidCollaterals[collateralIndex].Amount = borrowInterestPayment
		}
	}
	return sdk.ZeroInt()
}
