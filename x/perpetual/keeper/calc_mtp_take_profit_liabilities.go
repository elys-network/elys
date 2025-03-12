package keeper

import (
	"cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp types.MTP) (math.Int, error) {
	if mtp.TakeProfitCustody.IsZero() {
		return math.ZeroInt(), nil
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return math.ZeroInt(), err
	}

	takeProfitLiabilities := math.ZeroInt()
	if mtp.Position == types.Position_LONG {
		// convert custody amount to base currency, takeProfitCustody is in trading asset, so convert to liabilities asset which is usdc
		// We are not using takeProfitLiabilities anywhere at the moment so weight balance bonus doesn't matter here
		takeProfitLiabilities = mtp.TakeProfitCustody.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
	} else {
		//  takeProfitCustody is in base currency, so convert to liabilities asset which is trading asset
		// We are not using takeProfitLiabilities anywhere at the moment so weight balance bonus doesn't matter here
		if tradingAssetPrice.IsZero() {
			return math.ZeroInt(), errors.New("trading asset price is zero while calculating takeProfitLiabilities")
		}
		takeProfitLiabilities = mtp.TakeProfitCustody.ToLegacyDec().Quo(tradingAssetPrice).TruncateInt()
	}

	return takeProfitLiabilities, nil
}
