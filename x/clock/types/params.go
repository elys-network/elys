package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const MinimumContractGasLimit = uint64(100_000)

var (
	KeyContractAddressesIdentifier = []byte("ContractAddresses")
	KeyContractGasLimitIdentifier  = []byte("ContractGasLimit")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		ContractAddresses: []string(nil),
		ContractGasLimit:  1_000_000_000, // 1 billion
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyContractAddressesIdentifier, &p.ContractAddresses, validateContractAddressesIdentifier),
		paramtypes.NewParamSetPair(KeyContractGasLimitIdentifier, &p.ContractGasLimit, validateContractGasLimitIdentifier),
	}
}

// NewParams creates a new Params object
func NewParams(
	contracts []string,
	contractGasLimit uint64,
) Params {
	return Params{
		ContractAddresses: contracts,
		ContractGasLimit:  contractGasLimit,
	}
}

// Validate performs basic validation.
func (p Params) Validate() error {

	if p.ContractGasLimit < MinimumContractGasLimit {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid contract gas limit: %d. Must be above %d", p.ContractGasLimit, MinimumContractGasLimit,
		)
	}

	contractAddressesMap := make(map[string]bool)

	for _, addr := range p.ContractAddresses {
		// Valid address check
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid contract address: %s", err.Error(),
			)
		}

		// duplicate address check
		if contractAddressesMap[addr] {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "duplicate contract address: %s", addr)
		} else {
			contractAddressesMap[addr] = true
		}

	}

	return nil
}

// validateContractAddressesIdentifier validates the ContractAddresses param
func validateContractAddressesIdentifier(v interface{}) error {
	contractAddressesIdentifier, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = contractAddressesIdentifier

	return nil
}

// validateContractAddressesIdentifier validates the ContractAddresses param
func validateContractGasLimitIdentifier(v interface{}) error {
	contractGasLimitIdentifier, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = contractGasLimitIdentifier

	return nil
}
