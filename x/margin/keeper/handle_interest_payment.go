package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
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

		entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if !found {
			ctx.Logger().Error(sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency).Error())
			return sdk.ZeroInt()
		}
		baseCurrency := entry.Denom

		// collateralAsset is in base currency
		if mtp.Collaterals[collateralIndex].Denom == baseCurrency {
			mtp.InterestUnpaidCollaterals[collateralIndex].Amount = interestPayment
		} else {
			// swap
			amtTokenIn := sdk.NewCoin(baseCurrency, interestPayment)
			interestPayment, err := k.EstimateSwap(ctx, amtTokenIn, collateralAsset, ammPool) // may need spot price here to not deduct fee
			if err != nil {
				return sdk.ZeroInt()
			}

			mtp.InterestUnpaidCollaterals[collateralIndex].Amount = interestPayment
		}
	}
	return sdk.ZeroInt()
}
