package types

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeUpdatePoolMultipliers string = "ProposalUpdatePoolMultipliers"
)

func init() {
	gov.RegisterProposalType(ProposalTypeUpdatePoolMultipliers)
}

// NewProposalAddAssetInfo creates a new ProposalAddAssetInfo instance
func NewProposalUpdatePoolMultipliers(
	title, description string,
	poolMultipliers []PoolMultipliers,
) gov.Content {
	return &ProposalUpdatePoolMultipliers{
		Title:           title,
		Description:     description,
		PoolMultipliers: poolMultipliers,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalUpdatePoolMultipliers{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalUpdatePoolMultipliers) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalUpdatePoolMultipliers) ProposalType() string {
	return ProposalTypeUpdatePoolMultipliers
}

// ValidateBasic validates the proposal
func (sup *ProposalUpdatePoolMultipliers) ValidateBasic() error {
	return gov.ValidateAbstract(sup)
}
