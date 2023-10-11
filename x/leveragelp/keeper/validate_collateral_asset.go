package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) ValidateCollateralAsset(collateralAsset string) error {
	if collateralAsset != paramtypes.BaseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "invalid collateral asset")
	}
	return nil
}
