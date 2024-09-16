package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PendingSpotOrderList: []SpotOrder{},
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
		if _, ok := pendingSpotOrderIdMap[elem.OrderId]; ok {
			return fmt.Errorf("duplicated id for pendingSpotOrder")
		}
		if elem.OrderId >= pendingSpotOrderCount {
			return fmt.Errorf("pendingSpotOrder id should be lower or equal than the last id")
		}
		pendingSpotOrderIdMap[elem.OrderId] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
