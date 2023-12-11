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
	KeyDepositDenom         = []byte("DepositDenom")
	KeyRedemptionRate       = []byte("RedemptionRate")
	KeyEpochLength          = []byte("EpochLength")
	KeyInterestRate         = []byte("InterestRate")
	KeyInterestRateMax      = []byte("InterestRateMax")
	KeyInterestRateMin      = []byte("InterestRateMin")
	KeyInterestRateIncrease = []byte("InterestRateIncrease")
	KeyInterestRateDecrease = []byte("InterestRateDecrease")
	KeyHealthGainFactor     = []byte("HealthGainFactor")
	KeyTotalValue           = []byte("TotalValue")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	depositDenom string,
	redemptionRate sdk.Dec,
	epochLength int64,
	interestRate sdk.Dec,
	interestRateMax sdk.Dec,
	interestRateMin sdk.Dec,
	interestRateIncrease sdk.Dec,
	interestRateDecrease sdk.Dec,
	healthGainFactor sdk.Dec,
	totalValue sdk.Int,
) Params {
	return Params{
		DepositDenom:         depositDenom,
		RedemptionRate:       redemptionRate,
		EpochLength:          epochLength,
		InterestRate:         interestRate,
		InterestRateMax:      interestRateMax,
		InterestRateMin:      interestRateMin,
		InterestRateIncrease: interestRateIncrease,
		InterestRateDecrease: interestRateDecrease,
		HealthGainFactor:     healthGainFactor,
		TotalValue:           totalValue,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"uusdc",                   // deposit denom
		sdk.OneDec(),              // default redemption rate
		1,                         // epoch length
		sdk.NewDecWithPrec(15, 2), // 15% - default interest
		sdk.NewDecWithPrec(17, 2), // 17% - max
		sdk.NewDecWithPrec(12, 2), // 12% - min
		sdk.NewDecWithPrec(1, 2),  // 1% - interest rate increase
		sdk.NewDecWithPrec(1, 2),  // 1% - interest rate decrease
		sdk.NewDec(1),             // health gain factor
		sdk.NewInt(0),             // total value - 0
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDepositDenom, &p.DepositDenom, validateDepositDenom),
		paramtypes.NewParamSetPair(KeyRedemptionRate, &p.RedemptionRate, validateRedemptionRate),
		paramtypes.NewParamSetPair(KeyInterestRate, &p.InterestRate, validateInterestRate),
		paramtypes.NewParamSetPair(KeyInterestRateMax, &p.InterestRateMax, validateInterestRateMax),
		paramtypes.NewParamSetPair(KeyInterestRateMin, &p.InterestRateMin, validateInterestRateMin),
		paramtypes.NewParamSetPair(KeyInterestRateIncrease, &p.InterestRateIncrease, validateInterestRateIncrease),
		paramtypes.NewParamSetPair(KeyInterestRateDecrease, &p.InterestRateDecrease, validateInterestRateDecrease),
		paramtypes.NewParamSetPair(KeyHealthGainFactor, &p.HealthGainFactor, validateHealthGainFactor),
		paramtypes.NewParamSetPair(KeyEpochLength, &p.EpochLength, validateEpochLength),
		paramtypes.NewParamSetPair(KeyTotalValue, &p.TotalValue, validateTotalValue),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateDepositDenom(p.DepositDenom); err != nil {
		return err
	}

	if err := validateRedemptionRate(p.RedemptionRate); err != nil {
		return err
	}
	if err := validateInterestRate(p.InterestRate); err != nil {
		return err
	}
	if err := validateInterestRateMax(p.InterestRateMax); err != nil {
		return err
	}
	if err := validateInterestRateMin(p.InterestRateMin); err != nil {
		return err
	}
	if err := validateInterestRateIncrease(p.InterestRateIncrease); err != nil {
		return err
	}
	if err := validateInterestRateDecrease(p.InterestRateDecrease); err != nil {
		return err
	}
	if err := validateHealthGainFactor(p.HealthGainFactor); err != nil {
		return err
	}
	if err := validateEpochLength(p.EpochLength); err != nil {
		return err
	}
	if err := validateTotalValue(p.TotalValue); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateDepositDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("deposit denom should not be empty")
	}

	return nil
}

func validateRedemptionRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid redemption rate type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("redemption rate must be not nil")
	}
	if v.LT(sdk.OneDec()) {
		return fmt.Errorf("redemption rate must be bigger than 1: %s", v)
	}

	return nil
}

func validateEpochLength(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateInterestRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest must be positive: %s", v)
	}

	return nil
}

func validateInterestRateMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest max must be positive: %s", v)
	}

	return nil
}

func validateInterestRateMin(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest min must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest min must be positive: %s", v)
	}

	return nil
}

func validateInterestRateIncrease(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest rate increase must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest rate increase must be positive: %s", v)
	}

	return nil
}

func validateInterestRateDecrease(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest rate decrease must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest rate decrease must be positive: %s", v)
	}

	return nil
}

func validateHealthGainFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("health gain factor must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("health gain factor must be positive: %s", v)
	}

	return nil
}

func validateTotalValue(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("total value must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("total value must be positive: %s", v)
	}

	return nil
}
