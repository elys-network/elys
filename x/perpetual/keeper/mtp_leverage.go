package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// GetEffectiveLeverage = custody / (custody - liabilities), has to be dimensionless
func (k Keeper) GetEffectiveLeverage(ctx sdk.Context, mtp types.MTP) (math.LegacyDec, error) {

	// using swaps will have fees but this is just user facing value for a query
	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return math.LegacyDec{}, err
	}

	if mtp.Position == types.Position_LONG {
		// custody is trading asset, liabilities are in usdc
		custodyInLiabilitiesAsset := mtp.Custody.ToLegacyDec().Mul(tradingAssetPrice)
		denominator := custodyInLiabilitiesAsset.Sub(mtp.Liabilities.ToLegacyDec())
		effectiveLeverage := custodyInLiabilitiesAsset.Quo(denominator)
		return effectiveLeverage, nil
	} else {
		// custody is usdc, liabilities are in trading asset
		liabilitiesInCustodyAsset := mtp.Liabilities.ToLegacyDec().Mul(tradingAssetPrice)
		denominator := mtp.Custody.ToLegacyDec().Sub(liabilitiesInCustodyAsset)
		effectiveLeverage := mtp.Custody.ToLegacyDec().Quo(denominator)
		return effectiveLeverage, nil
	}
}
