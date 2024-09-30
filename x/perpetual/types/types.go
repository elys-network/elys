package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
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

func NewMTP(signer, collateralAsset, tradingAsset, liabilitiesAsset, custodyAsset string, position Position, leverage, takeProfitPrice sdkmath.LegacyDec, poolId uint64) *MTP {
	return &MTP{
		Address:                        signer,
		CollateralAsset:                collateralAsset,
		TradingAsset:                   tradingAsset,
		LiabilitiesAsset:               liabilitiesAsset,
		CustodyAsset:                   custodyAsset,
		Collateral:                     sdkmath.ZeroInt(),
		Liabilities:                    sdkmath.ZeroInt(),
		BorrowInterestPaidCollateral:   sdkmath.ZeroInt(),
		BorrowInterestPaidCustody:      sdkmath.ZeroInt(),
		BorrowInterestUnpaidCollateral: sdkmath.ZeroInt(),
		Custody:                        sdkmath.ZeroInt(),
		TakeProfitLiabilities:          sdkmath.ZeroInt(),
		TakeProfitCustody:              sdkmath.ZeroInt(),
		Leverage:                       leverage,
		MtpHealth:                      sdkmath.LegacyZeroDec(),
		Position:                       position,
		Id:                             0,
		AmmPoolId:                      poolId,
		ConsolidateLeverage:            leverage,
		SumCollateral:                  sdkmath.ZeroInt(),
		TakeProfitPrice:                takeProfitPrice,
		TakeProfitBorrowRate:           sdkmath.LegacyOneDec(),
		FundingFeePaidCollateral:       sdkmath.ZeroInt(),
		FundingFeePaidCustody:          sdkmath.ZeroInt(),
		FundingFeeReceivedCollateral:   sdkmath.ZeroInt(),
		FundingFeeReceivedCustody:      sdkmath.ZeroInt(),
		OpenPrice:                      sdkmath.LegacyZeroDec(),
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
