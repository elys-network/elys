package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdatedLeverage(ctx sdk.Context, mtp types.MTP, pool *types.Pool) (sdk.Dec, error) {

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	var custody_in_usdc sdk.Dec
	if mtp.CustodyAsset != baseCurrency {
		price := k.oracleKeeper.EstimatePrice(ctx, mtp.CustodyAsset, baseCurrency)
		custody_in_usdc = math.LegacyDec(mtp.Custody).Mul(price)
	}else {
		custody_in_usdc = math.LegacyDec(mtp.Custody)
	}
	var denominator sdk.Dec
	if mtp.LiabilitiesAsset != baseCurrency {
		price := k.oracleKeeper.EstimatePrice(ctx, mtp.CustodyAsset, baseCurrency)
		denominator = custody_in_usdc.Sub(math.LegacyDec(mtp.Liabilities).Mul(price))
	}else {
		denominator = custody_in_usdc.Mul(math.LegacyDec(mtp.Liabilities))
	}
	if denominator.IsZero() {
		return sdk.ZeroDec(),  errors.New("")
	}
	updated_leverage := custody_in_usdc.Quo(denominator)

	return updated_leverage, nil
}
