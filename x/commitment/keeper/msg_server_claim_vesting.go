package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// ClaimVesting claims already vested amount
func (k msgServer) ClaimVesting(goCtx context.Context, msg *types.MsgClaimVesting) (*types.MsgClaimVestingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the Commitments for the sender
	commitments := k.GetCommitments(ctx, msg.Sender)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	newClaims := sdk.Coins{}
	for i, vesting := range commitments.VestingTokens {
		vestedSoFar := vesting.VestedSoFar(ctx)
		newClaim := vestedSoFar.Sub(vesting.ClaimedAmount)
		newClaims = newClaims.Add(sdk.NewCoin(vesting.Denom, newClaim))
		commitments.VestingTokens[i].ClaimedAmount = vestedSoFar
	}

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
