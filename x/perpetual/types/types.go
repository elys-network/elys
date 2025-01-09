package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func NewMTP(ctx sdk.Context, signer, collateralAsset, tradingAsset, liabilitiesAsset, custodyAsset string, position Position, takeProfitPrice sdkmath.LegacyDec, poolId uint64) *MTP {
	return &MTP{
		Address:                       signer,
		CollateralAsset:               collateralAsset,
		TradingAsset:                  tradingAsset,
		LiabilitiesAsset:              liabilitiesAsset,
		CustodyAsset:                  custodyAsset,
		Collateral:                    sdkmath.ZeroInt(),
		Liabilities:                   sdkmath.ZeroInt(),
		BorrowInterestPaidCustody:     sdkmath.ZeroInt(),
		BorrowInterestUnpaidLiability: sdkmath.ZeroInt(),
		Custody:                       sdkmath.ZeroInt(),
		TakeProfitLiabilities:         sdkmath.ZeroInt(),
		TakeProfitCustody:             sdkmath.ZeroInt(),
		MtpHealth:                     sdkmath.LegacyZeroDec(),
		Position:                      position,
		Id:                            0,
		AmmPoolId:                     poolId,
		TakeProfitPrice:               takeProfitPrice,
		TakeProfitBorrowFactor:        sdkmath.LegacyOneDec(),
		FundingFeePaidCustody:         sdkmath.ZeroInt(),
		FundingFeeReceivedCustody:     sdkmath.ZeroInt(),
		OpenPrice:                     sdkmath.LegacyZeroDec(),
		StopLossPrice:                 sdkmath.LegacyZeroDec(),
		LastInterestCalcTime:          uint64(ctx.BlockTime().Unix()),
		LastInterestCalcBlock:         uint64(ctx.BlockHeight()),
		LastFundingCalcTime:           uint64(ctx.BlockTime().Unix()),
		LastFundingCalcBlock:          uint64(ctx.BlockHeight()),
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

func (mtp *MTP) GetAndSetOpenPrice() {
	openPrice := sdkmath.LegacyZeroDec()
	if mtp.Position == Position_LONG {
		if mtp.CollateralAsset == mtp.TradingAsset {
			// open price = liabilities / (custody - collateral)
			denominator := mtp.Custody.Sub(mtp.Collateral).ToLegacyDec()
			if !denominator.IsZero() {
				openPrice = mtp.Liabilities.ToLegacyDec().Quo(denominator)
			}
		} else {
			// open price = (collateral + liabilities) / custody
			openPrice = (mtp.Collateral.Add(mtp.Liabilities)).ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())
		}
	} else {
		if mtp.Liabilities.IsZero() {
			mtp.OpenPrice = openPrice
		} else {
			// open price = (custody - collateral) / liabilities
			openPrice = (mtp.Custody.Sub(mtp.Collateral)).ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
		}
	}
	mtp.OpenPrice = openPrice
	return
}

func (mtp MTP) GetAccountAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(mtp.Address)
}

func (mtp MTP) GetBorrowInterestAmountAsCustodyAsset(tradingAssetPrice sdkmath.LegacyDec) (sdkmath.Int, error) {
	borrowInterestPaymentInCustody := sdkmath.ZeroInt()
	if mtp.Position == Position_LONG {
		if tradingAssetPrice.IsZero() {
			return sdkmath.ZeroInt(), errors.New("trading asset price is zero")
		}
		// liabilities are in usdc, custody is in trading asset
		borrowInterestPaymentInCustody = mtp.BorrowInterestUnpaidLiability.ToLegacyDec().Quo(tradingAssetPrice).TruncateInt()
	} else {
		// liabilities are in trading asset, custody is in usdc
		borrowInterestPaymentInCustody = mtp.BorrowInterestUnpaidLiability.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
	}
	return borrowInterestPaymentInCustody, nil
}
