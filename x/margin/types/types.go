package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func GetPositionFromString(s string) Position {
	switch s {
	case "long":
		return Position_LONG
	case "short":
		return Position_SHORT
	default:
		return Position_UNSPECIFIED
	}
}

func NewMTP(signer string, collateralAsset string, borrowAsset string, position Position, leverage sdk.Dec) *MTP {
	return &MTP{
		Address:                  signer,
		CollateralAsset:          collateralAsset,
		CollateralAmount:         sdk.ZeroInt(),
		Liabilities:              sdk.ZeroInt(),
		InterestPaidCollateral:   sdk.ZeroInt(),
		InterestPaidCustody:      sdk.ZeroInt(),
		InterestUnpaidCollateral: sdk.ZeroInt(),
		CustodyAsset:             borrowAsset,
		CustodyAmount:            sdk.ZeroInt(),
		Leverage:                 leverage,
		MtpHealth:                sdk.ZeroDec(),
		Position:                 position,
	}
}

func (mtp MTP) Validate() error {
	if mtp.CollateralAsset == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no asset specified")
	}
	if mtp.Address == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no address specified")
	}
	if mtp.Position == Position_UNSPECIFIED {
		return sdkerrors.Wrap(ErrMTPInvalid, "no position specified")
	}
	if mtp.Id == 0 {
		return sdkerrors.Wrap(ErrMTPInvalid, "no id specified")
	}

	return nil
}

// Generate a new margin collateral wallet per position
func NewMarginCollateralAddress(positionId uint64) sdk.AccAddress {
	key := append([]byte("margin_collateral"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}

// Generate a new margin custody wallet per position
func NewMarginCustodyAddress(positionId uint64) sdk.AccAddress {
	key := append([]byte("margin_collateral"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}
