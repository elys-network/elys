package types

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(
	minCommissionRate sdkmath.LegacyDec,
	maxVotingPower sdkmath.LegacyDec,
	minSelfDelegation sdkmath.Int,
	brokerAddress string,
	totalBlocksPerYear int64,
	rewardsDataLifeTime int64,
) Params {
	return Params{
		MinCommissionRate:   minCommissionRate,
		MaxVotingPower:      maxVotingPower,
		MinSelfDelegation:   minSelfDelegation,
		BrokerAddress:       brokerAddress,
		TotalBlocksPerYear:  totalBlocksPerYear,
		RewardsDataLifetime: rewardsDataLifeTime,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		sdkmath.LegacyNewDecWithPrec(5, 2),
		sdkmath.LegacyNewDec(100),
		sdkmath.OneInt(),
		authtypes.NewModuleAddress("zero").String(),
		6307200,
		86400, // 1 day
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMinCommissionRate(p.MinCommissionRate); err != nil {
		return err
	}
	if err := validateMaxVotingPower(p.MaxVotingPower); err != nil {
		return err
	}
	if err := validateMinSelfDelegation(p.MinSelfDelegation); err != nil {
		return err
	}
	if err := validateBrokerAddress(p.BrokerAddress); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMinCommissionRate(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return ErrInvalidMinCommissionRate
	}
	if v.IsNegative() {
		return ErrInvalidMinCommissionRate
	}
	return nil
}

func validateMaxVotingPower(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return ErrInvalidMaxVotingPower
	}
	if v.IsNegative() {
		return ErrInvalidMaxVotingPower
	}
	return nil
}

func validateMinSelfDelegation(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return ErrInvalidMinSelfDelegation
	}
	if v.IsNegative() {
		return ErrInvalidMinSelfDelegation
	}
	return nil
}

func validateBrokerAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return ErrInvalidBrokerAddress
	}
	if v == "" {
		return ErrInvalidBrokerAddress
	}
	return nil
}

// String implements the Stringer interface.
func (p LegacyParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
