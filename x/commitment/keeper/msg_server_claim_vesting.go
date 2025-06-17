package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v6/x/commitment/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

func (k msgServer) ClaimVesting(goCtx context.Context, msg *types.MsgClaimVesting) (*types.MsgClaimVestingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return k.Keeper.ClaimVesting(ctx, msg)
}

// ClaimVesting claims already vested amount
func (k Keeper) ClaimVesting(ctx sdk.Context, msg *types.MsgClaimVesting) (*types.MsgClaimVestingResponse, error) {
	// Get the Commitments for the sender
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	commitments := k.GetCommitments(ctx, sender)

	newClaims := sdk.Coins{}
	var updatedVestingTokens []*types.VestingTokens
	for _, vesting := range commitments.VestingTokens {
		vestedSoFar := vesting.VestedSoFar(ctx)                         // tokens unlocked
		newClaim := vestedSoFar.Sub(vesting.ClaimedAmount)              // tokens to mint or transfer
		newClaims = newClaims.Add(sdk.NewCoin(vesting.Denom, newClaim)) // adding coin to mint or transfer
		vesting.ClaimedAmount = vestedSoFar                             // updating claimed amount
		if !vesting.ClaimedAmount.Equal(vesting.TotalAmount) {          // if ClaimedAmount == TotalAmount, it would mean all tokens has been claimed and no need to keep the vesting tokens
			updatedVestingTokens = append(updatedVestingTokens, vesting)
		}
	}
	commitments.VestingTokens = updatedVestingTokens

	if newClaims.IsAllPositive() {
		// mint coins if vesting token is ELYS
		if newClaims.AmountOf(ptypes.Elys).IsPositive() {
			elysCoins := sdk.Coins{sdk.NewCoin(ptypes.Elys, newClaims.AmountOf(ptypes.Elys))}
			err := k.bankKeeper.MintCoins(ctx, types.ModuleName, elysCoins)
			if err != nil {
				return nil, err
			}
		}

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, newClaims)
		if err != nil {
			return nil, err
		}
	}

	k.SetCommitments(ctx, commitments)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimVesting,
			sdk.NewAttribute(types.AttributeCreator, msg.Sender),
			sdk.NewAttribute(types.AttributeAmount, newClaims.String()),
		),
	)

	return &types.MsgClaimVestingResponse{}, nil
}
