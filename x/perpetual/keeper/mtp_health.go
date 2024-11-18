package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// GetMTPHealth Health = custody / liabilities
// It's responsibility of outer function to update mtp.BorrowInterestUnpaidLiability using UpdateMTPBorrowInterestUnpaidLiability
func (k Keeper) GetMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) (math.LegacyDec, error) {
	if mtp.Liabilities.IsZero() {
		return math.LegacyMaxSortableDec, nil
	}

	// For long this unit is base currency, for short this is in trading asset
	totalLiabilities := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)

	// if short position, convert liabilities to base currency
	if mtp.Position == types.Position_SHORT {
		liabilitiesTokenOut := sdk.NewCoin(mtp.LiabilitiesAsset, totalLiabilities)
		var err error
		totalLiabilities, _, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesTokenOut, baseCurrency, ammPool, mtp.Address)
		if err != nil {
			return math.LegacyZeroDec(), err
		}

		if totalLiabilities.IsZero() {
			return math.LegacyZeroDec(), nil
		}
	}

	// Funding rate is removed as it's subtracted from custody at every epoch

	// For Long this is in trading asset (not base currency, so will have to swap), for Short this is in base currency
	custodyAmtInBaseCurrency := mtp.Custody

	if !custodyAmtInBaseCurrency.IsPositive() {
		return math.LegacyZeroDec(), nil
	}

	if mtp.Position == types.Position_LONG {
		custodyAmtTokenOut := sdk.NewCoin(mtp.CustodyAsset, custodyAmtInBaseCurrency)
		var err error
		custodyAmtInBaseCurrency, _, _, err = k.EstimateSwapGivenOut(ctx, custodyAmtTokenOut, baseCurrency, ammPool, mtp.Address)
		if err != nil {
			return math.LegacyZeroDec(), err
		}
	}

	// health = custody / liabilities
	lr := custodyAmtInBaseCurrency.ToLegacyDec().Quo(totalLiabilities.ToLegacyDec())
	return lr, nil
}
