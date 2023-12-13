package types

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AccountedPoolList: []AccountedPool{},
		// this line is used by starport scaffolding # genesis/types/default
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in accountedPool
	accountedPoolIndexMap := make(map[string]struct{})

	for _, elem := range gs.AccountedPoolList {
		index := string(AccountedPoolKey(elem.PoolId))
		if _, ok := accountedPoolIndexMap[index]; ok {
			return ErrDuplicatedAccountedPoolId
		}
		accountedPoolIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return nil
}
