package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

var maxEdenAmountToClaim = sdkmath.NewInt(1000000000000)

func (k msgServer) ClaimRewardProgram(goCtx context.Context, msg *types.MsgClaimRewardProgram) (*types.MsgClaimRewardProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.ClaimAddress)
	rewardProgram := k.GetRewardProgram(ctx, sender)

	if rewardProgram.Amount.IsZero() {
		return nil, types.ErrRewardProgramNotFound
	}

	edenAmount := rewardProgram.Amount

	params := k.GetParams(ctx)
	if !params.EnableClaim {
		return nil, types.ErrClaimNotEnabled
	}

	if rewardProgram.Claimed {
		return nil, types.ErrRewardProgramAlreadyClaimed
	}

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartRewardProgramClaimHeight {
		return nil, types.ErrRewardProgramNotStarted
	}

	if currentHeight > params.EndRewardProgramClaimHeight {
		return nil, types.ErrRewardProgramEnded
	}

	total := k.GetTotalRewardProgramClaimed(ctx)
	total.TotalEdenClaimed = total.TotalEdenClaimed.Add(edenAmount)
	k.SetTotalRewardProgramClaimed(ctx, total)

	// Add eden to commitment
	err := k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, edenAmount)))
	if err != nil {
		return nil, err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, edenAmount)))
	if err != nil {
		return nil, err
	}

	rewardProgram.Claimed = true
	k.SetRewardProgram(ctx, rewardProgram)

	// This will never be triggered
	if total.TotalEdenClaimed.GT(maxEdenAmountToClaim) {
		return nil, types.ErrMaxEdenAmountReached
	}

	// Emit event for reward program claim
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimRewardProgram,
			sdk.NewAttribute(types.AttributeKeyClaimAddress, msg.ClaimAddress),
			sdk.NewAttribute(types.AttributeKeyEdenAmount, edenAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTotalEdenClaimed, total.TotalEdenClaimed.String()),
		),
	)

	return &types.MsgClaimRewardProgramResponse{
		EdenAmount: edenAmount,
	}, nil
}

func (k Keeper) BurnAirdropWallet(ctx sdk.Context) error {
	// Burn 990,250,400,000 uelys from airdrop wallet
	// we have allocated this amount in rewards program
	airdropWallet := "elys1wk7jwkqt2h9cnpkst85j9n454e4y8znlgk842n"
	airdropWalletAddress, err := sdk.AccAddressFromBech32(airdropWallet)
	if err != nil {
		return err
	}

	// transfer to module account
	err = k.SendCoinsFromAccountToModule(ctx, airdropWalletAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(990250400000))))
	if err != nil {
		return err
	}

	amountToBurn := sdkmath.NewInt(990250400000)

	if ctx.ChainID() == "elysicstestnet-1" {
		amountToBurn = sdkmath.NewInt(1000000000)
	}

	// burn the coins
	err = k.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, amountToBurn)))
	if err != nil {
		return err
	}

	// Add one time event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"burn-airdrop-wallet",
			sdk.NewAttribute("airdrop-wallet", airdropWallet),
			sdk.NewAttribute("amount-to-burn", amountToBurn.String()),
		),
	)

	return nil
}
