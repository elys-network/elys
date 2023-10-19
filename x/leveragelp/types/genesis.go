package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PoolList:         []Pool{},
		MtpList:          []MTP{},
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
	// Check for duplicated index in mtp
	mtpIndexMap := make(map[string]struct{})

	for _, elem := range gs.MtpList {
		key := GetMTPKey(elem.Address, elem.Id)
		index := string(key)
		if _, ok := mtpIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pool")
		}
		mtpIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in mtp
	whitelistMap := make(map[string]struct{})
	for _, elem := range gs.AddressWhitelist {
		index := elem
		if _, ok := mtpIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pool")
		}
		whitelistMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
