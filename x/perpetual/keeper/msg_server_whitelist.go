package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

func (k msgServer) Whitelist(goCtx context.Context, msg *types.MsgWhitelist) (*types.MsgWhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	accAddress, err := sdk.AccAddressFromBech32(msg.WhitelistedAddress)
	if err != nil {
		return nil, err
	}

	k.Keeper.WhitelistAddress(ctx, accAddress)

	return &types.MsgWhitelistResponse{}, nil
}
