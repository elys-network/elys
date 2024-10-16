package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
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

func NewMTP(signer, collateralAsset, tradingAsset, liabilitiesAsset, custodyAsset string, position Position, takeProfitPrice sdk.Dec, poolId uint64) *MTP {
	return &MTP{
		Address:                       signer,
		CollateralAsset:               collateralAsset,
		TradingAsset:                  tradingAsset,
		LiabilitiesAsset:              liabilitiesAsset,
		CustodyAsset:                  custodyAsset,
		Collateral:                    sdk.ZeroInt(),
		Liabilities:                   sdk.ZeroInt(),
		BorrowInterestPaidCustody:     sdk.ZeroInt(),
		BorrowInterestUnpaidLiability: sdk.ZeroInt(),
		Custody:                       sdk.ZeroInt(),
		TakeProfitLiabilities:         sdk.ZeroInt(),
		TakeProfitCustody:             sdk.ZeroInt(),
		MtpHealth:                     sdk.ZeroDec(),
		Position:                      position,
		Id:                            0,
		AmmPoolId:                     poolId,
		TakeProfitPrice:               takeProfitPrice,
		TakeProfitBorrowFactor:        sdk.OneDec(),
		FundingFeePaidCollateral:      sdk.ZeroInt(),
		FundingFeePaidCustody:         sdk.ZeroInt(),
		FundingFeeReceivedCollateral:  sdk.ZeroInt(),
		FundingFeeReceivedCustody:     sdk.ZeroInt(),
		OpenPrice:                     sdk.ZeroDec(),
	}
}

func (mtp MTP) Validate() error {
	if mtp.CollateralAsset == "" {
		return errorsmod.Wrap(ErrMTPInvalid, "no collateral asset specified")
	}
	if mtp.CustodyAsset == "" {
		return errorsmod.Wrap(ErrMTPInvalid, "no custody asset specified")
	}
	if mtp.Address == "" {
		return errorsmod.Wrap(ErrMTPInvalid, "no address specified")
	}
	if mtp.Position == Position_UNSPECIFIED {
		return errorsmod.Wrap(ErrMTPInvalid, "no position specified")
	}
	if mtp.Id == 0 {
		return errorsmod.Wrap(ErrMTPInvalid, "no id specified")
	}

	return nil
}

func (mtp MTP) GetAccountAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(mtp.Address)
}

// Generate a new perpetual collateral wallet per position
func NewPerpetualCollateralAddress(positionId uint64) sdk.AccAddress {
	key := append([]byte("perpetual_collateral"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}

// Generate a new perpetual custody wallet per position
func NewPerpetualCustodyAddress(positionId uint64) sdk.AccAddress {
	key := append([]byte("perpetual_custody"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}
