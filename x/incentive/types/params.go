package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter keys
var (
	ParamStoreKeyCommunityTax               = []byte("communitytax")
	ParamStoreKeyWithdrawAddrEnabled        = []byte("withdrawaddrenabled")
	ParamStoreKeyRewardPortionForLps        = []byte("rewardportionforlps")
	ParamStoreKeyLPIncentives               = []byte("lpincentives")
	ParamStoreKeyStkIncentives              = []byte("stkincentives")
	ParamStoreKeyPoolInfos                  = []byte("poolinfos")
	ParamStoreKeyElysStakeTrackingRate      = []byte("elysstaketrackingrate")
	ParamStoreKeyDexRewardsStakers          = []byte("dexrewardsstakers")
	ParamStoreKeyDexRewardsLps              = []byte("dexrewardslps")
	ParamStoreKeyMaxEdenRewardAprForStakers = []byte("maxedenrewardaprstakers")
	ParamStoreKeyMaxEdenRewardAprForLPs     = []byte("maxedenrewardaprlps")
	ParamStoreKeyDistributionEpochLPs       = []byte("distributionepochlps")
	ParamStoreKeyDistributionEpochStakers   = []byte("distributionepochstakers")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		LpIncentives:          []IncentiveInfo(nil),
		StakeIncentives:       []IncentiveInfo(nil),
		CommunityTax:          sdk.NewDecWithPrec(2, 2), // 2%
		WithdrawAddrEnabled:   true,
		RewardPortionForLps:   sdk.NewDecWithPrec(65, 2),
		PoolInfos:             []PoolInfo(nil),
		ElysStakeTrackingRate: 10,
		DexRewardsStakers: DexRewardsTracker{
			NumBlocks:                     sdk.ZeroInt(),
			Amount:                        sdk.ZeroDec(),
			AmountCollectedByOtherTracker: sdk.ZeroDec(),
		},
		DexRewardsLps: DexRewardsTracker{
			NumBlocks:                     sdk.ZeroInt(),
			Amount:                        sdk.ZeroDec(),
			AmountCollectedByOtherTracker: sdk.ZeroDec(),
		},
		MaxEdenRewardAprStakers:             sdk.NewDecWithPrec(3, 1),
		MaxEdenRewardAprLps:                 sdk.NewDecWithPrec(3, 1),
		DistributionEpochForStakersInBlocks: 10,
		DistributionEpochForLpsInBlocks:     10,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyCommunityTax, &p.CommunityTax, validateCommunityTax),
		paramtypes.NewParamSetPair(ParamStoreKeyWithdrawAddrEnabled, &p.WithdrawAddrEnabled, validateWithdrawAddrEnabled),
		paramtypes.NewParamSetPair(ParamStoreKeyRewardPortionForLps, &p.RewardPortionForLps, validateRewardPortionForLps),
		paramtypes.NewParamSetPair(ParamStoreKeyLPIncentives, &p.LpIncentives, validateLPIncentives),
		paramtypes.NewParamSetPair(ParamStoreKeyStkIncentives, &p.StakeIncentives, validateStakeIncentives),
		paramtypes.NewParamSetPair(ParamStoreKeyPoolInfos, &p.PoolInfos, validatePoolInfos),
		paramtypes.NewParamSetPair(ParamStoreKeyElysStakeTrackingRate, &p.ElysStakeTrackingRate, validateElysStakeTrakcingRate),
		paramtypes.NewParamSetPair(ParamStoreKeyDexRewardsStakers, &p.DexRewardsStakers, validateDexRewardsStakers),
		paramtypes.NewParamSetPair(ParamStoreKeyDexRewardsLps, &p.DexRewardsLps, validateDexRewardsLps),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxEdenRewardAprForStakers, &p.MaxEdenRewardAprStakers, validateEdenRewardApr),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxEdenRewardAprForLPs, &p.MaxEdenRewardAprLps, validateEdenRewardApr),
		paramtypes.NewParamSetPair(ParamStoreKeyDistributionEpochLPs, &p.DistributionEpochForLpsInBlocks, validateDistributionEpochLps),
		paramtypes.NewParamSetPair(ParamStoreKeyDistributionEpochStakers, &p.DistributionEpochForStakersInBlocks, validateDistributionEpochStakers),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	if err := validateCommunityTax(p.CommunityTax); err != nil {
		return err
	}

	if err := validateWithdrawAddrEnabled(p.WithdrawAddrEnabled); err != nil {
		return err
	}

	if err := validateRewardPortionForLps(p.RewardPortionForLps); err != nil {
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

	if err := validateElysStakeTrakcingRate(p.ElysStakeTrackingRate); err != nil {
		return err
	}

	if err := validateDexRewardsStakers(p.DexRewardsStakers); err != nil {
		return err
	}

	if err := validateDexRewardsLps(p.DexRewardsLps); err != nil {
		return err
	}

	if err := validateDistributionEpochLps(p.DistributionEpochForLpsInBlocks); err != nil {
		return err
	}

	if err := validateDistributionEpochStakers(p.DistributionEpochForStakersInBlocks); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ValidateBasic performs basic validation on distribution parameters.
func (p Params) ValidateBasic() error {
	if p.CommunityTax.IsNegative() || p.CommunityTax.GT(sdk.OneDec()) {
		return fmt.Errorf(
			"community tax should be non-negative and less than one: %s", p.CommunityTax,
		)
	}

	return nil
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("community tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("community tax must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("community tax too large: %s", v)
	}

	return nil
}

func validateWithdrawAddrEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
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
	v, ok := i.([]IncentiveInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == nil {
		return nil
	}

	for _, vv := range v {
		if vv.EdenAmountPerYear.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("invalid eden amount per year: %v", vv)
		}

		if vv.TotalBlocksPerYear.LT(sdk.NewInt(1)) {
			return fmt.Errorf("invalid total blocks per year: %v", vv)
		}

		if vv.AllocationEpochInBlocks.LT(sdk.NewInt(0)) {
			return fmt.Errorf("invalid allocation epoch in blocks: %v", vv)
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

		if vv.EdenBoostApr.GT(sdk.NewDec(1)) || vv.EdenBoostApr.LT(sdk.ZeroDec()) {
			return fmt.Errorf("invalid eden boot apr: %v", vv)
		}
	}

	return nil
}

func validateStakeIncentives(i interface{}) error {
	v, ok := i.([]IncentiveInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == nil {
		return nil
	}

	for _, vv := range v {
		if vv.EdenAmountPerYear.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("invalid eden amount per year: %v", vv)
		}

		if vv.TotalBlocksPerYear.LT(sdk.NewInt(1)) {
			return fmt.Errorf("invalid total blocks per year: %v", vv)
		}

		if vv.AllocationEpochInBlocks.LT(sdk.NewInt(0)) {
			return fmt.Errorf("invalid allocation epoch in blocks: %v", vv)
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

		if vv.EdenBoostApr.GT(sdk.NewDec(1)) || vv.EdenBoostApr.LT(sdk.ZeroDec()) {
			return fmt.Errorf("invalid eden boot apr: %v", vv)
		}
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

func validateElysStakeTrakcingRate(i interface{}) error {
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

func validateDistributionEpochLps(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDistributionEpochStakers(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
