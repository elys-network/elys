package types

import (
	"errors"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

var EdenBoostApr = sdk.NewDec(1)

// NewParams creates a new Params instance
func NewParams(
	lpIncentives *IncentiveInfo,
	stkIncentives *IncentiveInfo,
	rewardPortionForLps sdk.Dec,
	rewardPortionForStakers sdk.Dec,
	poolInfos []PoolInfo,
	elysStakeSnapInterval int64,
	dexRewardsStakers DexRewardsTracker,
	dexRewardsLps DexRewardsTracker,
	maxEdenRewardAprStakers sdk.Dec,
	maxEdenRewardAprLps sdk.Dec,
	distributionInterval int64,
) Params {
	return Params{
		LpIncentives:            lpIncentives,
		StakeIncentives:         stkIncentives,
		RewardPortionForLps:     rewardPortionForLps,
		RewardPortionForStakers: rewardPortionForStakers,
		PoolInfos:               poolInfos,
		ElysStakeSnapInterval:   elysStakeSnapInterval,
		DexRewardsStakers:       dexRewardsStakers,
		DexRewardsLps:           dexRewardsLps,
		MaxEdenRewardAprStakers: maxEdenRewardAprStakers,
		MaxEdenRewardAprLps:     maxEdenRewardAprLps,
		DistributionInterval:    distributionInterval,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		nil,
		nil,
		sdk.NewDecWithPrec(60, 2),
		sdk.NewDecWithPrec(30, 2),
		[]PoolInfo(nil),
		10,
		DexRewardsTracker{
			NumBlocks:                     sdk.NewInt(1),
			Amount:                        sdk.ZeroDec(),
			AmountCollectedByOtherTracker: sdk.ZeroDec(),
		},
		DexRewardsTracker{
			NumBlocks:                     sdk.NewInt(1),
			Amount:                        sdk.ZeroDec(),
			AmountCollectedByOtherTracker: sdk.ZeroDec(),
		},
		sdk.NewDecWithPrec(3, 1),
		sdk.NewDecWithPrec(5, 1),
		10,
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

	if err := validateLPIncentives(p.LpIncentives); err != nil {
		return err
	}

	if err := validateStakeIncentives(p.StakeIncentives); err != nil {
		return err
	}

	if err := validatePoolInfos(p.PoolInfos); err != nil {
		return err
	}

	if err := validateElysStakeSnapInterval(p.ElysStakeSnapInterval); err != nil {
		return err
	}

	if err := validateDexRewardsStakers(p.DexRewardsStakers); err != nil {
		return err
	}

	if err := validateDexRewardsLps(p.DexRewardsLps); err != nil {
		return err
	}

	if err := validateDistributionInterval(p.DistributionInterval); err != nil {
		return err
	}

	if p.RewardPortionForLps.Add(p.RewardPortionForStakers).GT(sdk.NewDec(1)) {
		return errors.New("invalid rewards portion parameter")
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

func validateRewardPortionForStakers(i interface{}) error {
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

	if vv.EpochNumBlocks.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid number of blocks in epoch: %v", vv)
	}

	if vv.DistributionEpochInBlocks.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid distribution epoch in blocks: %v", vv)
	}

	if vv.CurrentEpochInBlocks.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid current epoch: %v", vv)
	}

	if vv.DistributionStartBlock.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid distribution epoch: %v", vv)
	}

	return nil
}

func validateStakeIncentives(i interface{}) error {
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

	if vv.EpochNumBlocks.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid number of blocks in epoch: %v", vv)
	}

	if vv.DistributionEpochInBlocks.LT(sdk.NewInt(0)) {
		return fmt.Errorf("invalid distribution epoch in blocks: %v", vv)
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

func validateElysStakeSnapInterval(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDexRewardsStakers(i interface{}) error {
	_, ok := i.(DexRewardsTracker)
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

func validateEdenRewardApr(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDistributionInterval(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
