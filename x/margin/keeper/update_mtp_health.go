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

	// include unpaid borrow interest in debt (from disabled incremental pay)
	for i := range mtp.Collaterals {
		if mtp.BorrowInterestUnpaidCollaterals[i].Amount.GT(sdk.ZeroInt()) {
			unpaidCollaterals := sdk.NewCoin(mtp.Collaterals[i].Denom, mtp.BorrowInterestUnpaidCollaterals[i].Amount)

			if mtp.Collaterals[i].Denom == baseCurrency {
				xl = xl.Add(mtp.BorrowInterestUnpaidCollaterals[i].Amount)
			} else {
				C, err := k.EstimateSwapGivenOut(ctx, unpaidCollaterals, baseCurrency, ammPool)
				if err != nil {
					return sdk.ZeroDec(), err
				}

				xl = xl.Add(C)
			}
		}
	}

	custodyAmtInBaseCurrency := sdk.ZeroInt()
	for i := range mtp.Custodies {
		custodyTokenIn := sdk.NewCoin(mtp.Custodies[i].Denom, mtp.Custodies[i].Amount)
		// All liabilty is in base currency
		C, err := k.EstimateSwapGivenOut(ctx, custodyTokenIn, baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		custodyAmtInBaseCurrency = custodyAmtInBaseCurrency.Add(C)
	}

	lr := sdk.NewDecFromBigInt(custodyAmtInBaseCurrency.BigInt()).Quo(sdk.NewDecFromBigInt(xl.BigInt()))

	return lr, nil
}
