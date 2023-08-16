package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) Whitelist(goCtx context.Context, msg *types.MsgWhitelist) (*types.MsgWhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.WhitelistAddress(ctx, msg.WhitelistedAddress)

	return &types.MsgWhitelistResponse{}, nil
}
