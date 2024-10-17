package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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
	collateral_in_usdc := math.LegacyDec(mtp.Collateral)
	if mtp.CollateralAsset != baseCurrency {
		price := k.amm.EstimatePrice(ctx, mtp.CollateralAsset, baseCurrency)
		collateral_in_usdc = math.LegacyDec(mtp.Liabilities).Mul(price)
	}
	liablites := math.LegacyDec(mtp.Liabilities)
	if mtp.LiabilitiesAsset != baseCurrency {
		price := k.amm.EstimatePrice(ctx, mtp.LiabilitiesAsset, baseCurrency)
		liablites = math.LegacyDec(mtp.Liabilities).Mul(price)
	}
	if collateral_in_usdc.IsZero() {
		return sdk.ZeroDec(),  nil
	}
	updated_leverage := liablites.Quo(collateral_in_usdc).Add(sdk.OneDec())

	return updated_leverage, nil
}
