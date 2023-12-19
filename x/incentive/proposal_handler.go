package incentive

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/x/incentive/keeper"
	"github.com/elys-network/elys/x/incentive/types"
)

func NewPoolInfoProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ProposalUpdatePoolMultipliers:
			return handleUpdatePoolMultipliers(ctx, k, c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized software upgrade proposal content type: %T", c)
		}
	}
}

// handleUpdatePoolMultipliers updates pool multiplier parameter
func handleUpdatePoolMultipliers(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalUpdatePoolMultipliers) error {
	k.UpdatePoolMultipliers(ctx, p.PoolMultipliers)
	return nil
}
