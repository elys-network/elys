package keeper

import (
	"fmt"

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

	priceBaseCurrency := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	if priceBaseCurrency.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", baseCurrency)
	}

	custodyInUsdc := mtp.Custody.ToLegacyDec()
	if mtp.CustodyAsset != baseCurrency {
		priceCustodyAsset := k.oracleKeeper.GetAssetPriceFromDenom(ctx, mtp.CustodyAsset)
		if priceBaseCurrency.IsZero() {
			return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", mtp.CustodyAsset)
		}
		custodyInUsdc = mtp.Custody.ToLegacyDec().Mul(priceCustodyAsset.Quo(priceBaseCurrency))
	}
	denominator := custodyInUsdc.Sub(mtp.Liabilities.ToLegacyDec())
	if mtp.LiabilitiesAsset != baseCurrency {
		priceCustodyAsset := k.oracleKeeper.GetAssetPriceFromDenom(ctx, mtp.LiabilitiesAsset)
		if priceBaseCurrency.IsZero() {
			return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", mtp.LiabilitiesAsset)
		}
		custodyInUsdc = mtp.Custody.ToLegacyDec().Mul(priceCustodyAsset.Quo(priceBaseCurrency))
	}
	if denominator.IsZero() {
		return sdk.ZeroDec(),  nil
	}
	effectiveLeverage := custodyInUsdc.Quo(denominator)

	return effectiveLeverage, nil
}
