package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) HandleInterestPayment(ctx sdk.Context, collateralAsset string, custodyAsset string, interestPayment sdk.Int, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) sdk.Int {
	incrementalInterestPaymentEnabled := k.GetIncrementalInterestPaymentEnabled(ctx)
	// if incremental payment on, pay interest
	if incrementalInterestPaymentEnabled {
		finalInterestPayment, err := k.IncrementalInterestPayment(ctx, collateralAsset, custodyAsset, interestPayment, mtp, pool, ammPool)
		if err != nil {
			ctx.Logger().Error(sdkerrors.Wrap(err, "error executing incremental interest payment").Error())
		} else {
			return finalInterestPayment
		}
	} else { // else update unpaid mtp interest
		collateralIndex, _ := k.GetMTPAssetIndex(mtp, collateralAsset, "")
		if collateralIndex < 0 {
			return sdk.ZeroInt()
		}

		// collateralAsset is in base currency
		if mtp.CollateralAssets[collateralIndex] == ptypes.BaseCurrency {
			mtp.InterestUnpaidCollaterals[collateralIndex] = interestPayment
		} else {
			// swap
			amtTokenIn := sdk.NewCoin(ptypes.BaseCurrency, interestPayment)
			interestPayment, err := k.EstimateSwap(ctx, amtTokenIn, collateralAsset, ammPool) // may need spot price here to not deduct fee
			if err != nil {
				return sdk.ZeroInt()
			}

			mtp.InterestUnpaidCollaterals[collateralIndex] = interestPayment
		}
	}
	return sdk.ZeroInt()
}
