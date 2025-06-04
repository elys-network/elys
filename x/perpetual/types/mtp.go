package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

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

func (mtp MTP) GetAccountAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(mtp.Address)
}

func (mtp MTP) GetBorrowInterestAmountAsCustodyAsset(tradingAssetPriceInBaseUnits osmomath.BigDec) (sdkmath.Int, error) {
	borrowInterestPaymentInCustody := sdkmath.ZeroInt()
	if mtp.Position == Position_LONG {
		if tradingAssetPriceInBaseUnits.IsZero() {
			return sdkmath.ZeroInt(), errors.New("trading asset price is zero in GetBorrowInterestAmountAsCustodyAsset")
		}
		// liabilities are in usdc, custody is in trading asset
		borrowInterestPaymentInCustody = mtp.GetBigDecBorrowInterestUnpaidLiability().Quo(tradingAssetPriceInBaseUnits).Dec().TruncateInt()
	} else {
		// liabilities are in trading asset, custody is in usdc
		borrowInterestPaymentInCustody = mtp.GetBigDecBorrowInterestUnpaidLiability().Mul(tradingAssetPriceInBaseUnits).Dec().TruncateInt()
	}
	return borrowInterestPaymentInCustody, nil
}

func (mtp MTP) CheckForStopLoss(tradingAssetPrice sdkmath.LegacyDec) bool {
	if mtp.StopLossPrice.IsNil() || mtp.StopLossPrice.IsZero() {
		return false
	}
	stopLossReached := false
	if mtp.Position == Position_LONG {
		stopLossReached = tradingAssetPrice.LTE(mtp.StopLossPrice)
	}
	if mtp.Position == Position_SHORT {
		stopLossReached = tradingAssetPrice.GTE(mtp.StopLossPrice)
	}
	return stopLossReached
}

func (mtp MTP) CheckForTakeProfit(tradingAssetPrice sdkmath.LegacyDec) bool {
	if mtp.TakeProfitPrice.IsNil() || mtp.TakeProfitPrice.IsZero() {
		return false
	}
	takeProfitReached := false
	if mtp.Position == Position_LONG {
		takeProfitReached = tradingAssetPrice.GTE(mtp.TakeProfitPrice)
	}
	if mtp.Position == Position_SHORT {
		takeProfitReached = tradingAssetPrice.LTE(mtp.TakeProfitPrice)
	}
	return takeProfitReached
}

func (mtp MTP) GetBigDecTakeProfitLiabilities() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(mtp.TakeProfitLiabilities)
}

func (mtp MTP) GetBigDecTakeProfitCustody() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(mtp.TakeProfitCustody)
}

func (mtp MTP) GetBigDecTakeProfitBorrowFactor() osmomath.BigDec {
	return osmomath.BigDecFromDec(mtp.TakeProfitBorrowFactor)
}

func (mtp MTP) GetBigDecTakeProfitPrice() osmomath.BigDec {
	return osmomath.BigDecFromDec(mtp.TakeProfitPrice)
}

func (mtp MTP) GetBigDecStopLossPrice() osmomath.BigDec {
	return osmomath.BigDecFromDec(mtp.StopLossPrice)
}

func (mtp MTP) GetBigDecBorrowInterestUnpaidLiability() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(mtp.BorrowInterestUnpaidLiability)
}

func (mtp MTP) GetBigDecLiabilities() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(mtp.Liabilities)
}

func (mtp MTP) GetBigDecCustody() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(mtp.Custody)
}

func (mtp MTP) GetBigDecCollateral() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(mtp.Collateral)
}

func (mtp MTP) GetBigDecOpenPrice() osmomath.BigDec {
	return osmomath.BigDecFromDec(mtp.OpenPrice)
}

func (mtp *MTP) UpdateMTPTakeProfitBorrowFactor() error {
	takeProfitBorrowFactor, err := mtp.CalcMTPTakeProfitBorrowFactor()
	if err != nil {
		return err
	}
	mtp.TakeProfitBorrowFactor = takeProfitBorrowFactor.Dec()
	return nil
}

func (mtp MTP) CalcMTPTakeProfitBorrowFactor() (osmomath.BigDec, error) {
	// Ensure mtp.Custody is not zero to avoid division by zero
	if mtp.Custody.IsZero() {
		return osmomath.ZeroBigDec(), ErrZeroCustodyAmount
	}

	// infinite for long, 0 for short
	if mtp.IsTakeProfitPriceInfinite() || mtp.TakeProfitPrice.IsZero() {
		return osmomath.OneBigDec(), nil
	}

	takeProfitBorrowFactor := osmomath.OneBigDec()
	if mtp.Position == Position_LONG {
		// takeProfitBorrowFactor = 1 - (liabilities / (custody * take profit price))
		takeProfitBorrowFactor = osmomath.OneBigDec().Sub(mtp.GetBigDecLiabilities().Quo(mtp.GetBigDecCustody().MulDec(mtp.TakeProfitPrice)))
	} else {
		// takeProfitBorrowFactor = 1 - ((liabilities  * take profit price) / custody)
		takeProfitBorrowFactor = osmomath.OneBigDec().Sub((mtp.GetBigDecLiabilities().MulDec(mtp.TakeProfitPrice)).Quo(mtp.GetBigDecCustody()))
	}

	return takeProfitBorrowFactor, nil
}

func (mtp MTP) IsTakeProfitPriceInfinite() bool {
	return mtp.TakeProfitPrice.Equal(TakeProfitPriceDefault)
}

func (mtp MTP) IsLong() bool {
	return mtp.Position == Position_LONG
}

func (mtp MTP) IsShort() bool {
	return mtp.Position == Position_SHORT
}
