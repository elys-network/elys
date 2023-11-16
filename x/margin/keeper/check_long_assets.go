package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) CheckLongAssets(ctx sdk.Context, collateralAsset string, borrowAsset string) error {
	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom
	if borrowAsset == baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset == borrowAsset && collateralAsset == baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset != borrowAsset && collateralAsset != baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	return nil
}
