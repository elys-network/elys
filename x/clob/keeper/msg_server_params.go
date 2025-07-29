package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check authority
	if msg.Authority != k.authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to set params")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventUpdateParams,
			sdk.NewAttribute(types.AttributeAuthority, msg.Authority),
		),
	)

	return &types.MsgUpdateParamsResponse{}, nil
}
