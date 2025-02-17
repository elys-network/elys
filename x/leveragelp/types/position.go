package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewPosition(signer string, collateral sdk.Coin, poolId uint64) *Position {
	return &Position{
		Address:           signer,
		Collateral:        collateral,
		Liabilities:       sdkmath.ZeroInt(),
		PositionHealth:    sdkmath.LegacyZeroDec(),
		AmmPoolId:         poolId,
		LeveragedLpAmount: sdkmath.ZeroInt(),
		StopLossPrice:     sdkmath.LegacyZeroDec(),
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

func (p Position) CheckStopLossReached(lpTokenPrice sdkmath.LegacyDec) bool {
	return !p.StopLossPrice.IsNil() && lpTokenPrice.LTE(p.StopLossPrice)
}
