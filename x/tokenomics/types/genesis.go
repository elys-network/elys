package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AirdropList:            []Airdrop{},
		GenesisInflation:       nil,
		TimeBasedInflationList: []TimeBasedInflation{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in airdrop
	airdropIndexMap := make(map[string]struct{})

	for _, elem := range gs.AirdropList {
		index := string(AirdropKey(elem.Intent))
		if _, ok := airdropIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for airdrop")
		}
		airdropIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in timeBasedInflation
	timeBasedInflationIndexMap := make(map[string]struct{})

	for _, elem := range gs.TimeBasedInflationList {
		index := string(TimeBasedInflationKey(elem.StartBlockHeight, elem.EndBlockHeight))
		if _, ok := timeBasedInflationIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for timeBasedInflation")
		}
		timeBasedInflationIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
