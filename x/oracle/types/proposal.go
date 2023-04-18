package types

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddAssetInfo    string = "AddAssetInfo"
	ProposalTypeRemoveAssetInfo string = "RemoveAssetInfo"
)

// NewProposalAddAssetInfo creates a new ProposalAddAssetInfo instance
func NewProposalAddAssetInfo(
	title, description string,
	denom string,
	display string,
	bandTicker string,
	elysTicker string,
) gov.Content {
	return &ProposalAddAssetInfo{
		Title:       title,
		Description: description,
		Denom:       denom,
		Display:     display,
		BandTicker:  bandTicker,
		ElysTicker:  elysTicker,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalAddAssetInfo{}

func init() {
	gov.RegisterProposalType(ProposalTypeAddAssetInfo)
	gov.RegisterProposalType(ProposalTypeRemoveAssetInfo)
}

// ProposalRoute gets the proposal's router key
func (sup *ProposalAddAssetInfo) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalAddAssetInfo) ProposalType() string { return ProposalTypeAddAssetInfo }

// ValidateBasic validates the proposal
func (sup *ProposalAddAssetInfo) ValidateBasic() error {
	return gov.ValidateAbstract(sup)
}

func NewProposalRemoveAssetInfo(title, description, denom string) gov.Content {
	return &ProposalRemoveAssetInfo{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalRemoveAssetInfo{}

func (csup *ProposalRemoveAssetInfo) ProposalRoute() string { return RouterKey }
func (csup *ProposalRemoveAssetInfo) ProposalType() string {
	return ProposalTypeRemoveAssetInfo
}

func (csup *ProposalRemoveAssetInfo) ValidateBasic() error {
	return gov.ValidateAbstract(csup)
}
