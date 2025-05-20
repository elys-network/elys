package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/commitment/types"
)

func (k msgServer) UpdateAirdropParams(goCtx context.Context, msg *types.MsgUpdateAirdropParams) (*types.MsgUpdateAirdropParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	params := k.GetParams(ctx)
	params.EnableClaim = msg.EnableClaim
	params.StartAirdropClaimHeight = msg.StartAirdropClaimHeight
	params.EndAirdropClaimHeight = msg.EndAirdropClaimHeight
	params.StartKolClaimHeight = msg.StartKolClaimHeight
	params.EndKolClaimHeight = msg.EndKolClaimHeight
	k.SetParams(ctx, params)

	return &types.MsgUpdateAirdropParamsResponse{}, nil
}
