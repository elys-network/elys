package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	epochtypes "github.com/elys-network/elys/x/epochs/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyLeverageMax                                    = []byte("LeverageMax")
	KeyBorrowInterestRateMax                          = []byte("BorrowInterestRateMax")
	KeyBorrowInterestRateMin                          = []byte("BorrowInterestRateMin")
	KeyBorrowInterestRateIncrease                     = []byte("BorrowInterestRateIncrease")
	KeyBorrowInterestRateDecrease                     = []byte("BorrowInterestRateDecrease")
	KeyHealthGainFactor                               = []byte("HealthGainFactor")
	KeyEpochLength                                    = []byte("EpochLength")
	KeyRemovalQueueThreshold                          = []byte("RemovalQueueThreshold")
	KeyMaxOpenPositions                               = []byte("MaxOpenPositions")
	KeyPoolOpenThreshold                              = []byte("PoolOpenThreshold")
	KeyForceCloseFundPercentage                       = []byte("ForceCloseFundPercentage")
	KeyForceCloseFundAddress                          = []byte("ForceCloseFundAddress")
	KeyIncrementalBorrowInterestPaymentFundPercentage = []byte("IncrementalBorrowInterestPaymentFundPercentage")
	KeyIncrementalBorrowInterestPaymentFundAddress    = []byte("IncrementalBorrowInterestPaymentFundAddress")
	KeySqModifier                                     = []byte("SqModifier")
	KeySafetyFactor                                   = []byte("SafetyFactor")
	KeyIncrementalBorrowInterestPaymentEnabled        = []byte("IncrementalBorrowInterestPaymentEnabled")
	KeyWhitelistingEnabled                            = []byte("WhitelistingEnabled")
	KeyInvariantCheckEpoch                            = []byte("InvariantCheckEpoch")
	KeyBrokerAddress                                  = []byte("BrokerAddress")
	KeyTakeProfitBorrowInterestRateMin                = []byte("TakeProfitBorrowInterestRateMin")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		TakeProfitBorrowInterestRateMin:                sdk.OneDec(),
		BorrowInterestRateDecrease:                     sdk.MustNewDecFromStr("0.000000000333333333"),
		BorrowInterestRateIncrease:                     sdk.MustNewDecFromStr("0.000000000333333333"),
		BorrowInterestRateMax:                          sdk.MustNewDecFromStr("0.000000270000000000"),
		BorrowInterestRateMin:                          sdk.MustNewDecFromStr("0.000000030000000000"),
		BrokerAddress:                                  "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l",
		EpochLength:                                    (int64)(1),
		ForceCloseFundAddress:                          "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l",
		ForceCloseFundPercentage:                       sdk.OneDec(),
		HealthGainFactor:                               sdk.MustNewDecFromStr("0.000000022000000000"),
		IncrementalBorrowInterestPaymentEnabled:        true,
		IncrementalBorrowInterestPaymentFundAddress:    "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l",
		IncrementalBorrowInterestPaymentFundPercentage: sdk.MustNewDecFromStr("0.350000000000000000"),
		InvariantCheckEpoch:                            epochtypes.DayEpochID,
		LeverageMax:                                    sdk.NewDec(10),
		MaxOpenPositions:                               (int64)(9999),
		PoolOpenThreshold:                              sdk.OneDec(),
		RemovalQueueThreshold:                          sdk.MustNewDecFromStr("0.350000000000000000"),
		SafetyFactor:                                   sdk.MustNewDecFromStr("1.050000000000000000"),
		SqModifier:                                     sdk.MustNewDecFromStr("10000000000000000000000000.000000000000000000"),
		WhitelistingEnabled:                            false,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLeverageMax, &p.LeverageMax, validateLeverageMax),
		paramtypes.NewParamSetPair(KeyBorrowInterestRateMax, &p.BorrowInterestRateMax, validateBorrowInterestRateMax),
		paramtypes.NewParamSetPair(KeyBorrowInterestRateMin, &p.BorrowInterestRateMin, validateBorrowInterestRateMin),
		paramtypes.NewParamSetPair(KeyBorrowInterestRateIncrease, &p.BorrowInterestRateIncrease, validateBorrowInterestRateIncrease),
		paramtypes.NewParamSetPair(KeyBorrowInterestRateDecrease, &p.BorrowInterestRateDecrease, validateBorrowInterestRateDecrease),
		paramtypes.NewParamSetPair(KeyHealthGainFactor, &p.HealthGainFactor, validateHealthGainFactor),
		paramtypes.NewParamSetPair(KeyEpochLength, &p.EpochLength, validateEpochLength),
		paramtypes.NewParamSetPair(KeyRemovalQueueThreshold, &p.RemovalQueueThreshold, validateRemovalQueueThreshold),
		paramtypes.NewParamSetPair(KeyMaxOpenPositions, &p.MaxOpenPositions, validateMaxOpenPositions),
		paramtypes.NewParamSetPair(KeyPoolOpenThreshold, &p.PoolOpenThreshold, validatePoolOpenThreshold),
		paramtypes.NewParamSetPair(KeyForceCloseFundPercentage, &p.ForceCloseFundPercentage, validateForceCloseFundPercentage),
		paramtypes.NewParamSetPair(KeyForceCloseFundAddress, &p.ForceCloseFundAddress, validateForceCloseFundAddress),
		paramtypes.NewParamSetPair(KeyIncrementalBorrowInterestPaymentFundPercentage, &p.IncrementalBorrowInterestPaymentFundPercentage, validateIncrementalBorrowInterestPaymentFundPercentage),
		paramtypes.NewParamSetPair(KeyIncrementalBorrowInterestPaymentFundAddress, &p.IncrementalBorrowInterestPaymentFundAddress, validateIncrementalBorrowInterestPaymentFundAddress),
		paramtypes.NewParamSetPair(KeySqModifier, &p.SqModifier, validateSqModifier),
		paramtypes.NewParamSetPair(KeySafetyFactor, &p.SafetyFactor, validateSafetyFactor),
		paramtypes.NewParamSetPair(KeyIncrementalBorrowInterestPaymentEnabled, &p.IncrementalBorrowInterestPaymentEnabled, validateIncrementalBorrowInterestPaymentEnabled),
		paramtypes.NewParamSetPair(KeyWhitelistingEnabled, &p.WhitelistingEnabled, validateWhitelistingEnabled),
		paramtypes.NewParamSetPair(KeyInvariantCheckEpoch, &p.InvariantCheckEpoch, validateInvariantCheckEpoch),
		paramtypes.NewParamSetPair(KeyBrokerAddress, &p.BrokerAddress, validateBrokerAddress),
		paramtypes.NewParamSetPair(KeyTakeProfitBorrowInterestRateMin, &p.TakeProfitBorrowInterestRateMin, validateTakeProfitBorrowInterestRateMin),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateLeverageMax(p.LeverageMax); err != nil {
		return err
	}
	if err := validateBorrowInterestRateMax(p.BorrowInterestRateMax); err != nil {
		return err
	}
	if err := validateBorrowInterestRateMin(p.BorrowInterestRateMin); err != nil {
		return err
	}
	if err := validateBorrowInterestRateIncrease(p.BorrowInterestRateIncrease); err != nil {
		return err
	}
	if err := validateBorrowInterestRateDecrease(p.BorrowInterestRateDecrease); err != nil {
		return err
	}
	if err := validateHealthGainFactor(p.HealthGainFactor); err != nil {
		return err
	}
	if err := validateEpochLength(p.EpochLength); err != nil {
		return err
	}
	if err := validateRemovalQueueThreshold(p.RemovalQueueThreshold); err != nil {
		return err
	}
	if err := validateMaxOpenPositions(p.MaxOpenPositions); err != nil {
		return err
	}
	if err := validatePoolOpenThreshold(p.PoolOpenThreshold); err != nil {
		return err
	}
	if err := validateForceCloseFundPercentage(p.ForceCloseFundPercentage); err != nil {
		return err
	}
	if err := validateForceCloseFundAddress(p.ForceCloseFundAddress); err != nil {
		return err
	}
	if err := validateIncrementalBorrowInterestPaymentFundPercentage(p.IncrementalBorrowInterestPaymentFundPercentage); err != nil {
		return err
	}
	if err := validateIncrementalBorrowInterestPaymentFundAddress(p.IncrementalBorrowInterestPaymentFundAddress); err != nil {
		return err
	}
	if err := validateSqModifier(p.SqModifier); err != nil {
		return err
	}
	if err := validateSafetyFactor(p.SafetyFactor); err != nil {
		return err
	}
	if err := validateIncrementalBorrowInterestPaymentEnabled(p.IncrementalBorrowInterestPaymentEnabled); err != nil {
		return err
	}
	if err := validateWhitelistingEnabled(p.WhitelistingEnabled); err != nil {
		return err
	}
	if err := validateInvariantCheckEpoch(p.InvariantCheckEpoch); err != nil {
		return err
	}
	if err := validateBrokerAddress(p.BrokerAddress); err != nil {
		return err
	}
	if err := validateTakeProfitBorrowInterestRateMin(p.TakeProfitBorrowInterestRateMin); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateLeverageMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("leverage max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("leverage max must be positive: %s", v)
	}
	if v.GT(sdk.NewDec(10)) {
		return fmt.Errorf("leverage max too large: %s", v)
	}

	return nil
}

func validateBorrowInterestRateMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest max must be positive: %s", v)
	}

	return nil
}

func validateBorrowInterestRateMin(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest min must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest min must be positive: %s", v)
	}

	return nil
}

func validateBorrowInterestRateIncrease(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest rate increase must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest rate increase must be positive: %s", v)
	}

	return nil
}

func validateBorrowInterestRateDecrease(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest rate decrease must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest rate decrease must be positive: %s", v)
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

func validateEpochLength(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateRemovalQueueThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("removal queue threashold must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("removal queue threashold must be positive: %s", v)
	}

	return nil
}

func validateMaxOpenPositions(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateForceCloseFundPercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("force close fund percentage must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("force close fund percentage must be positive: %s", v)
	}

	return nil
}

func validateForceCloseFundAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateIncrementalBorrowInterestPaymentFundPercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("incremental borrow interest payment fund percentage must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("incremental borrow interest payment fund percentage must be positive: %s", v)
	}

	return nil
}

func validateIncrementalBorrowInterestPaymentFundAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateSqModifier(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("sq modifier must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("sq modifier must be positive: %s", v)
	}

	return nil
}

func validateSafetyFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("safety factor must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("safety factor must be positive: %s", v)
	}

	return nil
}

func validateIncrementalBorrowInterestPaymentEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateWhitelistingEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validatePoolOpenThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("pool open threshold must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("pool open threshold must be positive: %s", v)
	}

	return nil
}

func validateInvariantCheckEpoch(i interface{}) error {
	epoch, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if epoch != epochtypes.DayEpochID && epoch != epochtypes.WeekEpochID && epoch != epochtypes.HourEpochID {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBrokerAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateTakeProfitBorrowInterestRateMin(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("take profit borrow interest rate min must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("take profit borrow interest rate min must be positive: %s", v)
	}

	return nil
}
