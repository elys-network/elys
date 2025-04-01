package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func DefaultParams() Params {
	return Params{
		DerivativeMarketInstantListingFee:            sdk.NewCoin(ptypes.Elys, math.NewInt(0)),
		DefaultDerivativeMakerFeeRate:                utils.ZeroDec,
		DefaultDerivativeTakerFeeRate:                utils.ZeroDec,
		DefaultInitialMarginRatio:                    math.NewDecWithExp(5, -2),
		DefaultMaintenanceMarginRatio:                math.NewDecWithExp(1, -1),
		DefaultFundingInterval:                       10000,
		FundingMultiple:                              100,
		RelayerFeeShareRate:                          utils.ZeroDec,
		DefaultHourlyFundingRateCap:                  math.NewDecWithExp(1, -1),
		DefaultHourlyInterestRate:                    math.NewDecWithExp(1, -1),
		MaxDerivativeOrderSideCount:                  10,
		TradingRewardsVestingDuration:                10000,
		LiquidatorRewardShareRate:                    utils.ZeroDec,
		DerivativeAtomicMarketOrderFeeMultiplier:     utils.ZeroDec,
		MinimalProtocolFeeRate:                       utils.ZeroDec,
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
