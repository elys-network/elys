package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdatedLeverage(ctx sdk.Context, mtp types.MTP) (sdk.Dec, error) {
	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return sdk.ZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	collateral_in_usdc := mtp.Collateral.ToLegacyDec()
	if mtp.CollateralAsset != baseCurrency {
		price := k.amm.EstimatePrice(ctx, mtp.CollateralAsset, baseCurrency)
		collateral_in_usdc = mtp.Collateral.ToLegacyDec().Mul(price)
	}
	liablites := mtp.Liabilities.ToLegacyDec()
	if mtp.LiabilitiesAsset != baseCurrency {
		price := k.amm.EstimatePrice(ctx, mtp.LiabilitiesAsset, baseCurrency)
		liablites = mtp.Liabilities.ToLegacyDec().Mul(price)
	}
	if collateral_in_usdc.IsZero() {
		return sdk.ZeroDec(),  nil
	}
	updated_leverage := liablites.Quo(collateral_in_usdc).Add(sdk.OneDec())

	return updated_leverage, nil
}
