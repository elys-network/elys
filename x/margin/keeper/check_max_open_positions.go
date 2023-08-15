package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	if k.PositionChecker.GetOpenMTPCount(ctx) >= (uint64)(k.PositionChecker.GetMaxOpenPositions(ctx)) {
		return sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}
