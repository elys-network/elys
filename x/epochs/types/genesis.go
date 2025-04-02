package types

import (
	"fmt"
	"time"
)

// NewGenesisState creates a new genesis state instance
func NewGenesisState(epochs []EpochInfo) *GenesisState {
	return &GenesisState{Epochs: epochs}
}

// DefaultGenesisState returns the default epochs genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState([]EpochInfo{
		{
			Identifier:              BandEpochID,
			Duration:                time.Second * 15,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			EpochCountingStarted:    false,
		},
		{
			Identifier:              TenDaysEpochID,
			Duration:                time.Second * 864000,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			EpochCountingStarted:    false,
		},
		{
			Identifier:              FiveMinutesEpochID,
			Duration:                time.Second * 300,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			EpochCountingStarted:    false,
		},
		EightHourEpochInfo,
	})
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	epochIdentifiers := make(map[string]bool)

	for _, epoch := range gs.Epochs {
		if epochIdentifiers[epoch.Identifier] {
			return fmt.Errorf("duplicated epoch entry %s", epoch.Identifier)
		}
		if err := epoch.Validate(); err != nil {
			return err
		}
		epochIdentifiers[epoch.Identifier] = true
	}

	return nil
}
