package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (k msgServer) Dewhitelist(goCtx context.Context, msg *types.MsgDewhitelist) (*types.MsgDewhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	accAddress, err := sdk.AccAddressFromBech32(msg.WhitelistedAddress)
	if err != nil {
		return nil, err
	}
	k.Keeper.DewhitelistAddress(ctx, accAddress)

	return &types.MsgDewhitelistResponse{}, nil
}
