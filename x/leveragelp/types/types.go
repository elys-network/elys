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

func NewMTP(signer string, collateralAsset string, borrowAsset string, position Position, leverage sdk.Dec, poolId uint64) *MTP {
	return &MTP{
		Address:                   signer,
		CollateralAssets:          []string{collateralAsset},
		CollateralAmounts:         []sdk.Int{sdk.ZeroInt()},
		Liabilities:               sdk.ZeroInt(),
		InterestPaidCollaterals:   []sdk.Int{sdk.ZeroInt()},
		InterestPaidCustodys:      []sdk.Int{sdk.ZeroInt()},
		InterestUnpaidCollaterals: []sdk.Int{sdk.ZeroInt()},
		CustodyAssets:             []string{borrowAsset},
		CustodyAmounts:            []sdk.Int{sdk.ZeroInt()},
		Leverages:                 []sdk.Dec{leverage},
		MtpHealth:                 sdk.ZeroDec(),
		Position:                  position,
		AmmPoolId:                 poolId,
		ConsolidateLeverage:       leverage,
		SumCollateral:             sdk.ZeroInt(),
	}
}

func (mtp MTP) Validate() error {
	if len(mtp.CollateralAssets) < 1 {
		return sdkerrors.Wrap(ErrMTPInvalid, "no asset specified")
	}
	for _, asset := range mtp.CollateralAssets {
		if asset == "" {
			return sdkerrors.Wrap(ErrMTPInvalid, "no asset specified")
		}
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
