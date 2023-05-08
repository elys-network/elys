package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyBandEpoch         = []byte("BandEpoch")
	KeyClientID          = []byte("ClientID")
	KeyOracleScriptID    = []byte("OracleScriptID")
	KeyMultiplier        = []byte("Multiplier")
	KeyBandChannelSource = []byte("BandChannelSource")
	KeyAskCount          = []byte("AskCount")
	KeyMinCount          = []byte("MinCount")
	KeyFeeLimit          = []byte("FeeLimit")
	KeyPrepareGas        = []byte("PrepareGas")
	KeyExecuteGas        = []byte("ExecuteGas")
	KeyModuleAdmin       = []byte("ModuleAdmin")
	KeyPriceExpiryTime   = []byte("PriceExpiryTime")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	bandEpoch string,
	oracleScriptID OracleScriptID,
	multiplier uint64,
	bandChannelSrc string,
	askCount uint64,
	minCount uint64,
	feeLimit sdk.Coins,
	prepareGas uint64,
	executeGas uint64,
	priceExpiryTime uint64,
) Params {
	return Params{
		BandEpoch:         bandEpoch,
		ClientID:          BandPriceClientIDKey,
		OracleScriptID:    uint64(oracleScriptID),
		Multiplier:        multiplier,
		BandChannelSource: bandChannelSrc,
		AskCount:          askCount,
		MinCount:          minCount,
		FeeLimit:          feeLimit,
		PrepareGas:        prepareGas,
		ExecuteGas:        executeGas,
		PriceExpiryTime:   priceExpiryTime,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"band_epoch",
		37,
		18,          // decimal 18
		"channel-1", // used on dockernet
		4,
		3,
		sdk.NewCoins(sdk.NewInt64Coin("uband", 30)),
		600000,
		600000,
		86400, // 1 day old data
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBandEpoch, &p.BandEpoch, validateBandEpoch),
		paramtypes.NewParamSetPair(KeyClientID, &p.ClientID, validateClientID),
		paramtypes.NewParamSetPair(KeyOracleScriptID, &p.OracleScriptID, validateOracleScriptID),
		paramtypes.NewParamSetPair(KeyMultiplier, &p.Multiplier, validateMultiplier),
		paramtypes.NewParamSetPair(KeyBandChannelSource, &p.BandChannelSource, validateChannel),
		paramtypes.NewParamSetPair(KeyAskCount, &p.AskCount, validateCount),
		paramtypes.NewParamSetPair(KeyMinCount, &p.MinCount, validateCount),
		paramtypes.NewParamSetPair(KeyFeeLimit, &p.FeeLimit, validateFeeLimit),
		paramtypes.NewParamSetPair(KeyPrepareGas, &p.PrepareGas, validateGas),
		paramtypes.NewParamSetPair(KeyExecuteGas, &p.ExecuteGas, validateGas),
		paramtypes.NewParamSetPair(KeyPriceExpiryTime, &p.PriceExpiryTime, validatePriceExpiryTime),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {

	if err := validateBandEpoch(p.BandEpoch); err != nil {
		return err
	}
	if err := validateClientID(p.ClientID); err != nil {
		return err
	}
	if err := validateOracleScriptID(p.OracleScriptID); err != nil {
		return err
	}
	if err := validateMultiplier(p.Multiplier); err != nil {
		return err
	}
	if err := validateChannel(p.BandChannelSource); err != nil {
		return err
	}
	if err := validateCount(p.AskCount); err != nil {
		return err
	}
	if err := validateCount(p.MinCount); err != nil {
		return err
	}
	if err := validateFeeLimit(p.FeeLimit); err != nil {
		return err
	}
	if err := validateGas(p.PrepareGas); err != nil {
		return err
	}
	if err := validateGas(p.ExecuteGas); err != nil {
		return err
	}
	if err := validatePriceExpiryTime(p.PriceExpiryTime); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateBandEpoch(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid type for band epoch: %T", i)
	}

	if v == "" {
		return fmt.Errorf("band epoch must not be empty: %s", v)
	}
	return nil
}

func validateClientID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid type for client id: %T", i)
	}

	if v == "" {
		return fmt.Errorf("clientID must not be empty: %s", v)
	}
	return nil
}

func validateOracleScriptID(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid type for oracle script id: %T", i)
	}

	return nil
}

func validateMultiplier(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid type for multiplier: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("multiplier should be positive: %d", v)
	}

	return nil
}

func validateChannel(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid type for channel: %T", i)
	}

	if v == "" {
		return fmt.Errorf("channel should not be empty: %s", v)
	}

	return nil
}

func validateCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid type for count: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("count should not be zero: %d", v)
	}

	return nil
}

func validateFeeLimit(i interface{}) error {
	_, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid type for fee limit: %T", i)
	}

	return nil
}

func validateGas(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid type for gas: %T", i)
	}

	return nil
}

func validatePriceExpiryTime(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid type for module admin: %T", i)
	}

	return nil
}
