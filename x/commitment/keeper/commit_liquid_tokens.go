package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
)

// CommitLiquidTokens commit the tokens from user's balance
func (k Keeper) CommitLiquidTokens(ctx sdk.Context, addr sdk.AccAddress, denom string, amount math.Int, lockUntil uint64) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.CommitEnabled {
		return errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", denom)
	}

	depositCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// send the deposited coins to the module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, depositCoins)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("unable to send deposit tokens: %v", depositCoins))
	}

	// Update total commitment
	params := k.GetParams(ctx)
	params.TotalCommitted = params.TotalCommitted.Add(depositCoins...)
	k.SetParams(ctx, params)

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, addr)

	// Update the commitments
	commitments.AddCommittedTokens(denom, amount, lockUntil)
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	err = k.CommitmentChanged(ctx, addr, sdk.Coins{sdk.NewCoin(denom, amount)})
	if err != nil {
		return err
	}

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, addr.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}
