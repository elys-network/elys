package types

import (
	"errors"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PoolList:           []Pool{},
		DenomLiquidityList: []DenomLiquidity{},
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
		index := string(PoolKey(elem.PoolId))
		if _, ok := poolIndexMap[index]; ok {
			return errors.New("duplicated index for pool")
		}
		poolIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in denomLiquidity
	denomLiquidityIndexMap := make(map[string]struct{})

	for _, elem := range gs.DenomLiquidityList {
		index := string(DenomLiquidityKey(elem.Denom))
		if _, ok := denomLiquidityIndexMap[index]; ok {
			return errors.New("duplicated index for denomLiquidity")
		}
		denomLiquidityIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
