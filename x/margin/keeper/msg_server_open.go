package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.Open(ctx, msg)
}
