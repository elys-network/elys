package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp *types.MTP, baseCurrency string) (math.Int, error) {
	if mtp.TakeProfitCustody.IsZero() {
		return math.ZeroInt(), nil
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return math.ZeroInt(), err
	}

	// build take profit custody coin
	takeProfitCustody := sdk.NewCoin(mtp.CustodyAsset, mtp.TakeProfitCustody)

	takeProfitLiabilities := math.ZeroInt()
	if mtp.Position == types.Position_LONG {
		// convert custody amount to base currency, takeProfitCustody is in trading asset, so convert to liabilities asset which is usdc
		takeProfitLiabilities, _, err = k.EstimateSwapGivenIn(ctx, takeProfitCustody, baseCurrency, ammPool)
		if err != nil {
			return math.ZeroInt(), errorsmod.Wrapf(err, "unable to swap takeProfitCustody to baseCurrency for takeProfitLiabilities")
		}
	} else {
		//  takeProfitCustody is in base currency, so convert to liabilities asset which is trading asset
		takeProfitLiabilities, _, err = k.EstimateSwapGivenIn(ctx, takeProfitCustody, mtp.LiabilitiesAsset, ammPool)
		if err != nil {
			return sdk.ZeroInt(), errorsmod.Wrapf(err, "unable to swap takeProfitCustody to LiabilitiesAsset for takeProfitLiabilities")
		}
	}

	return takeProfitLiabilities, nil
}
