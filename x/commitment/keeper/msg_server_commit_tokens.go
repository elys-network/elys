package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) CommitTokens(goCtx context.Context, msg *types.MsgCommitTokens) (*types.MsgCommitTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCommitTokensResponse{}, nil
}
