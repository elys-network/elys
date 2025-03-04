package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func DefaultParams() Params {
	return Params{
		DerivativeMarketInstantListingFee:            sdk.NewCoin(ptypes.Elys, math.NewInt(0)),
		DefaultDerivativeMakerFeeRate:                math.LegacyZeroDec(),
		DefaultDerivativeTakerFeeRate:                math.LegacyZeroDec(),
		DefaultInitialMarginRatio:                    math.LegacyMustNewDecFromStr("0.05"),
		DefaultMaintenanceMarginRatio:                math.LegacyMustNewDecFromStr("0.1"),
		DefaultFundingInterval:                       10000,
		FundingMultiple:                              100,
		RelayerFeeShareRate:                          math.LegacyZeroDec(),
		DefaultHourlyFundingRateCap:                  math.LegacyMustNewDecFromStr("0.1"),
		DefaultHourlyInterestRate:                    math.LegacyMustNewDecFromStr("0.1"),
		MaxDerivativeOrderSideCount:                  10,
		TradingRewardsVestingDuration:                10000,
		LiquidatorRewardShareRate:                    math.LegacyZeroDec(),
		DerivativeAtomicMarketOrderFeeMultiplier:     math.LegacyZeroDec(),
		MinimalProtocolFeeRate:                       math.LegacyZeroDec(),
		IsInstantDerivativeMarketLaunchEnabled:       true,
		PostOnlyModeHeightThreshold:                  0,
		MarginDecreasePriceTimestampThresholdSeconds: 0,
		ExchangeAdmins:                               nil,
		MaxSubAccounts:                               3,
	}
}

func (p Params) Validate() error {
	return nil
}
