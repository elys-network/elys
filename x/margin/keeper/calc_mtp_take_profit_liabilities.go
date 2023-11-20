package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp *types.MTP, baseCurrency string) (sdk.Int, error) {
	takeProfitLiabilities := sdk.ZeroInt()
	if types.IsTakeProfitPriceInifite(mtp) {
		return takeProfitLiabilities, nil
	}
	for _, takeProfitCustody := range mtp.TakeProfitCustodies {
		takeProfitCustodyAsset := takeProfitCustody.Denom
		// Retrieve AmmPool
		ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, takeProfitCustodyAsset)
		if err != nil {
			return sdk.ZeroInt(), err
		}

		// convert custody amount to base currency
		C, err := k.EstimateSwapGivenOut(ctx, takeProfitCustody, baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		takeProfitLiabilities = takeProfitLiabilities.Add(C)
	}

	return takeProfitLiabilities, nil
}
