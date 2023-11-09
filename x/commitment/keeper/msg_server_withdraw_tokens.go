package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// WithdrawTokens withdraw first from unclaimed and if it requires more, withdraw from committed store
func (k msgServer) WithdrawTokens(goCtx context.Context, msg *types.MsgWithdrawTokens) (*types.MsgWithdrawTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Withdraw tokens
	err := k.ProcessWithdrawTokens(ctx, msg.Creator, msg.Denom, msg.Amount)

	return &types.MsgWithdrawTokensResponse{}, err
}
