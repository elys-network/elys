package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) EmitCloseEvent(ctx sdk.Context, mtp *types.MTP, repayAmount math.Int) {
	ctx.EventManager().EmitEvent(types.GenerateCloseEvent(mtp, repayAmount))
}
