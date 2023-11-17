package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CheckShortAssets(ctx sdk.Context, collateralAsset string, borrowAsset string, baseCurrency string) error {
	// You shouldn't be shorting the base currency (like USDC).
	if borrowAsset == baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "cannot short the base currency")
	}

	// If both the collateralAsset and borrowAsset are the same, it doesn't make sense.
	if collateralAsset == borrowAsset {
		return sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")
	}

	// The collateral for a short must be the base currency.
	if collateralAsset != baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset for a short position must be the base currency")
	}

	return nil
}
