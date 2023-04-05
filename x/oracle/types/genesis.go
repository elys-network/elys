package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId:       PortID,
		Params:       DefaultParams(),
		AssetInfos:   []AssetInfo{},
		Prices:       []Price{},
		PriceFeeders: []PriceFeeder{},
		// this line is used by starport scaffolding # genesis/types/default
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	// Check for duplicated index in assetInfo
	assetInfoIndexMap := make(map[string]struct{})

	for _, elem := range gs.AssetInfos {
		index := string(AssetInfoKey(elem.Denom))
		if _, ok := assetInfoIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for assetInfo")
		}
		assetInfoIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in price
	priceIndexMap := make(map[string]struct{})

	for _, elem := range gs.Prices {
		index := string(PriceKey(elem.Asset, elem.Source, elem.Timestamp))
		if _, ok := priceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for price")
		}
		priceIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in priceFeeder
	priceFeederIndexMap := make(map[string]struct{})

	for _, elem := range gs.PriceFeeders {
		index := string(PriceFeederKey(elem.Feeder))
		if _, ok := priceFeederIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for priceFeeder")
		}
		priceFeederIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
