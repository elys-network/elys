package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp *types.MTP, baseCurrency string) (sdk.Int, error) {
	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// build take profit custody coin
	takeProfitCustody := sdk.NewCoin(mtp.CustodyAsset, mtp.TakeProfitCustody)

	// convert custody amount to base currency
	takeProfitLiabilities, err := k.EstimateSwap(ctx, takeProfitCustody, baseCurrency, ammPool)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	return takeProfitLiabilities, nil
}
