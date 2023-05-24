package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeUpdateMinCommissionRate string = "UpdateMinCommissionRate"
	ProposalTypeUpdateMaxVotingPower    string = "UpdateMaxVotingPower"
	ProposalTypeUpdateMinSelfDelegation string = "UpdateMinSelfDelegation"
)

func init() {
	gov.RegisterProposalType(ProposalTypeUpdateMinCommissionRate)
	gov.RegisterProposalType(ProposalTypeUpdateMaxVotingPower)
	gov.RegisterProposalType(ProposalTypeUpdateMinSelfDelegation)
}

// NewProposalUpdateMinCommission creates a new ProposalUpdateMinCommission instance
func NewProposalUpdateMinCommission(
	title, description string,
	minCommission string,
) gov.Content {
	return &ProposalUpdateMinCommission{
		Title:         title,
		Description:   description,
		MinCommission: minCommission,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalUpdateMinCommission{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalUpdateMinCommission) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalUpdateMinCommission) ProposalType() string {
	return ProposalTypeUpdateMinCommissionRate
}

// ValidateBasic validates the proposal
func (sup *ProposalUpdateMinCommission) ValidateBasic() error {
	_, err := sdk.NewDecFromStr(sup.MinCommission)
	if err != nil {
		return err
	}
	// sup.MinCommission
	return gov.ValidateAbstract(sup)
}

// NewProposalUpdateMaxVotingPower creates a new ProposalUpdateMaxVotingPower instance
func NewProposalUpdateMaxVotingPower(
	title, description string,
	maxVotingPower string,
) gov.Content {
	return &ProposalUpdateMaxVotingPower{
		Title:          title,
		Description:    description,
		MaxVotingPower: maxVotingPower,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalUpdateMaxVotingPower{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalUpdateMaxVotingPower) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalUpdateMaxVotingPower) ProposalType() string {
	return ProposalTypeUpdateMaxVotingPower
}

// ValidateBasic validates the proposal
func (sup *ProposalUpdateMaxVotingPower) ValidateBasic() error {
	_, err := sdk.NewDecFromStr(sup.MaxVotingPower)
	if err != nil {
		return err
	}
	return gov.ValidateAbstract(sup)
}

// NewProposalUpdateMinSelfDelegation creates a new NewProposalUpdateMinSelfDelegation instance
func NewProposalUpdateMinSelfDelegation(
	title, description string,
	minSelfDelegation string,
) gov.Content {
	return &ProposalUpdateMinSelfDelegation{
		Title:             title,
		Description:       description,
		MinSelfDelegation: minSelfDelegation,
	}
}

// Implements Proposal Interface
var _ gov.Content = &ProposalUpdateMinSelfDelegation{}

// ProposalRoute gets the proposal's router key
func (sup *ProposalUpdateMinSelfDelegation) ProposalRoute() string { return RouterKey }

// ProposalType is "SoftwareUpgrade"
func (sup *ProposalUpdateMinSelfDelegation) ProposalType() string {
	return ProposalTypeUpdateMinSelfDelegation
}

// ValidateBasic validates the proposal
func (sup *ProposalUpdateMinSelfDelegation) ValidateBasic() error {
	_, ok := sdk.NewIntFromString(sup.MinSelfDelegation)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	return gov.ValidateAbstract(sup)
}
