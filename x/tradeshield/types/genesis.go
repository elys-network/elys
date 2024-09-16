package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PendingSpotOrderList: []PendingSpotOrder{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in pendingSpotOrder
	pendingSpotOrderIdMap := make(map[uint64]bool)
	pendingSpotOrderCount := gs.GetPendingSpotOrderCount()
	for _, elem := range gs.PendingSpotOrderList {
		if _, ok := pendingSpotOrderIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for pendingSpotOrder")
		}
		if elem.Id >= pendingSpotOrderCount {
			return fmt.Errorf("pendingSpotOrder id should be lower or equal than the last id")
		}
		pendingSpotOrderIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
