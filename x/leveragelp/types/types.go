package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMTP(signer string, collateral sdk.Coin, leverage sdk.Dec, poolId uint64) *MTP {
	return &MTP{
		Address:           signer,
		Collateral:        collateral,
		Liabilities:       sdk.ZeroInt(),
		InterestPaid:      sdk.ZeroInt(),
		Leverage:          leverage,
		MtpHealth:         sdk.ZeroDec(),
		AmmPoolId:         poolId,
		LeveragedLpAmount: sdk.ZeroInt(),
	}
}

func (mtp MTP) Validate() error {
	if mtp.Address == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no address specified")
	}
	if mtp.Id == 0 {
		return sdkerrors.Wrap(ErrMTPInvalid, "no id specified")
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
