package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}
	// include unpaid interest in debt (from disabled incremental pay)
	for i := range mtp.CollateralAssets {
		if mtp.InterestUnpaidCollaterals[i].GT(sdk.ZeroInt()) {
			unpaidCollaterals := sdk.NewCoin(mtp.CollateralAssets[i], mtp.InterestUnpaidCollaterals[i])

			if mtp.CollateralAssets[i] == ptypes.BaseCurrency {
				xl = xl.Add(mtp.InterestUnpaidCollaterals[i])
			} else {
				C, err := k.EstimateSwapGivenOut(ctx, unpaidCollaterals, ptypes.BaseCurrency, ammPool)
				if err != nil {
					return sdk.ZeroDec(), err
				}

				xl = xl.Add(C)
			}
		}
	}

	custodyAmtInBaseCurrency := sdk.ZeroInt()
	for i := range mtp.CustodyAssets {
		custodyTokenIn := sdk.NewCoin(mtp.CustodyAssets[i], mtp.CustodyAmounts[i])
		// All liabilty is in base currency
		C, err := k.EstimateSwapGivenOut(ctx, custodyTokenIn, ptypes.BaseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		custodyAmtInBaseCurrency = custodyAmtInBaseCurrency.Add(C)
	}

	lr := sdk.NewDecFromBigInt(custodyAmtInBaseCurrency.BigInt()).Quo(sdk.NewDecFromBigInt(xl.BigInt()))

	return lr, nil
}
