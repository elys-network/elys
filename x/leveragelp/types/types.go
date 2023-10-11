package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMTP(signer string, collateralAsset string, leverage sdk.Dec, poolId uint64) *MTP {
	return &MTP{
		Address:                 signer,
		CollateralAssets:        []string{collateralAsset},
		CollateralAmounts:       []sdk.Int{sdk.ZeroInt()},
		Liabilities:             sdk.ZeroInt(),
		InterestPaidCollaterals: []sdk.Int{sdk.ZeroInt()},
		Leverages:               []sdk.Dec{leverage},
		MtpHealth:               sdk.ZeroDec(),
		AmmPoolId:               poolId,
		ConsolidateLeverage:     leverage,
		SumCollateral:           sdk.ZeroInt(),
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
