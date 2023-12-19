package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func NewPosition(signer string, collateral sdk.Coin, leverage sdk.Dec, poolId uint64) *Position {
	return &Position{
		Address:           signer,
		Collateral:        collateral,
		Liabilities:       sdk.ZeroInt(),
		InterestPaid:      sdk.ZeroInt(),
		Leverage:          leverage,
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

// Generate a new leveragelp collateral wallet per position
func NewLeveragelpCollateralAddress(positionId uint64) sdk.AccAddress {
	key := append([]byte("leveragelp_collateral"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}

// Generate a new leveragelp custody wallet per position
func NewLeveragelpCustodyAddress(positionId uint64) sdk.AccAddress {
	key := append([]byte("leveragelp_collateral"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}
