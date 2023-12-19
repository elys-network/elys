package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) Dewhitelist(goCtx context.Context, msg *types.MsgDewhitelist) (*types.MsgDewhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.DewhitelistAddress(ctx, msg.WhitelistedAddress)

	return &types.MsgDewhitelistResponse{}, nil
}
