package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/parameter/types"
)

func NewParameterChangeProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ProposalUpdateMinCommission:
			return handleUpdateMinCommission(ctx, k, c)
		case *types.ProposalUpdateMaxVotingPower:
			return handleUpdateMaxVotingPower(ctx, k, c)
		case *types.ProposalUpdateMinSelfDelegation:
			return handleUpdateMinSelfDelegation(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized software upgrade proposal content type: %T", c)
		}
	}
}

// Update min commission
func handleUpdateMinCommission(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalUpdateMinCommission) error {
	minComission, err := sdk.NewDecFromStr(p.MinCommission)
	if err != nil {
		return err
	}

	params, bfound := k.GetAnteHandlerParam(ctx)
	if !bfound {
		return nil
	}
	params.MinCommissionRate = minComission
	k.SetAnteHandlerParam(ctx, params)

	return nil
}

// Update max voting power
func handleUpdateMaxVotingPower(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalUpdateMaxVotingPower) error {
	maxVotingPower, err := sdk.NewDecFromStr(p.MaxVotingPower)
	if err != nil {
		return err
	}

	params, bfound := k.GetAnteHandlerParam(ctx)
	if !bfound {
		return nil
	}
	params.MaxVotingPower = maxVotingPower
	k.SetAnteHandlerParam(ctx, params)

	return nil
}

// Update min self delegation
func handleUpdateMinSelfDelegation(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalUpdateMinSelfDelegation) error {
	minSelfDelegation, ok := sdk.NewIntFromString(p.MinSelfDelegation)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	params, bfound := k.GetAnteHandlerParam(ctx)
	if !bfound {
		return nil
	}
	params.MinSelfDelegation = minSelfDelegation
	k.SetAnteHandlerParam(ctx, params)

	return nil
}
