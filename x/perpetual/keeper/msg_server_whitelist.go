package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) Whitelist(goCtx context.Context, msg *types.MsgWhitelist) (*types.MsgWhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.WhitelistAddress(ctx, msg.WhitelistedAddress)

	return &types.MsgWhitelistResponse{}, nil
}
