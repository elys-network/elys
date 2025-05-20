package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// GetMTPHealth Health = custody / liabilities
// It's responsibility of outer function to update mtp.BorrowInterestUnpaidLiability using UpdateMTPBorrowInterestUnpaidLiability
func (k Keeper) GetMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) (osmomath.BigDec, error) {

	if mtp.Custody.LTE(math.ZeroInt()) {
		return osmomath.ZeroBigDec(), nil
	}

	if mtp.Liabilities.IsZero() {
		maxDec := osmomath.OneBigDec().Quo(osmomath.SmallestBigDec())
		return maxDec, nil
	}

	// For long this unit is base currency, for short this is in trading asset
	// We do not consider here funding fee because it has been / should be already subtracted from mtp.Custody, the custody amt can be <= 0, then above it returns 0
	totalLiabilities := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)

	// if short position, convert liabilities to base currency
	if mtp.Position == types.Position_SHORT {
		liabilitiesTokenOut := sdk.NewCoin(mtp.LiabilitiesAsset, totalLiabilities)
		var err error
		totalLiabilities, _, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesTokenOut, baseCurrency, ammPool, mtp.Address)
		if err != nil {
			return osmomath.ZeroBigDec(), err
		}

		if totalLiabilities.IsZero() {
			return osmomath.ZeroBigDec(), nil
		}
	}

	// Funding rate is removed as it's subtracted from custody at every epoch

	// For Long this is in trading asset (not base currency, so will have to swap), for Short this is in base currency
	custodyAmtInBaseCurrency := mtp.Custody

	if !custodyAmtInBaseCurrency.IsPositive() {
		return osmomath.ZeroBigDec(), nil
	}

	if mtp.Position == types.Position_LONG {
		custodyAmtTokenOut := sdk.NewCoin(mtp.CustodyAsset, custodyAmtInBaseCurrency)
		var err error
		custodyAmtInBaseCurrency, _, _, err = k.EstimateSwapGivenOut(ctx, custodyAmtTokenOut, baseCurrency, ammPool, mtp.Address)
		if err != nil {
			return osmomath.ZeroBigDec(), err
		}
	}

	// health = custody / liabilities
	lr := osmomath.BigDecFromSDKInt(custodyAmtInBaseCurrency).Quo(osmomath.BigDecFromSDKInt(totalLiabilities))
	return lr, nil
}
