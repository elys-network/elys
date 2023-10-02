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
	KeyDepositDenom   = []byte("DepositDenom")
	KeyRedemptionRate = []byte("RedemptionRate")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(depositDenom string, redemptionRate sdk.Dec) Params {
	return Params{
		DepositDenom:   depositDenom,
		RedemptionRate: redemptionRate,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams("uusdc", sdk.OneDec())
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDepositDenom, &p.DepositDenom, validateDepositDenom),
		paramtypes.NewParamSetPair(KeyRedemptionRate, &p.RedemptionRate, validateRedemptionRate),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateDepositDenom(p.DepositDenom); err != nil {
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
