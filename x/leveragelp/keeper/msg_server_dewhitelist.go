package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
)

func (k msgServer) Dewhitelist(goCtx context.Context, msg *types.MsgDewhitelist) (*types.MsgDewhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	whitelistAddress := sdk.MustAccAddressFromBech32(msg.WhitelistedAddress)
	k.Keeper.DewhitelistAddress(ctx, whitelistAddress)

	return &types.MsgDewhitelistResponse{}, nil
}
