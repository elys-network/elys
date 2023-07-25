package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeUpdateParams string = "UpdateParams"
	ProposalTypeUpdatePools  string = "UpdatePools"
	ProposalTypeWhitelist    string = "Whitelist"
	ProposalTypeDewhitelist  string = "Dewhitelist"
)

func init() {
	gov.RegisterProposalType(ProposalTypeUpdateParams)
	gov.RegisterProposalType(ProposalTypeUpdatePools)
	gov.RegisterProposalType(ProposalTypeWhitelist)
	gov.RegisterProposalType(ProposalTypeDewhitelist)
}

// NewProposalUpdateParams
func NewProposalUpdateParams(
	title, description string,
	LeverageMax sdk.Dec,
	InterestRateMax sdk.Dec,
	InterestRateMin sdk.Dec,
	InterestRateIncrease sdk.Dec,
	InterestRateDecrease sdk.Dec,
	HealthGainFactor sdk.Dec,
	EpochLength uint64,
	RemovalQueueThreshold sdk.Dec,
	MaxOpenPositions uint64,
	PoolOpenThreshold sdk.Dec,
	ForceCloseFundPercentage sdk.Dec,
	ForceCloseFundAddress string,
	IncrementalInterestPaymentFundPercentage sdk.Dec,
	IncrementalInterestPaymentFundAddress string,
	SqModifier sdk.Dec,
	SafetyFactor sdk.Dec,
	IncrementalInterestPaymentEnabled bool,
	WhitelistingEnabled bool,
) gov.Content {
	params := Params{
		LeverageMax:                              LeverageMax,
		InterestRateMax:                          InterestRateMax,
		InterestRateMin:                          InterestRateMin,
		InterestRateIncrease:                     InterestRateIncrease,
		InterestRateDecrease:                     InterestRateDecrease,
		HealthGainFactor:                         HealthGainFactor,
		EpochLength:                              EpochLength,
		RemovalQueueThreshold:                    RemovalQueueThreshold,
		MaxOpenPositions:                         MaxOpenPositions,
		PoolOpenThreshold:                        PoolOpenThreshold,
		ForceCloseFundPercentage:                 ForceCloseFundPercentage,
		ForceCloseFundAddress:                    ForceCloseFundAddress,
		IncrementalInterestPaymentFundPercentage: IncrementalInterestPaymentFundPercentage,
		IncrementalInterestPaymentFundAddress:    IncrementalInterestPaymentFundAddress,
		SqModifier:                               SqModifier,
		SafetyFactor:                             SafetyFactor,
		IncrementalInterestPaymentEnabled:        IncrementalInterestPaymentEnabled,
		WhitelistingEnabled:                      WhitelistingEnabled,
	}

	return &ProposalUpdateParams{
		Title:       title,
		Description: description,
		Params:      &params,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalUpdateParams{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalUpdateParams) ProposalType() string { return ProposalTypeUpdateParams }

// ValidateBasic validates the proposal
func (sup *ProposalUpdateParams) ValidateBasic() error {
	return gov.ValidateAbstract(sup)
}

// NewProposalUpdatePools
func NewProposalUpdatePools(
	title, description string,
	pools []string,
	close_pools []string,
) gov.Content {
	return &ProposalUpdatePools{
		Title:       title,
		Description: description,
		Pools:       pools,
		ClosedPools: close_pools,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalUpdatePools{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalUpdatePools) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalUpdatePools) ProposalType() string { return ProposalTypeUpdatePools }

// ValidateBasic validates the proposal
func (sup *ProposalUpdatePools) ValidateBasic() error {
	return gov.ValidateAbstract(sup)
}

// NewProposalWhitelist
func NewProposalWhitelist(
	title, description string,
	WhitelistedAddress string,
) gov.Content {
	return &ProposalWhitelist{
		Title:              title,
		Description:        description,
		WhitelistedAddress: WhitelistedAddress,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalWhitelist{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalWhitelist) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalWhitelist) ProposalType() string { return ProposalTypeWhitelist }

// ValidateBasic validates the proposal
func (sup *ProposalWhitelist) ValidateBasic() error {
	return gov.ValidateAbstract(sup)
}

// NewProposalDewhitelist
func NewProposalDewhitelist(
	title, description string,
	WhitelistedAddress string,
) gov.Content {
	return &ProposalDewhitelist{
		Title:              title,
		Description:        description,
		WhitelistedAddress: WhitelistedAddress,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalDewhitelist{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalDewhitelist) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalDewhitelist) ProposalType() string { return ProposalTypeDewhitelist }

// ValidateBasic validates the proposal
func (sup *ProposalDewhitelist) ValidateBasic() error {
	return gov.ValidateAbstract(sup)
}
