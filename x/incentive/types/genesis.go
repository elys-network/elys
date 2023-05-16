package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

//nolint:interfacer
func NewGenesisState(
	params Params, fp FeePool, dwis []DelegatorWithdrawInfo, pp sdk.ConsAddress, r []ValidatorOutstandingRewardsRecord,
	acc []ValidatorAccumulatedCommissionRecord, historical []ValidatorHistoricalRewardsRecord,
	cur []ValidatorCurrentRewardsRecord, dels []DelegatorStartingInfoRecord, slashes []ValidatorSlashEventRecord,
) *GenesisState {
	return &GenesisState{
		Params:                          params,
		FeePool:                         fp,
		DelegatorWithdrawInfos:          dwis,
		PreviousProposer:                pp.String(),
		OutstandingRewards:              r,
		ValidatorAccumulatedCommissions: acc,
		ValidatorHistoricalRewards:      historical,
		ValidatorCurrentRewards:         cur,
		DelegatorStartingInfos:          dels,
		ValidatorSlashEvents:            slashes,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:                          DefaultParams(),
		FeePool:                         InitialFeePool(),
		DelegatorWithdrawInfos:          []DelegatorWithdrawInfo{},
		PreviousProposer:                "",
		OutstandingRewards:              []ValidatorOutstandingRewardsRecord{},
		ValidatorAccumulatedCommissions: []ValidatorAccumulatedCommissionRecord{},
		ValidatorHistoricalRewards:      []ValidatorHistoricalRewardsRecord{},
		ValidatorCurrentRewards:         []ValidatorCurrentRewardsRecord{},
		DelegatorStartingInfos:          []DelegatorStartingInfoRecord{},
		ValidatorSlashEvents:            []ValidatorSlashEventRecord{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	return gs.FeePool.ValidateGenesis()
}
