package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	etypes "github.com/elys-network/elys/x/epochs/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter keys
var (
	ParamStoreKeyCommunityTax          = []byte("communitytax")
	ParamStoreKeyWithdrawAddrEnabled   = []byte("withdrawaddrenabled")
	ParamStoreKeyRewardPortionForLps   = []byte("rewardportionforlps")
	ParamStoreKeyLPIncentives          = []byte("lpincentives")
	ParamStoreKeyStkIncentives         = []byte("stkincentives")
	ParamStoreKeyPoolInfos             = []byte("poolinfos")
	ParamStoreKeyElysStakeTrackingRate = []byte("elysstaketrackingrate")
	ParamStoreKeyAprEdenElys           = []byte("apredenelys")
	ParamStoreKeyAprUsdc               = []byte("aprusdc")
	ParamStoreKeyAprEdenB              = []byte("apredenb")
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
		AprEdenElys: AprEdenElys{
			Uusdc:  sdk.ZeroInt(),
			Ueden:  sdk.ZeroInt(),
			Uedenb: sdk.NewInt(100),
		},
		AprUsdc: AprUsdc{
			Uusdc: sdk.ZeroInt(),
			Ueden: sdk.ZeroInt(),
		},
		AprEdenB: AprEdenB{
			Uedenb: sdk.NewInt(100),
		},
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
		paramtypes.NewParamSetPair(ParamStoreKeyAprEdenElys, &p.AprEdenElys, validateAprEdenElys),
		paramtypes.NewParamSetPair(ParamStoreKeyAprUsdc, &p.AprUsdc, validateAprUsdc),
		paramtypes.NewParamSetPair(ParamStoreKeyAprEdenB, &p.AprEdenB, validateAprEdenB),
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

	if err := validateAprEdenElys(p.AprEdenElys); err != nil {
		return err
	}

	if err := validateAprUsdc(p.AprEdenElys); err != nil {
		return err
	}

	if err := validateAprEdenB(p.AprEdenElys); err != nil {
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
		if vv.Amount.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("invalid amount: %v", vv)
		}

		if vv.EpochIdentifier != etypes.WeekEpochID &&
			vv.EpochIdentifier != etypes.DayEpochID &&
			vv.EpochIdentifier != etypes.HourEpochID {
			return fmt.Errorf("invalid epoch: %v", vv)
		}

		if vv.NumEpochs < 1 {
			return fmt.Errorf("invalid num epoch: %v", vv)
		}

		if vv.CurrentEpoch < 0 {
			return fmt.Errorf("invalid current epoch: %v", vv)
		}

		if vv.EdenBoostApr < 1 {
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
		if vv.Amount.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("invalid amount: %v", vv)
		}

		if vv.EpochIdentifier != etypes.WeekEpochID &&
			vv.EpochIdentifier != etypes.DayEpochID &&
			vv.EpochIdentifier != etypes.HourEpochID {
			return fmt.Errorf("invalid epoch: %v", vv)
		}

		if vv.NumEpochs < 1 {
			return fmt.Errorf("invalid num epoch: %v", vv)
		}

		if vv.CurrentEpoch < 0 {
			return fmt.Errorf("invalid current epoch: %v", vv)
		}

		if vv.EdenBoostApr < 1 {
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

func validateAprEdenElys(i interface{}) error {
	_, ok := i.(AprEdenElys)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateAprUsdc(i interface{}) error {
	_, ok := i.(AprUsdc)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateAprEdenB(i interface{}) error {
	_, ok := i.(AprEdenB)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
