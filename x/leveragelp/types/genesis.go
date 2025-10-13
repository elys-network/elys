package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PoolList:         []Pool{},
		PositionList:     []Position{},
		AddressWhitelist: []string{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in pool
	poolIndexMap := make(map[string]struct{})

	for _, elem := range gs.PoolList {
		index := string(PoolKey(elem.AmmPoolId))
		if _, ok := poolIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pool")
		}
		poolIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in position
	positionIndexMap := make(map[string]struct{})

	for _, elem := range gs.PositionList {
		key := GetPositionKey(elem.AmmPoolId, sdk.MustAccAddressFromBech32(elem.Address), elem.Id)
		index := string(key)
		if _, ok := positionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pool")
		}
		positionIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in position
	whitelistMap := make(map[string]struct{})
	for _, elem := range gs.AddressWhitelist {
		index := elem
		if _, ok := whitelistMap[index]; ok {
			return fmt.Errorf("duplicated index for pool")
		}
		whitelistMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
