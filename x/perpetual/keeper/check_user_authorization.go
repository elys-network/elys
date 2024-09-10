package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckUserAuthorization(ctx sdk.Context, msg *types.MsgOpen) error {
	if !k.IsWhitelistingEnabled(ctx) {
		return fmt.Errorf("whitelisting not enabled")
	}
	if !k.CheckIfWhitelisted(ctx, sdk.MustAccAddressFromBech32(msg.Creator)) {
		return errorsmod.Wrap(types.ErrUnauthorised, "unauthorised")
	}
	return nil
}
