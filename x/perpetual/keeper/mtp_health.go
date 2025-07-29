package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// GetMTPHealth Health = custody / liabilities
// It's responsibility of outer function to update mtp.BorrowInterestUnpaidLiability using UpdateMTPBorrowInterestUnpaidLiability
func (k Keeper) GetMTPHealth(ctx sdk.Context, mtp types.MTP) (math.LegacyDec, error) {

	if mtp.Custody.LTE(math.ZeroInt()) {
		return math.LegacyZeroDec(), nil
	}

	if mtp.Liabilities.IsZero() {
		maxDec := math.LegacyOneDec().Quo(math.LegacySmallestDec())
		return maxDec, nil
	}

	_, tradingAssetPriceDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return math.LegacyDec{}, err
	}

	// For long this unit is base currency, for short this is in trading asset
	// We do not consider here funding fee because it has been / should be already subtracted from mtp.Custody, the custody amt can be <= 0, then above it returns 0
	totalLiabilities := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability).ToLegacyDec()

	// if short position, convert liabilities to base currency
	if mtp.Position == types.Position_SHORT {
		totalLiabilities = tradingAssetPriceDenomRatio.Mul(osmomath.BigDecFromDec(totalLiabilities)).Dec()
		if totalLiabilities.IsZero() {
			return math.LegacyZeroDec(), nil
		}
	}

	// Funding rate is removed as it's subtracted from custody at every epoch

	// For Long this is in trading asset (not base currency, so will have to swap), for Short this is in base currency
	custodyAmtInBaseCurrency := mtp.Custody.ToLegacyDec()

	if !custodyAmtInBaseCurrency.IsPositive() {
		return math.LegacyZeroDec(), nil
	}

	if mtp.Position == types.Position_LONG {
		custodyAmtInBaseCurrency = tradingAssetPriceDenomRatio.Mul(osmomath.BigDecFromDec(custodyAmtInBaseCurrency)).Dec()
	}

	// health = custody / liabilities
	return custodyAmtInBaseCurrency.Quo(totalLiabilities), nil
}
