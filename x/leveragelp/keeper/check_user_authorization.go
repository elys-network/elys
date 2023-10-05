package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CheckUserAuthorization(ctx sdk.Context, msg *types.MsgOpen) error {
	if k.AuthorizationChecker.IsWhitelistingEnabled(ctx) && !k.AuthorizationChecker.CheckIfWhitelisted(ctx, msg.Creator) {
		return sdkerrors.Wrap(types.ErrUnauthorised, "unauthorised")
	}
	return nil
}
