package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// GetEffectiveLeverage = custody / (custody - liabilities), has to be dimensionless
func (k Keeper) GetEffectiveLeverage(ctx sdk.Context, mtp types.MTP) (math.LegacyDec, error) {

	// using swaps will have fees but this is just user facing value for a query
	_, tradingAssetPriceDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return math.LegacyDec{}, err
	}

	if mtp.Position == types.Position_LONG {
		// custody is trading asset, liabilities are in usdc
		custodyInLiabilitiesAsset := mtp.GetBigDecCustody().Mul(tradingAssetPriceDenomRatio)
		denominator := custodyInLiabilitiesAsset.Sub(mtp.GetBigDecLiabilities())
		effectiveLeverage := custodyInLiabilitiesAsset.Quo(denominator)
		return effectiveLeverage.Dec(), nil
	} else {
		// custody is usdc, liabilities are in trading asset
		liabilitiesInCustodyAsset := mtp.GetBigDecLiabilities().Mul(tradingAssetPriceDenomRatio)
		denominator := mtp.GetBigDecCustody().Sub(liabilitiesInCustodyAsset)
		effectiveLeverage := mtp.GetBigDecCustody().Quo(denominator)
		// We subtract here 1 because we added 1 while opening position for shorts
		effectiveLeverage = effectiveLeverage.Sub(osmomath.OneBigDec())
		return effectiveLeverage.Dec(), nil
	}
}
