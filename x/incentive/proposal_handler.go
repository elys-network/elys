package incentive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/app/params"
	"github.com/elys-network/elys/x/incentive/types"
)

func NewIncentiveProposalHandler() govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ProposalUpdateMinCommission:
			return handleUpdateMinCommission(ctx, c)
		case *types.ProposalUpdateMaxVotingPower:
			return handleUpdateMaxVotingPower(ctx, c)
		case *types.ProposalUpdateMinSelfDelegation:
			return handleUpdateMinSelfDelegation(ctx, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized software upgrade proposal content type: %T", c)
		}
	}
}

// Update min commission
func handleUpdateMinCommission(ctx sdk.Context, p *types.ProposalUpdateMinCommission) error {
	minComission, err := sdk.NewDecFromStr(p.MinCommission)
	if err != nil {
		return err
	}

	params.MinCommissionRate = minComission
	return nil
}

// Update max voting power
func handleUpdateMaxVotingPower(ctx sdk.Context, p *types.ProposalUpdateMaxVotingPower) error {
	maxVotingPower, err := sdk.NewDecFromStr(p.MaxVotingPower)
	if err != nil {
		return err
	}

	params.MaxVotingPower = maxVotingPower
	return nil
}

// Update min self delegation
func handleUpdateMinSelfDelegation(ctx sdk.Context, p *types.ProposalUpdateMinSelfDelegation) error {
	minSelfDelegation, ok := sdk.NewIntFromString(p.MinSelfDelegation)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}
	params.MinSelfDelegation = minSelfDelegation
	return nil
}
