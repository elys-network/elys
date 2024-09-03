package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(
	lpIncentives *IncentiveInfo,
	rewardPortionForLps sdk.Dec,
	rewardPortionForStakers sdk.Dec,
	maxEdenRewardAprLps sdk.Dec,
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
		sdk.NewDecWithPrec(60, 2),
		sdk.NewDecWithPrec(25, 2),
		sdk.NewDecWithPrec(5, 1),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateRewardPortionForLps(p.RewardPortionForLps); err != nil {
		return err
	}

	if err := validateLPIncentives(p.LpIncentives); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateRewardPortionForLps(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("reward percent for lp must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("reward percent for lp must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("reward percent for lp too large: %s", v)
	}

	return nil
}

func validateLPIncentives(i interface{}) error {
	vv, ok := i.(*IncentiveInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if vv == nil {
		return nil
	}

	if vv.EdenAmountPerYear.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("invalid eden amount per year: %v", vv)
	}

	if vv.BlocksDistributed.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid BlocksDistributed: %v", vv)
	}

	if vv.DistributionStartBlock.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid DistributionStartBlock: %v", vv)
	}

	return nil
}

func validateDexRewardsLps(i interface{}) error {
	_, ok := i.(DexRewardsTracker)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
