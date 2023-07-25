package margin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
)

func NewMarginProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ProposalUpdateParams:
			return handleUpdateParamsProposal(ctx, k, c)

		case *types.ProposalUpdatePools:
			return handleUpdatePoolsProposal(ctx, k, c)

		case *types.ProposalWhitelist:
			return handleWhitelistProposal(ctx, k, c)

		case *types.ProposalDewhitelist:
			return handleDewhitelistProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized software upgrade proposal content type: %T", c)
		}
	}
}

// Handle updating margin param
func handleUpdateParamsProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalUpdateParams) error {
	k.SetParams(ctx, *p.Params)
	return nil
}

// Handle updating margin pools
func handleUpdatePoolsProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalUpdatePools) error {
	return nil
}

// Handle whitelisting
func handleWhitelistProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalWhitelist) error {
	return nil
}

// Handle de whitelisting
func handleDewhitelistProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalDewhitelist) error {
	return nil
}
