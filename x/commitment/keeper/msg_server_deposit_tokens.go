package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) DepositTokens(goCtx context.Context, msg *types.MsgDepositTokens) (*types.MsgDepositTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositCoins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))
	// send the deposited coins to the module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.Creator), types.ModuleName, depositCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to send deposit tokens")
	}
	// burn the deposited coins
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, depositCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to burn deposit tokens")
	}

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		commitments = types.Commitments{
			Creator:           msg.Creator,
			CommittedTokens:   []*types.CommittedTokens{},
			UncommittedTokens: []*types.UncommittedTokens{},
		}
	}
	// Get the uncommitted tokens for the creator
	uncommittedToken, _ := commitments.GetUncommittedTokensForDenom(msg.Denom)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedTokens = append(uncommittedTokens, &types.UncommittedTokens{
			Denom:  msg.Denom,
			Amount: sdk.ZeroInt(),
		})
		commitments.UncommittedTokens = uncommittedTokens
	}
	// Update the uncommitted tokens amount
	uncommittedToken.Amount = uncommittedToken.Amount.Add(msg.Amount)

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	return &types.MsgDepositTokensResponse{}, nil
}
