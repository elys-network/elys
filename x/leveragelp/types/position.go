package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewPosition(signer string, collateral sdk.Coin, poolId uint64) *Position {
	return &Position{
		Address:           signer,
		Collateral:        collateral,
		Liabilities:       sdk.ZeroInt(),
		PositionHealth:    sdk.ZeroDec(),
		AmmPoolId:         poolId,
		LeveragedLpAmount: sdk.ZeroInt(),
	}
}

func (position Position) Validate() error {
	if position.Address == "" {
		return errorsmod.Wrap(ErrPositionInvalid, "no address specified")
	}
	if position.Id == 0 {
		return errorsmod.Wrap(ErrPositionInvalid, "no id specified")
	}

	return nil
}

func (position Position) GetOwnerAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(position.Address)
}

func GetPositionAddress(positionId uint64) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("leveragelp/%d", positionId))
}

// Get Position address
func (p Position) GetPositionAddress() sdk.AccAddress {
	return GetPositionAddress(p.Id)
}
