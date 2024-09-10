package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}

	pendingBorrowInterest := k.GetBorrowInterest(ctx, &mtp, ammPool)
	mtp.BorrowInterestUnpaidCollateral = mtp.BorrowInterestUnpaidCollateral.Add(pendingBorrowInterest)

	// if short position, convert liabilities to base currency
	if mtp.Position == types.Position_SHORT {
		liabilities := sdk.NewCoin(mtp.LiabilitiesAsset, xl)
		var err error
		xl, err = k.EstimateSwapGivenOut(ctx, liabilities, baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroDec(), err
		}

		if xl.IsZero() {
			return sdk.ZeroDec(), nil
		}
	}

	// include unpaid borrow interest in debt (from disabled incremental pay)
	if mtp.BorrowInterestUnpaidCollateral.IsPositive() {
		unpaidCollateral := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)

		if mtp.CollateralAsset == baseCurrency {
			xl = xl.Add(mtp.BorrowInterestUnpaidCollateral)
		} else {
			C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateral, baseCurrency, ammPool)
			if err != nil {
				return sdk.ZeroDec(), err
			}

			xl = xl.Add(C)
		}
	}

	// Funding rate payment consideration
	// get funding rate
	fundingRate := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.AmmPoolId)
	var takeAmountCustodyAmount sdk.Int
	// if funding rate is zero, return
	if fundingRate.IsZero() {
		takeAmountCustodyAmount = sdk.ZeroInt()
	} else if (fundingRate.IsNegative() && mtp.Position == types.Position_LONG) || (fundingRate.IsPositive() && mtp.Position == types.Position_SHORT) {
		takeAmountCustodyAmount = sdk.ZeroInt()
	} else {
		// Calculate the take amount in custody asset
		takeAmountCustodyAmount = types.CalcTakeAmount(mtp.Custody, fundingRate)
	}

	// if short position, custody asset is already in base currency
	custodyAmtInBaseCurrency := mtp.Custody.Sub(takeAmountCustodyAmount)

	if mtp.Position == types.Position_LONG {
		custodyAmt := sdk.NewCoin(mtp.CustodyAsset, mtp.Custody)
		var err error
		custodyAmtInBaseCurrency, err = k.EstimateSwapGivenOut(ctx, custodyAmt, baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroDec(), err
		}
	}

	lr := sdk.NewDecFromBigInt(custodyAmtInBaseCurrency.BigInt()).Quo(sdk.NewDecFromBigInt(xl.BigInt()))

	return lr, nil
}
