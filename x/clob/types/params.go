package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func DefaultParams() Params {
	return Params{
		DerivativeMarketInstantListingFee:        sdk.NewCoin(ptypes.Elys, math.NewInt(0)),
		DefaultDerivativeMakerFeeRate:            math.LegacyZeroDec(),
		DefaultDerivativeTakerFeeRate:            math.LegacyZeroDec(),
		DefaultInitialMarginRatio:                math.LegacyNewDecWithPrec(5, 2),
		DefaultMaintenanceMarginRatio:            math.LegacyNewDecWithPrec(1, 1),
		DefaultFundingInterval:                   10000,
		FundingMultiple:                          100,
		RelayerFeeShareRate:                      math.LegacyZeroDec(),
		DefaultHourlyFundingRateCap:              math.LegacyNewDecWithPrec(1, 1),
		DefaultHourlyInterestRate:                math.LegacyNewDecWithPrec(1, 1),
		MaxDerivativeOrderSideCount:              10,
		TradingRewardsVestingDuration:            10000,
		LiquidatorRewardShareRate:                math.LegacyZeroDec(),
		DerivativeAtomicMarketOrderFeeMultiplier: math.LegacyZeroDec(),
		MinimalProtocolFeeRate:                   math.LegacyZeroDec(),
		ExchangeAdmins:                           nil,
		EpochIdentifier:                          epochstypes.EightHoursEpochID,
	}
}

func (p Params) Validate() error {
	return nil
}
