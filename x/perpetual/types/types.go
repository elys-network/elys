package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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

func NewMTP(ctx sdk.Context, signer, collateralAsset, tradingAsset, liabilitiesAsset, custodyAsset string, position Position, takeProfitPrice sdk.Dec, poolId uint64) *MTP {
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
		FundingFeePaidCustody:         sdk.ZeroInt(),
		FundingFeeReceivedCustody:     sdk.ZeroInt(),
		OpenPrice:                     sdk.ZeroDec(),
		StopLossPrice:                 math.LegacyZeroDec(),
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
	openPrice := math.LegacyZeroDec()
	// open price = (collateral + liabilities) / custody
	if mtp.Position == Position_LONG {
		if mtp.CollateralAsset == mtp.TradingAsset {
			openPrice = mtp.Liabilities.ToLegacyDec().Quo(mtp.Custody.Sub(mtp.Collateral).ToLegacyDec())
		} else {
			openPrice = (mtp.Collateral.Add(mtp.Liabilities)).ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())
		}
	} else {
		// open price = (custody - collateral) / liabilities
		openPrice = (mtp.Custody.Sub(mtp.Collateral)).ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
	}
	mtp.OpenPrice = openPrice
	return
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
