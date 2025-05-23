package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
)

func (k msgServer) UpdateEnableVestNow(goCtx context.Context, msg *types.MsgUpdateEnableVestNow) (*types.MsgUpdateEnableVestNowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	params := k.GetParams(ctx)
	params.EnableVestNow = msg.EnableVestNow

	k.SetParams(ctx, params)
	return &types.MsgUpdateEnableVestNowResponse{}, nil
}
