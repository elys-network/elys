package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) DepositTokens(goCtx context.Context, msg *types.MsgDepositTokens) (*types.MsgDepositTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDepositTokensResponse{}, nil
}
