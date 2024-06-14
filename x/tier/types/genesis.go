package types

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortfolioList: []Portfolio{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
// No data at genesis
func (gs GenesisState) Validate() error {
	// Check for duplicated index in portfolio
	// portfolioIndexMap := make(map[string]struct{})

	// for _, elem := range gs.PortfolioList {
	// 	index := string(PortfolioKey(elem.Index))
	// 	if _, ok := portfolioIndexMap[index]; ok {
	// 		return fmt.Errorf("duplicated index for portfolio")
	// 	}
	// 	portfolioIndexMap[index] = struct{}{}
	// }
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
