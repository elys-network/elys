package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewPosition(signer string, collateral sdk.Coin, leverage sdk.Dec, poolId uint64) *Position {
	return &Position{
		Address:           signer,
		Collateral:        collateral,
		Liabilities:       sdk.ZeroInt(),
		Leverage:          leverage,
		PositionHealth:    sdk.ZeroDec(),
		AmmPoolId:         poolId,
		LeveragedLpAmount: sdk.ZeroInt(),
		StopLossPrice:     sdk.ZeroDec(),
	}
}

func (position Position) Validate() error {
	if _, err := sdk.AccAddressFromBech32(position.Address); err != nil {
		return errorsmod.Wrap(ErrPositionInvalid, err.Error())
	}
	if position.Id == 0 {
		return errorsmod.Wrap(ErrPositionInvalid, "position id cannot be 0")
	}
	if position.LeveragedLpAmount.IsNegative() {
		return errorsmod.Wrap(ErrPositionInvalid, "leveraged lp amount cannot be negative")
	}
	if !position.Leverage.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(ErrPositionInvalid, "leverage must be greater than 1")
	}
	if !position.Collateral.IsValid() {
		return errorsmod.Wrap(ErrPositionInvalid, "invalid collateral coin")
	}
	if position.StopLossPrice.IsNegative() {
		return errorsmod.Wrapf(ErrPositionInvalid, "stop loss price cannot be negative")
	}
	if position.Liabilities.IsNegative() {
		return errorsmod.Wrap(ErrPositionInvalid, "liabilities cannot be negative")
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
