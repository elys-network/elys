package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) HandleBorrowInterestPayment(ctx sdk.Context, borrowInterestPayment sdk.Int, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) sdk.Int {
	incrementalBorrowInterestPaymentEnabled := k.GetIncrementalBorrowInterestPaymentEnabled(ctx)
	// if incremental payment on, pay interest
	if incrementalBorrowInterestPaymentEnabled {
		finalBorrowInterestPayment, err := k.IncrementalBorrowInterestPayment(ctx, borrowInterestPayment, mtp, pool, ammPool, baseCurrency)
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrap(err, "error executing incremental borrow interest payment").Error())
		} else {
			return finalBorrowInterestPayment
		}
	} else { // else update unpaid mtp interest
		// collateralAsset is not in base currency
		if mtp.CollateralAsset != baseCurrency {
			// swap
			amtTokenIn := sdk.NewCoin(baseCurrency, borrowInterestPayment)
			var err error
			borrowInterestPayment, err = k.EstimateSwap(ctx, amtTokenIn, mtp.CollateralAsset, ammPool) // may need spot price here to not deduct fee
			if err != nil {
				return sdk.ZeroInt()
			}
		}

		mtp.BorrowInterestUnpaidCollateral = borrowInterestPayment
	}
	return sdk.ZeroInt()
}
