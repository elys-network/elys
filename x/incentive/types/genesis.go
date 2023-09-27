package types

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

//nolint:interfacer
func NewGenesisState(
	params Params, fp FeePool,
) *GenesisState {
	return &GenesisState{
		Params:  params,
		FeePool: fp,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:  DefaultParams(),
		FeePool: InitialFeePool(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return gs.FeePool.ValidateGenesis()
}
