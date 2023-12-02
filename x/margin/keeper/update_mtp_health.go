package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool, baseCurrency string) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}

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

	// if short position, custody asset is already in base currency
	custodyAmtInBaseCurrency := mtp.Custody

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
