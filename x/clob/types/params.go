package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/elys-network/elys/v6/x/epochs/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
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
	// Validate fee rates
	if p.DefaultDerivativeMakerFeeRate.IsNegative() {
		return fmt.Errorf("maker fee rate cannot be negative")
	}
	if p.DefaultDerivativeTakerFeeRate.IsNegative() {
		return fmt.Errorf("taker fee rate cannot be negative")
	}

	// Validate margin ratios
	if p.DefaultInitialMarginRatio.IsNil() || p.DefaultInitialMarginRatio.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("initial margin ratio must be positive")
	}
	if p.DefaultMaintenanceMarginRatio.IsNil() || p.DefaultMaintenanceMarginRatio.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("maintenance margin ratio must be positive")
	}
	if p.DefaultMaintenanceMarginRatio.GTE(p.DefaultInitialMarginRatio) {
		return fmt.Errorf("maintenance margin ratio must be less than initial margin ratio")
	}

	// Validate other parameters
	if p.DefaultFundingInterval == 0 {
		return fmt.Errorf("funding interval must be positive")
	}
	if p.MaxDerivativeOrderSideCount == 0 {
		return fmt.Errorf("max derivative order side count must be positive")
	}

	// Validate exchange admins
	for _, admin := range p.ExchangeAdmins {
		if _, err := sdk.AccAddressFromBech32(admin); err != nil {
			return fmt.Errorf("invalid exchange admin address %s: %v", admin, err)
		}
	}

	// Validate fee shares and rates
	if p.RelayerFeeShareRate.IsNegative() || p.RelayerFeeShareRate.GT(math.LegacyOneDec()) {
		return fmt.Errorf("relayer fee share rate must be between 0 and 1")
	}
	if p.LiquidatorRewardShareRate.IsNegative() || p.LiquidatorRewardShareRate.GT(math.LegacyOneDec()) {
		return fmt.Errorf("liquidator reward share rate must be between 0 and 1")
	}

	// Validate listing fee
	if err := p.DerivativeMarketInstantListingFee.Validate(); err != nil {
		return fmt.Errorf("invalid listing fee: %v", err)
	}

	// Validate epoch identifier
	if p.EpochIdentifier == "" {
		return fmt.Errorf("epoch identifier cannot be empty")
	}

	return nil
}
