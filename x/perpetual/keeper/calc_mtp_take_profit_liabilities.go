package keeper

import (
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp types.MTP) (math.Int, error) {
	if mtp.TakeProfitCustody.IsZero() {
		return math.ZeroInt(), nil
	}

	// tradingAssetPriceDenomRatio will also be 0 when tradingAssetPrice is 0
	_, tradingAssetPriceDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return math.ZeroInt(), err
	}

	takeProfitLiabilities := math.ZeroInt()
	if mtp.Position == types.Position_LONG {
		// convert custody amount to base currency, takeProfitCustody is in trading asset, so convert to liabilities asset which is usdc
		// We are not using takeProfitLiabilities anywhere at the moment so weight balance bonus doesn't matter here
		takeProfitLiabilities = mtp.GetBigDecTakeProfitCustody().Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
	} else {
		//  takeProfitCustody is in base currency, so convert to liabilities asset which is trading asset
		// We are not using takeProfitLiabilities anywhere at the moment so weight balance bonus doesn't matter here
		if tradingAssetPriceDenomRatio.IsZero() {
			return math.ZeroInt(), errors.New("trading asset price is zero while calculating takeProfitLiabilities")
		}
		takeProfitLiabilities = mtp.GetBigDecTakeProfitCustody().Quo(tradingAssetPriceDenomRatio).Dec().TruncateInt()
	}

	return takeProfitLiabilities, nil
}
