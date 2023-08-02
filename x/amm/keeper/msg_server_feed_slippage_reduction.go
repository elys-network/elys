package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) FeedSlippageReduction(goCtx context.Context, msg *types.MsgFeedSlippageReduction) (*types.MsgFeedSlippageReductionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgFeedSlippageReductionResponse{}, nil
}
