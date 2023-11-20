package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CheckShortAssets(collateralAsset string, borrowAsset string, baseCurrency string) error {
	// You shouldn't be shorting the base currency (like USDC).
	if borrowAsset == baseCurrency {
		return sdkerrors.Wrap(ErrInvalidBorrowingAsset, "cannot short the base currency")
	}

	// If both the collateralAsset and borrowAsset are the same, it doesn't make sense.
	if collateralAsset == borrowAsset {
		return sdkerrors.Wrap(ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")
	}

	// The collateral for a short must be the base currency.
	if collateralAsset != baseCurrency {
		return sdkerrors.Wrap(ErrInvalidCollateralAsset, "collateral asset for a short position must be the base currency")
	}

	return nil
}
