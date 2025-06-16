package keeper

import (
	"cosmossdk.io/math"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetAndSetOpenPrice(ctx sdk.Context, mtp *types.MTP, isLeverageZero bool) error {
	if isLeverageZero {
		mtp.OpenPrice = math.LegacyZeroDec()
		return nil
	}
	openPriceDenomRatio := osmomath.ZeroBigDec()
	if mtp.Position == types.Position_LONG {
		if mtp.CollateralAsset == mtp.TradingAsset {
			// open price = liabilities / (custody - collateral)
			denominator := osmomath.BigDecFromSDKInt(mtp.Custody.Sub(mtp.Collateral))
			if denominator.IsZero() {
				return errors.New("(custody - collateral) is zero while calculating open price")
			}
			openPriceDenomRatio = mtp.GetBigDecLiabilities().Quo(denominator)
		} else {
			if mtp.Custody.IsZero() {
				return errors.New("custody is zero while calculating open price")
			}
			// open price = (collateral + liabilities) / custody
			openPriceDenomRatio = osmomath.BigDecFromSDKInt(mtp.Collateral.Add(mtp.Liabilities)).Quo(mtp.GetBigDecCustody())
		}
	} else {
		if mtp.Liabilities.IsZero() {
			// This is special case, when just adding collateral, liabilities of a new MTP before consolidating will be 0
			mtp.OpenPrice = math.LegacyZeroDec()
			return nil
		}
		// open price = (custody - collateral) / liabilities
		openPriceDenomRatio = osmomath.BigDecFromSDKInt(mtp.Custody.Sub(mtp.Collateral)).Quo(mtp.GetBigDecLiabilities())
	}

	openPrice, err := k.ConvertDenomRatioPriceToUSDPrice(ctx, openPriceDenomRatio, mtp.TradingAsset)
	if err != nil {
		return err
	}

	mtp.OpenPrice = openPrice
	return nil
}
