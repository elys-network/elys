package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) WithdrawTokens(goCtx context.Context, msg *types.MsgWithdrawTokens) (*types.MsgWithdrawTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Withdraw tokens
	err := k.ProcessWithdrawTokens(ctx, msg.Creator, msg.Denom, msg.Amount)

	return &types.MsgWithdrawTokensResponse{}, err
}
