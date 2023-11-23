package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) SwapByDenom(goCtx context.Context, msg *types.MsgSwapByDenom) (*types.MsgSwapByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSwapByDenomResponse{}, nil
}
