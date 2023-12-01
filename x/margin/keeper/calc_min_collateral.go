package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
)

// CalcMinCollateral calculates the minimum collateral required to open a position
func (k Keeper) CalcMinCollateral(ctx sdk.Context, position types.Position, leverage sdk.Dec, tradingAsset string, collateralAsset string) (sdk.Dec, error) {
	// leverage must be greater than 1
	if leverage.LTE(sdk.NewDec(1)) {
		return sdk.ZeroDec(), sdkerrors.Wrapf(types.ErrInvalidLeverage, "leverage must be greater than 1")
	}

	// get collateral asset decimals
	entry, found := k.apKeeper.GetEntry(ctx, collateralAsset)
	if !found {
		return sdk.ZeroDec(), sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", collateralAsset)
	}
	collateralDecimals := entry.Decimals

	// get min borrow rate
	borrowInterestRateMin := k.GetBorrowInterestRateMin(ctx)

	// min_collateral = ( 1 / (( leverage - 1 ) * borrow_interest_rate_min )) / 10 ^ collateral_decimals
	minCollateral := sdk.NewDec(1).Quo(
		leverage.Sub(sdk.NewDec(1)).Mul(borrowInterestRateMin),
	).Quo(
		sdk.NewDec(10).Power(collateralDecimals),
	)

	return minCollateral, nil
}
