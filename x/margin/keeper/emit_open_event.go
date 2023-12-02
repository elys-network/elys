package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) EmitOpenEvent(ctx sdk.Context, mtp *types.MTP) {
	ctx.EventManager().EmitEvent(types.GenerateOpenEvent(mtp))
}
