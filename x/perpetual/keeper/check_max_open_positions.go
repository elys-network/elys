package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	if k.GetOpenMTPCount(ctx) >= k.GetMaxOpenPositions(ctx) {
		return errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}
