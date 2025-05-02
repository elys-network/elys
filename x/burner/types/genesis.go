package types

import "errors"

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		HistoryList: []History{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in history
	historyIndexMap := make(map[uint64]struct{})

	for _, elem := range gs.HistoryList {
		if _, ok := historyIndexMap[elem.Block]; ok {
			return errors.New("duplicated index for history")
		}
		historyIndexMap[elem.Block] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
