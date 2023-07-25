package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	anteParam := AnteHandlerParam{
		MinCommissionRate: sdk.NewDecWithPrec(5, 2),
		MaxVotingPower:    sdk.NewDecWithPrec(66, 1),
		MinSelfDelegation: sdk.NewInt(1)}

	return &GenesisState{
		AnteHandlerParam: anteParam,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
