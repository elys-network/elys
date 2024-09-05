package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fmt "fmt"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DebtList:     []Debt{},
		InterestList: []InterestBlock{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in debt
	debtIndexMap := make(map[string]struct{})
	for _, elem := range gs.DebtList {
		index := elem.Address
		if _, ok := debtIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for debt")
		}
		debtIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in interest
	interestIndexMap := make(map[string]struct{})
	for _, elem := range gs.InterestList {
		index := string(sdk.Uint64ToBigEndian(uint64(elem.BlockTime)))
		if _, ok := interestIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for interest")
		}
		interestIndexMap[index] = struct{}{}
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
