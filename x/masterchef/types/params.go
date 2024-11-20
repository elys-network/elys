package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(
	lpIncentives *IncentiveInfo,
	rewardPortionForLps sdkmath.LegacyDec,
	rewardPortionForStakers sdkmath.LegacyDec,
	maxEdenRewardAprLps sdkmath.LegacyDec,
	protocolRevenueAddress string,
) Params {
	return Params{
		LpIncentives:            lpIncentives,
		RewardPortionForLps:     rewardPortionForLps,
		RewardPortionForStakers: rewardPortionForStakers,
		MaxEdenRewardAprLps:     maxEdenRewardAprLps,
		SupportedRewardDenoms:   nil,
		ProtocolRevenueAddress:  protocolRevenueAddress,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		nil,
		sdkmath.LegacyNewDecWithPrec(60, 2),
		sdkmath.LegacyNewDecWithPrec(25, 2),
		sdkmath.LegacyNewDecWithPrec(5, 1),
		authtypes.NewModuleAddress("protocol-revenue-address").String(), // TODO: Change it in genesis in mainnet launch
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateRewardPortionForLps(p.RewardPortionForLps); err != nil {
		return err
	}

	if err := validateRewardPortionForStakers(p.RewardPortionForStakers); err != nil {
		return err
	}

	if err := validateMaxEdenRewardAprLps(p.MaxEdenRewardAprLps); err != nil {
		return err
	}

	if err := validateLPIncentives(p.LpIncentives); err != nil {
		return err
	}

	if err := validateProtocolRevenueAddress(p.ProtocolRevenueAddress); err != nil {
		return err
	}

	if err := validateSupportedRewardDenoms(p.SupportedRewardDenoms); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// String implements the Stringer interface.
func (p LegacyParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateRewardPortionForLps(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("reward percent for lp must not be nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("reward percent for lp must be positive: %s", v)
	}
	if v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("reward percent for lp too large: %s", v)
	}

	return nil
}

func validateRewardPortionForStakers(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("reward percent for stakers must not be nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("reward percent for stakers must be positive: %s", v)
	}
	if v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("reward percent for stakers too large: %s", v)
	}

	return nil
}

func validateLPIncentives(i interface{}) error {
	vv, ok := i.(*IncentiveInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// not checking for nil as LPIncentives nil is allowed in abci.go
	if vv != nil && vv.EdenAmountPerYear.LTE(sdkmath.ZeroInt()) {
		return fmt.Errorf("invalid eden amount per year: %v", vv)
	}

	if vv != nil && vv.BlocksDistributed < 0 {
		return fmt.Errorf("invalid BlocksDistributed: %v", vv)
	}

	return nil
}

func validateMaxEdenRewardAprLps(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("MaxEdenRewardAprLps must not be nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("MaxEdenRewardAprLps must be positive: %s", v)
	}

	return nil
}

func validateProtocolRevenueAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("ProtocolRevenueAddres cannot be empty")
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid ProtocolRevenueAddress %s: %v", v, err)
	}

	return nil
}

func validateSupportedRewardDenoms(i interface{}) error {
	v, ok := i.([]*SupportedRewardDenom)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, supportedRewardDenom := range v {
		if err := sdk.ValidateDenom(supportedRewardDenom.Denom); err != nil {
			return fmt.Errorf("invalid reward denom %s: %v", supportedRewardDenom.Denom, err)
		}
		if supportedRewardDenom.MinAmount.IsNil() {
			return fmt.Errorf("reward denom minimum amount cannot be nil: %s", supportedRewardDenom.Denom)
		}
		if supportedRewardDenom.MinAmount.IsNegative() {
			return fmt.Errorf("reward denom(%s) minimum amount cannot be negative: %s", supportedRewardDenom.Denom, supportedRewardDenom.MinAmount.String())
		}
	}

	return nil
}
