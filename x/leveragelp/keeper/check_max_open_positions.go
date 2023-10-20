package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	if k.GetOpenMTPCount(ctx) >= (uint64)(k.GetMaxOpenPositions(ctx)) {
		return sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}
