package types

import (
	"fmt"

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

	anteParamList := make([]AnteHandlerParam, 0)
	anteParamList = append(anteParamList, anteParam)
	return &GenesisState{
		AnteHandlerParamList: anteParamList,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in anteHandlerParam
	anteHandlerParamIndexMap := make(map[string]struct{})

	for range gs.AnteHandlerParamList {
		index := string(AnteHandlerParamKey(AnteStoreKey))
		if _, ok := anteHandlerParamIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for anteHandlerParam")
		}
		anteHandlerParamIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
