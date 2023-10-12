package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}

	// TODO: calculate the value of leveragelp tokens holding
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
