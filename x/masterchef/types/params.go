package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	lpIncentives *IncentiveInfo,
	rewardPortionForLps sdk.Dec,
	poolInfos []PoolInfo,
	dexRewardsLps DexRewardsTracker,
	maxEdenRewardAprLps sdk.Dec,
) Params {
	return Params{
		LpIncentives:          lpIncentives,
		RewardPortionForLps:   rewardPortionForLps,
		PoolInfos:             poolInfos,
		DexRewardsLps:         dexRewardsLps,
		MaxEdenRewardAprLps:   maxEdenRewardAprLps,
		SupportedRewardDenoms: nil,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		nil,
		sdk.NewDecWithPrec(60, 2),
		[]PoolInfo(nil),
		DexRewardsTracker{
			NumBlocks: sdk.NewInt(1),
			Amount:    sdk.ZeroDec(),
		},
		sdk.NewDecWithPrec(5, 1),
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateRewardPortionForLps(p.RewardPortionForLps); err != nil {
		return err
	}

	if err := validateLPIncentives(p.LpIncentives); err != nil {
		return err
	}

	if err := validatePoolInfos(p.PoolInfos); err != nil {
		return err
	}

	if err := validateDexRewardsLps(p.DexRewardsLps); err != nil {
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

	if vv.TotalBlocksPerYear.LT(sdk.NewInt(1)) {
		return fmt.Errorf("invalid total blocks per year: %v", vv)
	}

	if vv.CurrentEpochInBlocks.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid current epoch: %v", vv)
	}

	if vv.DistributionStartBlock.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid distribution epoch: %v", vv)
	}

	return nil
}

func validatePoolInfos(i interface{}) error {
	_, ok := i.([]PoolInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
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
