package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	// If set to -1, no limit on how many positions can be open
	if k.GetMaxOpenPositions(ctx) < 0 {
		return nil
	}
	if k.GetOpenMTPCount(ctx) >= uint64(k.GetMaxOpenPositions(ctx)) {
		return errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}
