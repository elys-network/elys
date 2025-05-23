package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/utils"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) UpdateOpenPrice(ctx sdk.Context, mtp *types.MTP) error {
	err := k.GetAndSetOpenPrice(ctx, mtp)
	if err != nil {
		return err
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAndSetOpenPrice(ctx sdk.Context, mtp *types.MTP) error {
	openPrice := osmomath.ZeroBigDec()
	if mtp.Position == types.Position_LONG {
		if mtp.CollateralAsset == mtp.TradingAsset {
			// open price = liabilities / (custody - collateral)
			denominator := osmomath.BigDecFromSDKInt(mtp.Custody.Sub(mtp.Collateral))
			if denominator.IsZero() {
				return errors.New("(custody - collateral) is zero while calculating open price")
			}
			openPrice = mtp.GetBigDecLiabilities().Quo(denominator)
		} else {
			if mtp.Custody.IsZero() {
				return errors.New("custody is zero while calculating open price")
			}
			// open price = (collateral + liabilities) / custody
			openPrice = osmomath.BigDecFromSDKInt(mtp.Collateral.Add(mtp.Liabilities)).Quo(mtp.GetBigDecCustody())
		}
	} else {
		if mtp.Liabilities.IsZero() {
			// This is special case, when just adding collateral, liabilities of a new MTP before consolidating will be 0
			mtp.OpenPrice = osmomath.ZeroBigDec().Dec()
			return nil
		}
		// open price = (custody - collateral) / liabilities
		openPrice = osmomath.BigDecFromSDKInt(mtp.Custody.Sub(mtp.Collateral)).Quo(mtp.GetBigDecLiabilities())
	}

	// right now units of open price is in base units, ex: uusdc per uatom, uusdc per wei, uusdc per satoshi
	// we need to convert to usd per atom, usd per eth, usd per btc
	baseCurrencyDenom := mtp.LiabilitiesAsset
	if mtp.IsShort() {
		baseCurrencyDenom = mtp.CollateralAsset
	}
	baseCurrencyDenomPrice, err := k.GetDenomPrice(ctx, baseCurrencyDenom)
	if err != nil {
		return err
	}

	// Now the units are usd per uatom, usd per wei, usd per sat
	openPrice = openPrice.Mul(baseCurrencyDenomPrice) // Now the units are usd per uatom, usd per wei, usd per sat

	decimal, err := k.GetDenomDecimal(ctx, mtp.TradingAsset)
	if err != nil {
		return err
	}

	// Multiply by 10^decimal of taring asset
	openPrice = openPrice.MulInt64(utils.Pow10Int64(decimal))

	mtp.OpenPrice = openPrice.Dec()
	return nil
}
