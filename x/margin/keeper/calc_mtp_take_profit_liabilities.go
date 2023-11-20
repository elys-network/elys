package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp *types.MTP, baseCurrency string) (sdk.Int, error) {
	takeProfitLiabilities := sdk.ZeroInt()

	for custodyIndex, custody := range mtp.Custodies {
		custodyAsset := custody.Denom
		// Retrieve AmmPool
		ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			return sdk.ZeroInt(), err
		}

		// update mtp take profit liabilities
		C, err := k.EstimateSwapGivenOut(ctx, mtp.TakeProfitCustodies[custodyIndex], baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		takeProfitLiabilities = takeProfitLiabilities.Add(C)
	}

	return takeProfitLiabilities, nil
}
