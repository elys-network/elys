package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdatedLeverage(ctx sdk.Context, mtp types.MTP) (sdk.Dec, error) {

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	custodyInUsdc := mtp.Custody.ToLegacyDec()
	if mtp.CustodyAsset != baseCurrency {
		priceCustodyAsset := k.oracleKeeper.GetAssetPriceFromDenom(ctx, mtp.CustodyAsset)
		custodyInUsdc = mtp.Custody.ToLegacyDec().Mul(priceCustodyAsset)
	}
	var denominator sdk.Dec
	if mtp.LiabilitiesAsset != baseCurrency {
		priceLiablitiesAsset := k.oracleKeeper.GetAssetPriceFromDenom(ctx, mtp.LiabilitiesAsset)
		denominator = custodyInUsdc.Sub(mtp.Liabilities.ToLegacyDec().Mul(priceLiablitiesAsset))
	} else {
		denominator = custodyInUsdc.Sub(mtp.Liabilities.ToLegacyDec())
	}
	if denominator.IsZero() {
		return sdk.ZeroDec(), nil
	}
	effectiveLeverage := custodyInUsdc.Quo(denominator)

	return effectiveLeverage, nil
}
