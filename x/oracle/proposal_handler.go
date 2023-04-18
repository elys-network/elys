package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

func NewAssetInfoProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ProposalAddAssetInfo:
			return handleSoftwareUpgradeProposal(ctx, k, c)

		case *types.ProposalRemoveAssetInfo:
			return handleCancelSoftwareUpgradeProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized software upgrade proposal content type: %T", c)
		}
	}
}

func handleSoftwareUpgradeProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalAddAssetInfo) error {
	k.SetAssetInfo(ctx, types.AssetInfo{
		Denom:      p.Denom,
		Display:    p.Display,
		BandTicker: p.BandTicker,
		ElysTicker: p.ElysTicker,
	})
	return nil
}

func handleCancelSoftwareUpgradeProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ProposalRemoveAssetInfo) error {
	k.RemoveAssetInfo(ctx, p.Denom)
	return nil
}
