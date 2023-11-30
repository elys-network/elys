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

func NewMTP(signer, collateralAsset, tradingAsset, liabilitiesAsset, custodyAsset string, position Position, leverage, takeProfitPrice sdk.Dec, poolId uint64) *MTP {
	return &MTP{
		Address:                        signer,
		CollateralAsset:                collateralAsset,
		TradingAsset:                   tradingAsset,
		LiabilitiesAsset:               liabilitiesAsset,
		CustodyAsset:                   custodyAsset,
		Collateral:                     sdk.ZeroInt(),
		Liabilities:                    sdk.ZeroInt(),
		BorrowInterestPaidCollateral:   sdk.ZeroInt(),
		BorrowInterestPaidCustody:      sdk.ZeroInt(),
		BorrowInterestUnpaidCollateral: sdk.ZeroInt(),
		Custody:                        sdk.ZeroInt(),
		TakeProfitLiabilities:          sdk.ZeroInt(),
		TakeProfitCustody:              sdk.ZeroInt(),
		Leverage:                       leverage,
		MtpHealth:                      sdk.ZeroDec(),
		Position:                       position,
		AmmPoolId:                      poolId,
		ConsolidateLeverage:            leverage,
		SumCollateral:                  sdk.ZeroInt(),
		TakeProfitPrice:                takeProfitPrice,
		TakeProfitBorrowRate:           sdk.OneDec(),
		FundingFeePaidCollateral:       sdk.ZeroInt(),
		FundingFeePaidCustody:          sdk.ZeroInt(),
		FundingFeeReceivedCollateral:   sdk.ZeroInt(),
		FundingFeeReceivedCustody:      sdk.ZeroInt(),
	}
}

func (mtp MTP) Validate() error {
	if mtp.CollateralAsset == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no collateral asset specified")
	}
	if mtp.CustodyAsset == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no custody asset specified")
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
	key := append([]byte("margin_custody"), sdk.Uint64ToBigEndian(positionId)...)
	return address.Module(ModuleName, key)
}
